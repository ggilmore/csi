package memtable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
)

const (
	maxLevel          = 16
	ascendProbability = 0.5
)

type SkipList struct {
	level int
	Start *Node
	End   *Node

	Size             int
	ssTableBlockSize int
}

type SkipListOptions struct {
	SSTableBlockSize int
}

func NewSkipList(options SkipListOptions) *SkipList {
	o := SkipListOptions{
		SSTableBlockSize: defaultSSTableBlockSize,
	}

	if options.SSTableBlockSize > 0 {
		o.SSTableBlockSize = options.SSTableBlockSize
	}

	start := &Node{
		isStart: true,
	}
	end := &Node{
		isEnd: true,
	}

	// seed the skip list by linking start->end on all levels
	for i := 0; i < maxLevel; i++ {
		start.forward[i] = end
	}

	return &SkipList{
		level: 1,
		Start: start,
		End:   end,

		ssTableBlockSize: o.SSTableBlockSize,
	}
}

func (s *SkipList) Get(key []byte) (value []byte, err error) {
	newNode := &Node{
		Key: key,
	}

	currentNode := s.Start
	// starting from the top and going down, find the rightmost node
	// that's still less than "key"
	for i := s.level; i >= 1; i-- {
		index := i - 1
		for {
			next := currentNode.forward[index]
			if next == nil || NodeCompare(next, newNode) >= 0 {
				break
			}

			currentNode = next
		}
	}

	// is the next node after this one "key"?
	currentNode = currentNode.forward[0]
	if currentNode != nil && NodeCompare(currentNode, newNode) == 0 {
		return currentNode.Value, nil

	}

	return nil, KeyNotFound
}

