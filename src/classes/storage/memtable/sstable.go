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

	start, end, eligible := s.findClosestBlock(key)
	if !eligible {
		return nil, KeyNotFound
	}

	return s.findKeyInBlock(start, end, key)
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
	startOffset, _, eligible := s.findClosestBlock(start)
	if !eligible {
		// quick hack - the iterator will return false immediately if we don't g
		startOffset = s.entriesEndOffset
	}

	return &SSIterator{
		reader: io.NewSectionReader(s.reader, int64(startOffset), int64(s.entriesEndOffset-startOffset)),
		start:  start,
		limit:  limit,
	}, nil
}

// findClosestBlock returns the starting and ending offsets of the first block that
// could contain "key" in its entries, along with a boolean that indicates whether
// an eligible block was found.
func (s *SSTable) findClosestBlock(key []byte) (startOffset, endOffset uint32, found bool) {
	endOffset = s.entriesEndOffset
	for i := len(s.sparseIndex) - 1; i >= 0; i-- {
		indexEntry := s.sparseIndex[i]
		if bytes.Compare(indexEntry.key, key) <= 0 {
			return indexEntry.offset, endOffset, true
		}

		endOffset = indexEntry.offset
	}

	return 0, 0, false
}

func (s *SSTable) findKeyInBlock(startOffset, endOffset uint32, targetKey []byte) (value []byte, err error) {
	r := io.NewSectionReader(s.reader, int64(startOffset), int64(endOffset-startOffset))

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

		var keyDeleted bool
		err := binary.Read(r, binary.LittleEndian, &keyDeleted)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while reading keyDeleted")
			}

			return nil, fmt.Errorf("reading keyDeleted: %s", err)
		}

		comparison := bytes.Compare(key, targetKey)
		if comparison > 0 {
			// key > targetKey - the key's not in the block if we haven't found it by now
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

		if comparison == 0 && !keyDeleted {
			// key == targetKey - read off the value and return it

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

		// if we reached here, then either key < targetKey or the key has been deleted.
		// In either case, we need to skip past this value.
		_, err = r.Seek(int64(valueLength), io.SeekCurrent)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("block corruption - hit end of block while seeking past value")
			}

			return nil, fmt.Errorf("seeking past value: %s", err)
		}

		if keyDeleted {
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

type SSIterator struct {
	reader *io.SectionReader

	start []byte
	limit []byte

	key   []byte
	value []byte

	err error
}

func (s *SSIterator) Next() bool {
	for {
		var keyLength uint32
		err := binary.Read(s.reader, binary.LittleEndian, &keyLength)
		if err != nil {
			if err != io.EOF {
				s.err = fmt.Errorf("reading key length: %w", err)
			}

			return false
		}

		key := make([]byte, keyLength)
		err = binary.Read(s.reader, binary.LittleEndian, &key)
		if err != nil {
			s.err = fmt.Errorf("reading key: %w", err)

			if err == io.EOF {
				s.err = fmt.Errorf("block corruption - hit end of block while reading key")
			}

			return false
		}

		if bytes.Compare(key, s.limit) >= 0 {
			// if we're past the limit, then there are no more elements to yield
			// (doesn't matter whether this element is deleted)
			return false
		}

		var isDeleted bool
		err = binary.Read(s.reader, binary.LittleEndian, &isDeleted)
		if err != nil {
			s.err = fmt.Errorf("reading isDeleted: %w", err)

			if err == io.EOF {
				s.err = fmt.Errorf("block corruption - hit end of block while reading isDeleted")
			}

			return false
		}

		var valueLength uint32
		err = binary.Read(s.reader, binary.LittleEndian, &valueLength)
		if err != nil {
			s.err = fmt.Errorf("reading value length: %s", err)

			if err == io.EOF {
				s.err = fmt.Errorf("block corruption - hit end of block while reading value length")
			}

			return false
		}

		if bytes.Compare(key, s.start) >= 0 && !isDeleted {
			// key >= s.start and is present, read value and store it
			value := make([]byte, valueLength)
			err = binary.Read(s.reader, binary.LittleEndian, &value)
			if err != nil {
				s.err = fmt.Errorf("reading value for key %q: %w", key, err)

				if err == io.EOF {
					s.err = fmt.Errorf("block corruption - hit end of block while reading value for key %q", key)
				}

				return false
			}

			s.key = key
			s.value = value

			return true
		}

		// if we reached here, then either the current key is <= than the start or the key is deleted.
		// In either case, we need to skip past the value.
		_, err = s.reader.Seek(int64(valueLength), io.SeekCurrent)
		if err != nil {
			s.err = fmt.Errorf("seeking past value for key %q: %w", key, err)

			if err == io.EOF {
				s.err = fmt.Errorf("block corruption - hit end of block while seeking past value for key %q: %w", key, err)
			}

			return false
		}
	}
}

func (s *SSIterator) Error() error {
	return s.err
}

func (s *SSIterator) Key() []byte {
	return s.key
}

func (s *SSIterator) Value() []byte {
	return s.value
}

type KeyDeletedError struct{}

func (k KeyDeletedError) Error() string {
	return "key has been deleted"
}

func (k KeyDeletedError) Is(target error) bool {
	return target == KeyNotFound
}
