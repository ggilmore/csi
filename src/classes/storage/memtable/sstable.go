package memtable

import (
	"bytes"
	"encoding/binary"
	"errors"
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

// ReadAtSeekerCloser combines the io.ReaderAt, io.Seeker, and io.Closer interfaces.
type ReadAtSeekerCloser interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type SSTable struct {
	// offset of the first byte after the
	// sorted list of key-value entries
	entriesLength uint32

	reader      ReadAtSeekerCloser
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

func (s *SSTable) Get(key []byte) (value []byte, err error) {
	//if len(s.sparseIndex) == 0 {
	//	// SSTable is empty
	//	return nil, KeyNotFound
	//}
	//
	//switch len(s.sparseIndex) {
	//case 0:
	//	return nil, KeyNotFound
	//case 1:
	//	return s.findKeyInBlock(s.sparseIndex[0].offset, s.entriesLength, key)
	//default:
	//
	//}

	// use the super slow non-indexed method
	return s.findKeyInBlock(0, s.entriesLength, key)

}

// Plan: Just use _get_ and throw away the value

func (s *SSTable) Has(key []byte) (ret bool, err error) {
	_, err = s.Get(key)

	if errors.Is(err, KeyNotFound) {
		return false, nil
	}

	return true, nil
}

// Plan: Use the logic from get to find the closest key that's <= start.
// from there, the iterator and keep yielding values until the current key
// >= limit, or the end of the entries list is reached.

func (s *SSTable) RangeScan(start, limit []byte) (Iterator, error) {
	//TODO implement me
	panic("implement me")
}

func (s *SSTable) findKeyInBlock(start, limit uint32, needle []byte) (value []byte, err error) {
	r := io.NewSectionReader(s.reader, int64(start), int64(limit-start))

	for {
		var keyLength uint32
		err = binary.Read(r, binary.LittleEndian, &keyLength)
		if err != nil {
			if err == io.EOF {
				// we reached the end of the block
				return nil, KeyNotFound
			}

			return nil, fmt.Errorf("reading key length: %s", err)
		}

		key := make([]byte, keyLength)
		err = binary.Read(r, binary.LittleEndian, &key)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while reading key")
			}

			return nil, fmt.Errorf("reading key: %s", err)
		}

		comparison := bytes.Compare(key, needle)
		if comparison > 0 {
			// current key is greater than the needle - the key's not in the block
			// if we haven't found it by now
			return nil, KeyNotFound
		}

		var valueLength uint32
		err = binary.Read(r, binary.LittleEndian, &valueLength)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while reading value length")
			}

			return nil, fmt.Errorf("reading value length: %s", err)
		}

		if comparison == 0 {
			// the key == needle, read off the value and return it
			value := make([]byte, valueLength)
			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("block corruption - hit end of block while reading value")
				}

				return nil, fmt.Errorf("reading value: %s", err)
			}

			return value, nil
		}

		// if we reached here, then we're  reading values that are smaller than the key.
		// We need to keep scanning the other entries in the block. Just skip past the next value then move on.
		//
		// Note, we could combine this case with the key == needles case, but
		// allocating the value slice would create a ton of small unnecessary allocations
		_, err := r.Seek(int64(valueLength), io.SeekCurrent)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while seeking past value")
			}

			return nil, fmt.Errorf("seeking past value: %s", err)
		}
	}
}

func OpenSStable(r ReadAtSeekerCloser) (*SSTable, error) {
	indexEnd, err := r.Seek(-4, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seeking to end of file to read index start location: %s", err)
	}

	var indexStart uint32
	err = binary.Read(r, binary.LittleEndian, &indexStart)
	if err != nil {
		return nil, fmt.Errorf("reading index start location: %s", err)
	}

	if indexStart > uint32(indexEnd) {
		return nil, fmt.Errorf("corrupted ss table file: index end (%d) > index start (%d)", uint32(indexEnd), indexStart)
	}

	if uint32(indexEnd) == indexStart {
		// this ss table must be empty - we have no index entries
		return &SSTable{
			reader: r,
		}, nil
	}

	var sparseIndex []sparseIndexEntry

	_, err = r.Seek(int64(indexStart), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seeking to start of index: %s", err)
	}

	sectionReader := io.NewSectionReader(r, int64(indexStart), indexEnd-int64(indexStart))

	for {
		// read all
		var keyLength uint32
		err := binary.Read(sectionReader, binary.LittleEndian, &keyLength)
		if err != nil {
			if err == io.EOF {
				// we've hit the end
				break
			}

			return nil, fmt.Errorf("reading key length in index: %s", err)
		}

		key := make([]byte, keyLength)
		_, err = sectionReader.Read(key)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("index corruption - EOF while reading key")
			}
			return nil, fmt.Errorf("reading key: %s", err)
		}

		var blockOffset uint32
		err = binary.Read(sectionReader, binary.LittleEndian, &blockOffset)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("index corruption - EOF while reading value")
			}

			return nil, fmt.Errorf("reading block offset for key %s: %s", key, err)
		}

		sparseIndex = append(sparseIndex, sparseIndexEntry{
			key:    key,
			offset: blockOffset,
		})
	}

	return &SSTable{
		reader:        r,
		entriesLength: indexStart,
		sparseIndex:   sparseIndex,
	}, nil
}
