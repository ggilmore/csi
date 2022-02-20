package memtable

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type CombinedSkipAndSS struct {
	ssTables []*SSTable

	skipList         *SkipList
	deletionSkipList *SkipList

	SkipListSizeThreshold int
	SSTableBlockSize      int
	SSTableDir            string
}

type CombinedSkipAndSSOptions struct {
	SkipListSizeThreshold int

	SSTableBlockSize int
	SSTableDir       string
}

func NewCombinedSkipAndSS(options CombinedSkipAndSSOptions) (*CombinedSkipAndSS, error) {
	o := CombinedSkipAndSSOptions{
		SkipListSizeThreshold: 1 << 21, // ~2MB
		SSTableDir:            os.TempDir(),
		SSTableBlockSize:      defaultSSTableBlockSize,
	}

	if options.SkipListSizeThreshold > 0 {
		o.SkipListSizeThreshold = options.SkipListSizeThreshold
	}

	if options.SSTableDir != "" {
		o.SSTableDir = options.SSTableDir
	}

	if options.SSTableBlockSize > 0 {
		o.SSTableBlockSize = options.SSTableBlockSize
	}

	out := CombinedSkipAndSS{
		skipList:         NewSkipList(SkipListOptions{SSTableBlockSize: o.SSTableBlockSize}),
		deletionSkipList: NewSkipList(SkipListOptions{SSTableBlockSize: o.SSTableBlockSize}),

		SSTableDir:            filepath.Clean(o.SSTableDir),
		SkipListSizeThreshold: o.SkipListSizeThreshold,
		SSTableBlockSize:      o.SSTableBlockSize,
	}

	// if the index directory exists, load in all the existing sstables
	if _, err := os.Stat(filepath.Clean(out.SSTableDir)); err == nil {
		var sstables []*SSTable

		ssTableFiles, err := filepath.Glob(filepath.Join(out.SSTableDir, "*.sstable"))
		if err != nil {
			return nil, fmt.Errorf("listing all sstables in index dir %q: %w", options.SSTableDir, err)
		}

		// sort in increasing numerical order
		sort.Strings(ssTableFiles)

		for _, f := range ssTableFiles {
			file, err := os.Open(f)
			if err != nil {
				return nil, fmt.Errorf("opening sstable underlying file %q: %w", f, err)
			}

			ss, err := OpenSStable(file)
			if err != nil {
				return nil, fmt.Errorf("parsing sstable file %q: %w", file.Name(), err)
			}

			sstables = append(sstables, ss)
		}

		out.ssTables = sstables
	}

	return &out, nil
}

func (c *CombinedSkipAndSS) Get(key []byte) (value []byte, err error) {
	contains, err := c.deletionSkipList.Has(key)
	if err != nil {
		return nil, fmt.Errorf("checking deletion skip list: %w", err)
	}

	if contains {
		return nil, KeyNotFound
	}

	// data sources are checked from newest to oldest
	dbs := []ImmutableDB{c.skipList}
	for i := len(c.ssTables) - 1; i >= 0; i-- {
		dbs = append(dbs, c.ssTables[i])
	}

	for i, db := range dbs {
		value, err := db.Get(key)
		if err != nil {
			if errors.Is(err, KeyDeletedError{}) {
				// short circuit if we find out that the key is
				// deleted
				return nil, KeyNotFound
			}

			if errors.Is(err, KeyNotFound) {
				continue
			}

			return nil, fmt.Errorf("querying db #%d: %w", i, err)
		}

		return value, nil
	}

	return nil, KeyNotFound
}

func (c *CombinedSkipAndSS) Has(key []byte) (ret bool, err error) {
	_, err = c.Get(key)
	if err != nil {
		if errors.Is(err, KeyNotFound) {
			return false, err
		}

		return false, err
	}

	return true, nil
}

func (c *CombinedSkipAndSS) Put(key, value []byte) error {
	err := c.deletionSkipList.Delete(key)
	if err != nil {
		return fmt.Errorf("removing key from deletion skiplist: %w", err)
	}

	err = c.skipList.Put(key, value)
	if err != nil {
		return fmt.Errorf("inserting into skiplist: %w", err)
	}

	if (c.skipList.Size + c.deletionSkipList.Size) >= c.SkipListSizeThreshold {
		err := os.MkdirAll(c.SSTableDir, 0666)
		if err != nil {
			return fmt.Errorf("creating directory %q: %w", c.SSTableDir, err)
		}

		tempFile, err := os.CreateTemp(c.SSTableDir, "new-*.sstable")
		if err != nil {
			return fmt.Errorf("creating temporary file for sstable %w", err)
		}

		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		err = c.flushtoSStable(tempFile)
		if err != nil {
			return fmt.Errorf("flushing sstable to file %q: %w", tempFile.Name(), err)
		}

		finalPath := filepath.Join(c.SSTableDir, fmt.Sprintf("%d.sstable", len(c.ssTables)))
		err = os.Rename(tempFile.Name(), finalPath)
		if err != nil {
			return fmt.Errorf("renaming temporary sstable file %q to final sstable file %q: %w", tempFile.Name(), finalPath, err)
		}

		tempFile.Close()
		file, err := os.Open(finalPath)
		if err != nil {
			return fmt.Errorf("opening sstable file %q: %w", finalPath, err)
		}

		ssTable, err := OpenSStable(file)
		if err != nil {
			return fmt.Errorf("opening sstable from file %q: %s", file.Name(), err)
		}

		c.ssTables = append(c.ssTables, ssTable)
		c.skipList = NewSkipList(SkipListOptions{})
		c.deletionSkipList = NewSkipList(SkipListOptions{})
	}

	return nil
}