func (s *SkipList) Has(key []byte) (ret bool, err error) {
	_, err = s.Get(key)
	if err != nil {
		if errors.Is(err, KeyNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *SkipList) Put(key, value []byte) error {
	newNode := &Node{
		Key:   key,
		Value: value,
	}

	currentNode := s.Start
	update := [maxLevel]*Node{}

	// starting from the top and going down, find the rightmost node
	// that's still less than "key"
	for i := s.level; i >= 1; i-- {
		index := i - 1
		for {
			next := currentNode.forward[index]
			if next == nil || NodeCompare(next, newNode) >= 0 {
				break
			}

			currentNode = next
		}

		// remember the nodes that we used to "descend" a level
		update[index] = currentNode
	}

	// peek at the next node after this one
	currentNode = currentNode.forward[0]
	if currentNode != nil && NodeCompare(currentNode, newNode) == 0 {
		// current node contains your key, just update it
		currentNode.Value = newNode.Value
		return nil
	}

	// we didn't find our key, insert it now

	// figure out what level this new node is at
	newLevel := randomLevel()
	if newLevel > s.level {

		// if this new node is higher than all the other ones,
		// we need to make sure that the starting node's forward pointers
		// are updated to point to the node we're adding
		for i := s.level + 1; i <= newLevel; i++ {
			index := i - 1
			update[index] = s.Start
		}

		s.level = newLevel
	}

	// splice ourselves in between the nodes we used to descend and the nodes
	// they were pointing to before
	for i := 1; i <= newLevel; i++ {
		index := i - 1
		newNode.forward[index] = update[index].forward[index]
		update[index].forward[index] = newNode
	}

	s.Size += len(key) + len(value)

	return nil
}

func (s *SkipList) Delete(key []byte) error {
	newNode := &Node{
		Key: key,
	}

	currentNode := s.Start
	update := [maxLevel]*Node{}

	// starting from the top and going down, find the rightmost node
	// that's still less than "key"
	for i := s.level; i >= 1; i-- {
		index := i - 1
		for {
			next := currentNode.forward[index]
			if next == nil || NodeCompare(next, newNode) >= 0 {
				break
			}

			currentNode = next
		}

		// remember the nodes that we used to "descend" a level
		update[index] = currentNode
	}

	// peek at the next node after this one
	currentNode = currentNode.forward[0]
	if currentNode == nil || NodeCompare(currentNode, newNode) != 0 {
		// we don't have the key - there is nothing to do
		return nil
	}

	// we have the key, splice ourselves out
	for i := 1; i <= s.level; i++ {
		index := i - 1
		if update[index].forward[index] != currentNode {
			break
		}

		update[index].forward[index] = currentNode.forward[index]
	}

	for s.level >= 1 && s.Start.forward[s.level-1] == s.End {
		s.level--
	}

	s.Size -= len(currentNode.Key) + len(currentNode.Value)

	return nil
}

func (s *SkipList) RangeScan(start, limit []byte) (Iterator, error) {
	fakeNode := &Node{
		Key: start,
	}

	currentNode := s.Start
	// starting from the top and going down, find the rightmost node
	// that's still less than "start"
	for i := s.level; i >= 1; i-- {
		index := i - 1
		for {
			next := currentNode.forward[index]
			if next == nil || NodeCompare(next, fakeNode) >= 0 {
				break
			}

			currentNode = next
		}
	}

	return &SkipIterator{
		skip:    s,
		current: currentNode,
		limit:   limit,
	}, nil
}

// Each SSTable looks like this:
//
// There are three sections:
// entries_list - list of key-value pairs, each field is prefixed by its length: [key_length (uint32)][key][value_length (uint32)][value]...
// sparse_index - sorted list of pairs [key][block_offset_in_file (uint32)]. As described in DDIA, this is a much sparser list
// 				  than a full hashtable. I create a new entry in the index every 4k bytes.
// index_location: location in file where the sparse index starts (offset (uint32))
//
// since the sparse_index is of variable length, putting it at the bottom of the file means that you can write it after you
// know all the offsets of the key-value entries without to recalculate them (writing the offset at the beggining would require
// you to bump the offsets by the length of the index
//
// you can read the file by 1) reading the 4 byte index_location at the end of the file, 2) parsing the sparse_index
//                          3) doing a for loop over the index to find the keys you want

func (s *SkipList) flushSSTable(w io.Writer) error {
	writer := offsetWriter{
		wrappedWriter: w,
	}

	var sparseIndex []sparseIndexEntry

	nextCheckpointBytes := uint32(0)
	node := s.Start.forward[0]
	for node != s.End {
		// write [key_length][key][isDeleted][value_length][value] for each entry
		startingOffset := writer.Offset

		if nextCheckpointBytes <= startingOffset || node.forward[0] == s.End {
			sparseIndex = append(sparseIndex, sparseIndexEntry{
				key:    node.Key,
				offset: startingOffset,
			})

			nextCheckpointBytes = startingOffset + uint32(s.ssTableBlockSize)
		}

		err := writer.WriteUint32(uint32(len(node.Key)))
		if err != nil {
			return fmt.Errorf("writing length (%d) of key %q: %w", len(node.Key), node.Key, err)
		}

		err = writer.Write(node.Key)
		if err != nil {
			return fmt.Errorf("writing key %q: %w", node.Key, err)
		}

		err = writer.WriteUint32(uint32(len(node.Value)))
		if err != nil {
			return fmt.Errorf("writing length (%d) of value %q: %w", len(node.Value), node.Value, err)
		}

		err = writer.Write(node.Value)
		if err != nil {
			return fmt.Errorf("writing value %q: %w", node.Key, err)
		}

		node = node.forward[0]
	}

	sparseIndexOffset := writer.Offset

	for _, e := range sparseIndex {
		// write [key_length][key][offset] for all entries in the sparse index
		err := writer.WriteUint32(uint32(len(e.key)))
		if err != nil {
			return fmt.Errorf("writing length (%d) of key %q in sparse index: %w", len(e.key), e.key, err)
		}

		err = writer.Write(e.key)
		if err != nil {
			return fmt.Errorf("writing key %q in sparse index: %w", e.key, err)
		}

		err = writer.WriteUint32(e.offset)
		if err != nil {
			return fmt.Errorf("writing offset %d in sparse index: %w", e.offset, err)
		}
	}

	err := writer.WriteUint32(sparseIndexOffset)
	if err != nil {
		return fmt.Errorf("writing starting offset (%d) for sparse index: %s", sparseIndexOffset, err)
	}

	return nil
}

type Node struct {
	// isStart indicates that this node is the sentinel "starting" node
	isStart bool

	// isEnd indicates that this node is the sentinel "ending" node
	isEnd bool

	// isDeleted indicates that the key for this node has been deleted
	isDeleted bool

	Key   []byte
	Value []byte

	forward [maxLevel]*Node
}

// NodeCompare compares the keys of two nodes and uses the same return
// value semantics as bytes.Compare.
func NodeCompare(a, b *Node) int {
	// massage any nil arguments
	if a == nil {
		a = &Node{}
	}
	if b == nil {
		b = &Node{}
	}

	if a.isStart {
		if b.isStart {
			return 0
		}

		// everything is larger than the start node (other than the start node)
		return -1
	}

	if a.isEnd {
		if b.isEnd {
			return 0
		}

		// everything is smaller than the end node (other than the end node)
		return 1
	}

	// fall back to comparing the keys
	return bytes.Compare(a.Key, b.Key)
}

type SkipIterator struct {
	skip *SkipList

	start []byte
	limit []byte

	current *Node
}

func (s *SkipIterator) Next() bool {
	startOfRange := &Node{
		Key: s.start,
	}

	endOfRange := &Node{
		Key: s.limit,
	}

	// TODO handle the case where a node is deleted while we are range scanning
	// (is deleted flag?)
	// search from scratch every time?

	// peek at the next highest node and see if it's within [start, limit)
	next := s.current.forward[0]
	if next != nil && NodeCompare(next, startOfRange) >= 0 && NodeCompare(next, endOfRange) <= -1 {
		s.current = next
		return true
	}

	s.current = &Node{}
	return false
}

func (s *SkipIterator) Error() error {
	return nil
}

func (s *SkipIterator) Key() []byte {
	return s.current.Key
}

func (s *SkipIterator) Value() []byte {
	return s.current.Value
}

type sparseIndexEntry struct {
	key    []byte
	offset uint32
}

// offsetWriter is an io.Writer that keeps track of the current Offset
type offsetWriter struct {
	Offset uint32

	wrappedWriter io.Writer
}

func (o *offsetWriter) Write(b []byte) error {
	n, err := o.wrappedWriter.Write(b)
	if err != nil {
		return err
	}

	o.Offset += uint32(n)
	return nil
}

func (o *offsetWriter) WriteUint32(n uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], n)
	return o.Write(buf[:])
}

func (o *offsetWriter) WriteBool(b bool) error {
	buf := [1]byte{0}
	if b {
		buf[0] = 1
	}
	return o.Write(buf[:])
}

func randomLevel() int {
	level := 1

	for rand.Float64() < ascendProbability && level < maxLevel {
		level++
	}

	return level
}
