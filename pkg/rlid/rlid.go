package rlid

import "time"

/*
An RLID is a 10 byte totally ordered unique lexicographically sortable identifier.
The components are encoded in 10 octets where each component is encoded with the MSB
first (network byte order). The first 48 bits encode a timestamp at millisecond
granularity and the next 32 bits encode a montonically increasing sequence number that
allows for ~4.294e12 unique ids to be generated per second with timestamps that will
last until around 10,000 years after 1970.

RLIDs are string encoded using Crockford's base32 (5 bits per character) meaning that
each RLID is 16 characters long in string format.
*/
type RLID [10]byte

// Make returns an RLID with the current timestamp and the given sequence number.
// Make will panic if the current timestamp is too far in the future for an RLID.
func Make(seq uint32) (id RLID) {
	if err := id.SetTime(Now()); err != nil {
		panic(err)
	}

	if err := id.SetSequence(seq); err != nil {
		panic(err)
	}

	return id
}

//===========================================================================
// Time Functionality
//===========================================================================

// Time returns the Unix epoch time in milliseconds that is encoded in the ID.
// Use the top level Time function to convert into a time.Time.
func (id RLID) Time() uint64 {
	return uint64(id[5]) | uint64(id[4])<<8 | uint64(id[3])<<16 | uint64(id[2])<<24 | uint64(id[1])<<32 | uint64(id[0])<<40
}

// SetTime sets the time component of the RLID to the given Unix epoch time in milliseconds.
func (id *RLID) SetTime(ms uint64) error {
	if ms > MaxTime {
		return ErrOverTime
	}

	(*id)[0] = byte(ms >> 40)
	(*id)[1] = byte(ms >> 32)
	(*id)[2] = byte(ms >> 24)
	(*id)[3] = byte(ms >> 16)
	(*id)[4] = byte(ms >> 8)
	(*id)[5] = byte(ms)

	return nil
}

// The maximum time in milliseconds that can be represented in an RLID.
const MaxTime uint64 = 281474976710655

// Now is a convenience function to return the current time in Unix epoch time in
// milliseconds. Unix epoch time is the number of milliseconds since January 1, 1970 at
// midnight UTC. Note that Unix epoch time will always be in the UTC timezone.
func Now() uint64 {
	return uint64(time.Now().UTC().UnixMilli())
}

// Timestamp converts a time.Time to Unix epoch time in milliseconds.
// NOTE: RLIDs cannot store a full Unix millisecond epoch time so timestamps after
// 10889-08-02 05:31:50.656 UTC are undefined.
func Timestamp(t time.Time) uint64 {
	return uint64(t.UnixMilli())
}

// Time converts a Unix epoch time in milliseconds back to a time.Time.
func Time(ms uint64) time.Time {
	return time.UnixMilli(int64(ms))
}

//===========================================================================
// Sequence Functionality
//===========================================================================

// Sequence returns the montonically increasing sequence number component.
func (id RLID) Sequence() uint32 {
	return uint32(id[9]) | uint32(id[8])<<8 | uint32(id[7])<<16 | uint32(id[6])<<24
}

// SetSequence sets the sequence component of the RLID to the given integer.
func (id *RLID) SetSequence(seq uint32) error {
	(*id)[6] = byte(seq >> 24)
	(*id)[7] = byte(seq >> 16)
	(*id)[8] = byte(seq >> 8)
	(*id)[9] = byte(seq)

	return nil
}

//===========================================================================
// Encoding
//===========================================================================
