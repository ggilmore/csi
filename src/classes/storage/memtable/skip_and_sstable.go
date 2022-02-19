package memtable

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type CombinedSkipAndSS struct {
	ssTables []*SSTable
	skipList *SkipList

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
		skipList: NewSkipList(SkipListOptions{SSTableBlockSize: o.SSTableBlockSize}),

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
	dbs := []ImmutableDB{c.skipList}
	for i := len(c.ssTables) - 1; i >= 0; i-- {
		dbs = append(dbs, c.ssTables[i])
	}

	for i, db := range dbs {
		value, err := db.Get(key)
		if err != nil {
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
	err := c.skipList.Put(key, value)
	if err != nil {
		return fmt.Errorf("inserting into skiplist: %w", err)
	}

	if c.skipList.Size >= c.SkipListSizeThreshold {
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

		err = c.skipList.flushSSTable(tempFile)
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
		c.skipList = NewSkipList(SkipListOptions{SSTableBlockSize: c.SSTableBlockSize})
	}

	return nil
}

func (c CombinedSkipAndSS) Delete(key []byte) error {
	//TODO implement me
	panic("implement me")
}

func (c CombinedSkipAndSS) RangeScan(start, limit []byte) (Iterator, error) {
	//TODO implement me
	panic("implement me")
}
