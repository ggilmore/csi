package memtable

import (
	"bytes"
	"errors"
	"math/rand"
)

const maxLevel = 16

var ascendProbability = 0.5

type SkipList struct {
	level int
	Start *Node
	End   *Node
}

func NewSkipList() *SkipList {
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

func randomLevel() int {
	level := 1

	for rand.Float64() < ascendProbability && level < maxLevel {
		level++
	}

	return level
}

type Node struct {
	// isStart indicates that this node is the sentinel "starting" node
	isStart bool

	// isEnd indicates that this node is the sentinel "ending" node
	isEnd bool

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
