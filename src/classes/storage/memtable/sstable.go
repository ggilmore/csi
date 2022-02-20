package memtable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const defaultSSTableBlockSize = 1 << 11

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
	entriesEndOffset uint32

	reader      ReadAtSeekerCloser
	sparseIndex []sparseIndexEntry
}

func (s *SSTable) Get(key []byte) (value []byte, err error) {
	if len(s.sparseIndex) == 0 {
		return nil, KeyNotFound
	}

	endOffset := s.entriesEndOffset
	for i := len(s.sparseIndex) - 1; i >= 0; i-- {
		indexEntry := s.sparseIndex[i]
		if bytes.Compare(indexEntry.key, key) <= 0 {
			return s.findKeyInBlock(indexEntry.offset, endOffset, key)
		}

		endOffset = indexEntry.offset
	}

	return nil, KeyNotFound
}

func (s *SSTable) Has(key []byte) (ret bool, err error) {
	_, err = s.Get(key)
	if err != nil {
		if errors.Is(err, KeyNotFound) {
			return false, nil
		}

		return false, err
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

type KeyDeletedError struct{}

func (k KeyDeletedError) Error() string {
	return "key has been deleted"
}

func (k KeyDeletedError) Is(target error) bool {
	return target == KeyNotFound
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

		var isDeleted bool
		err := binary.Read(r, binary.LittleEndian, &isDeleted)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while reading isDeleted")
			}

			return nil, fmt.Errorf("reading isDeleted: %s", err)
		}

		comparison := bytes.Compare(key, needle)
		if comparison > 0 {
			// key > needle - the key's not in the block if we haven't found it by now
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

		if comparison == 0 && !isDeleted {
			// key == needle - read off the value and return it

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

		// if we reached here, then we're reading values that are smaller than the key or the key has been deleted.
		// In either case, we need to skip past the value.
		_, err = r.Seek(int64(valueLength), io.SeekCurrent)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while seeking past value")
			}

			return nil, fmt.Errorf("seeking past value: %s", err)
		}

		if isDeleted {
			return nil, KeyDeletedError{}
		}
	}
}

func OpenSStable(r ReadAtSeekerCloser) (*SSTable, error) {
	indexEnd, err := r.Seek(-4, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("seeking to end of file to read index start location: %s", err)
	}

	var s uint32
	err = binary.Read(r, binary.LittleEndian, &s)
	if err != nil {
		return nil, fmt.Errorf("reading index start location: %s", err)
	}
	indexStart := int64(s)

	if indexStart > indexEnd {
		return nil, fmt.Errorf("corrupted ss table file: index end (%d) > index start (%d)", indexEnd, indexStart)
	}

	if indexEnd == indexStart {
		// this ss table must be empty - we have no index entries
		return &SSTable{
			reader: r,
		}, nil
	}

	var sparseIndex []sparseIndexEntry

	_, err = r.Seek(indexStart, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seeking to start of index: %s", err)
	}

	indexReader := io.NewSectionReader(r, indexStart, indexEnd-indexStart)

	for {
		var keyLength uint32
		err := binary.Read(indexReader, binary.LittleEndian, &keyLength)
		if err != nil {
			if err == io.EOF {
				// we've hit the end
				break
			}

			return nil, fmt.Errorf("reading key length in index: %s", err)
		}

		key := make([]byte, keyLength)
		_, err = indexReader.Read(key)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("index corruption - EOF while reading key")
			}
			return nil, fmt.Errorf("reading key: %s", err)
		}

		var blockOffset uint32
		err = binary.Read(indexReader, binary.LittleEndian, &blockOffset)
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
		reader:           r,
		entriesEndOffset: uint32(indexStart),
		sparseIndex:      sparseIndex,
	}, nil
}
