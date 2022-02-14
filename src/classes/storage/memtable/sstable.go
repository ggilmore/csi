package memtable

import (
	"encoding/binary"
	"fmt"
	"io"
)

const skipListBlockSize = 1 << 11

type ImmutableDB interface {
	// Get gets the value for the given key. It returns an error if the
	// DB does not contain the key.
	Get(key []byte) (value []byte, err error)

	// Has returns true if the DB contains the given key.
	Has(key []byte) (ret bool, err error)

	// RangeScan returns an Iterator (see below) for scanning through all
	// key-value pairs in the given range, ordered by key ascending.
	RangeScan(start, limit []byte) (Iterator, error)
}

type SSTable struct {
	// offset of the first byte after the
	// sorted list of key-value entries
	entriesLength uint32

	file        io.ReadSeekCloser
	sparseIndex []sparseIndexEntry
}

// Plan: step through the sparse index, looking for the first index key
// that is < "key". Use the offset from the index to Seek() there and read until
// the boundary of the next block, or the end of the entriesList.
//
// TODO: I now realized that I  never ensured that the last entry for the SStable
// was in the index. I should adjust my splitting logic to ensure that happens
// since it makes the stepping logic cleaner. Now, I might have to open the last
// block no matter what I do since I don't know what the last key in the index is.

func (S SSTable) Get(key []byte) (value []byte, err error) {
	//TODO implement me
	panic("implement me")
}

// Plan: Just use _get_ and throw away the value

func (S SSTable) Has(key []byte) (ret bool, err error) {
	//TODO implement me
	panic("implement me")
}

// Plan: Use the logic from get to find the closest key that's <= start.
// from there, the iterator and keep yielding values until the current key
// >= limit, or the end of the entries list is reached.

func (S SSTable) RangeScan(start, limit []byte) (Iterator, error) {
	//TODO implement me
	panic("implement me")
}

func OpenSStable(r io.ReadSeekCloser) (*SSTable, error) {
	indexEnd, err := r.Seek(-4, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seeking to end of file to read index start location: %s", err)
	}

	var indexStart uint32
	err = binary.Read(r, binary.LittleEndian, &indexStart)
	if err != nil {
		return nil, fmt.Errorf("reading index start location: %s", err)
	}

	if uint32(indexEnd) > indexStart {
		return nil, fmt.Errorf("corrupted ss table file: index end (%d) > index start (%d)", uint32(indexEnd), indexStart)
	}

	if uint32(indexEnd) == indexStart {
		// this ss table must be empty - we have no index entries
		return &SSTable{
			file: r,
		}, nil
	}

	var sparseIndex []sparseIndexEntry

	for i := indexStart; i < uint32(indexEnd); i++ {
		var keyLength uint32
		err := binary.Read(r, binary.LittleEndian, &keyLength)
		if err != nil {
			return nil, fmt.Errorf("failed to key length in index: %s", err)
		}

		key := make([]byte, keyLength)
		_, err = r.Read(key)
		if err != nil {
			return nil, fmt.Errorf("reading key: %s", err)
		}

		var blockOffset uint32
		err = binary.Read(r, binary.LittleEndian, &blockOffset)
		if err != nil {
			return nil, fmt.Errorf("reading block offset for key %s: %s", key, err)
		}

		sparseIndex = append(sparseIndex, sparseIndexEntry{
			key:    key,
			offset: blockOffset,
		})
	}

	return &SSTable{
		file:          r,
		entriesLength: indexStart,
		sparseIndex:   sparseIndex,
	}, nil
}