func (c *CombinedSkipAndSS) flushtoSStable(w io.Writer) error {
	writer := offsetWriter{
		wrappedWriter: w,
	}

	var sparseIndex []sparseIndexEntry

	var lastWrittenNode *Node
	lastStartingOffset := uint32(0)

	nextCheckpointBytes := uint32(0)

	var writeNode = func(n *Node, isDeleted bool) error {
		startingOffset := writer.Offset

		if nextCheckpointBytes <= startingOffset {
			sparseIndex = append(sparseIndex, sparseIndexEntry{
				key:    n.Key,
				offset: startingOffset,
			})

			nextCheckpointBytes = startingOffset + uint32(c.SSTableBlockSize)
		}

		// write [key_length][key][isDeleted][value_length][value]
		err := writer.WriteUint32(uint32(len(n.Key)))
		if err != nil {
			return fmt.Errorf("writing length (%d) of key %q: %w", len(n.Key), n.Key, err)
		}

		err = writer.Write(n.Key)
		if err != nil {
			return fmt.Errorf("writing key %q: %w", n.Key, err)
		}

		err = writer.WriteBool(isDeleted)
		if err != nil {
			return fmt.Errorf("writing isDeleted %t: %w", isDeleted, err)
		}

		err = writer.WriteUint32(uint32(len(n.Value)))
		if err != nil {
			return fmt.Errorf("writing length (%d) of value %q: %w", len(n.Value), n.Value, err)
		}

		err = writer.Write(n.Value)
		if err != nil {
			return fmt.Errorf("writing value %q: %w", n.Key, err)
		}

		lastWrittenNode = n
		lastStartingOffset = startingOffset
		return nil
	}

	currentNormalNode := c.skipList.Start.forward[0]
	currentDeletedNode := c.deletionSkipList.Start.forward[0]

	for currentNormalNode != c.skipList.End && currentDeletedNode != c.deletionSkipList.End {
		if bytes.Compare(currentNormalNode.Key, currentDeletedNode.Key) > 0 {
			err := writeNode(currentDeletedNode, true)
			if err != nil {
				return fmt.Errorf("writing deleted node: %s", err)
			}

			currentDeletedNode = currentDeletedNode.forward[0]
			continue
		}

		err := writeNode(currentNormalNode, true)
		if err != nil {
			return fmt.Errorf("writing normal node: %s", err)
		}

		currentNormalNode = currentNormalNode.forward[0]
	}

	for currentNormalNode != c.skipList.End {
		err := writeNode(currentNormalNode, true)
		if err != nil {
			return fmt.Errorf("writing normal node: %s", err)
		}

		currentNormalNode = currentNormalNode.forward[0]
	}

	for currentDeletedNode != c.deletionSkipList.End {
		err := writeNode(currentDeletedNode, true)
		if err != nil {
			return fmt.Errorf("writing deleted node: %s", err)
		}

		currentDeletedNode = currentDeletedNode.forward[0]
	}

	if len(sparseIndex) > 0 && bytes.Compare(sparseIndex[len(sparseIndex)-1].key, lastWrittenNode.Key) != 0 {
		sparseIndex = append(sparseIndex, sparseIndexEntry{
			key:    lastWrittenNode.Key,
			offset: lastStartingOffset,
		})
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

func (c CombinedSkipAndSS) Delete(key []byte) error {
	contains, err := c.skipList.Has(key)
	if err != nil {
		return fmt.Errorf("removing key from skiplist: %w", err)

	}

	if !contains {
		err := c.skipList.Delete(key)
		if err != nil {
			return fmt.Errorf("removing key from skiplist: %w", err)
		}

		return nil
	}

	err = c.deletionSkipList.Put(key, nil)
	if err != nil {
		return fmt.Errorf("puting key in deletion skiplist: %w", err)
	}

	return nil
}

func (c CombinedSkipAndSS) RangeScan(start, limit []byte) (Iterator, error) {
	//TODO implement me
	panic("implement me")
}
