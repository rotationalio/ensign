package meta

import (
	"github.com/rotationalio/ensign/pkg/ensign/store/errors"
	ldbiter "github.com/syndtr/goleveldb/leveldb/iterator"
)

// Segments ensure that different objects are stored contiguously in the database
// ordered by their project then their ID to make it easy to scan for objects.
type Segment [2]byte

// Segments currently in use by Ensign
var (
	TopicSegment      = Segment{0x74, 0x70}
	TopicNamesSegment = Segment{0x54, 0x6e}
	TopicInfoSegment  = Segment{0x54, 0x69}
	GroupSegment      = Segment{0x47, 0x50}
)

func (s Segment) String() string {
	switch s {
	case TopicSegment:
		return "topic"
	case TopicNamesSegment:
		return "topic_name"
	case TopicInfoSegment:
		return "topic_info"
	case GroupSegment:
		return "group"
	default:
		return "unknown"
	}
}

// SegmentIterator wraps a leveldb iterator but skips any keys that do not match the
// specified segment. Since objects in the meta database are organized by project, this
// allows you to fetch all of the objects (e.g. topics or topic infos) across all
// projects without having to know the project ID in advance.
type SegmentIterator struct {
	ldbiter.Iterator
	segment Segment
}

func (i *SegmentIterator) Error() error {
	if err := i.Iterator.Error(); err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (i *SegmentIterator) Next() bool {
	for i.Iterator.Next() {
		// If the current object matches the segment, then break; otherwise continue
		// until we find an object with the specified segment.
		if i.Segment() == i.segment {
			return true
		}
	}
	return false
}

func (i *SegmentIterator) Prev() bool {
	for i.Iterator.Prev() {
		// If the current object matches the segment, then break; otherwise continue
		// until we find an object with the specified segment.
		if i.Segment() == i.segment {
			return true
		}
	}
	return false
}

func (i *SegmentIterator) Segment() Segment {
	key := i.Key()
	if len(key) == 34 {
		return Segment{key[16], key[17]}
	}
	return Segment{}
}
