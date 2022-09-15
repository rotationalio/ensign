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

// Parse an RLID from an encoded string.
// ErrDataSize is returned if the length of the string is different from the expected
// encoded length of RLIDs. Invalid encodings produce undefined RLIDs. For a version
// that validates the RLID, use ParseStrict.
func Parse(rlid string) (id RLID, err error) {
	return id, id.Decode([]byte(rlid), false)
}

// Parse an RLID from an encoded string in strict mode.
// ErrDataSize is returned if the length of the string is different from the expected
// encoded length of the RLID. ErrInvalidCharacters is returned if the encoding is not
// valid or would produce an undefined RLID.
func ParseStrict(rlid string) (id RLID, err error) {
	return id, id.Decode([]byte(rlid), true)
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
	return time.UnixMilli(int64(ms)).In(time.UTC)
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

const (
	// EncodedSize is the length of a text encoded RLID.
	EncodedSize = 16

	// Encoding is the base 32 encoding alphabet used in ULID strings.
	Encoding = "0123456789abcdefghjkmnpqrstvwxyz"
)

// Byte to index table for O(1) lookups when decoding.
// We use 0xFF as sentinel value for invalid indexes.
var dec = [...]byte{
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x01,
	0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E,
	0x0F, 0x10, 0x11, 0xFF, 0x12, 0x13, 0xFF, 0x14, 0x15, 0xFF,
	0x16, 0x17, 0x18, 0x19, 0x1A, 0xFF, 0x1B, 0x1C, 0x1D, 0x1E,
	0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x0A, 0x0B, 0x0C,
	0x0D, 0x0E, 0x0F, 0x10, 0x11, 0xFF, 0x12, 0x13, 0xFF, 0x14,
	0x15, 0xFF, 0x16, 0x17, 0x18, 0x19, 0x1A, 0xFF, 0x1B, 0x1C,
	0x1D, 0x1E, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}

// Encode the RLID as a string to the given buffer.
// ErrBufferSize is returned if the dst is not equal to the encoded size.
func (id RLID) Encode(dst []byte) error {
	if len(dst) != EncodedSize {
		return ErrBufferSize
	}

	// Optimized unrolled loop for performance.
	// Note: the unrolled loop repeats every 8 steps and covers 5 bytes evenly
	dst[0] = Encoding[(id[0]&248)>>3]
	dst[1] = Encoding[((id[0]&7)<<2)|((id[1]&192)>>6)]
	dst[2] = Encoding[(id[1]&62)>>1]
	dst[3] = Encoding[((id[1]&1)<<4)|((id[2]&240)>>4)]
	dst[4] = Encoding[((id[2]&15)<<1)|((id[3]&128)>>7)]
	dst[5] = Encoding[(id[3]&124)>>2]
	dst[6] = Encoding[((id[3]&3)<<3)|((id[4]&224)>>5)]
	dst[7] = Encoding[id[4]&31]
	dst[8] = Encoding[(id[5]&248)>>3]
	dst[9] = Encoding[((id[5]&7)<<2)|((id[6]&192)>>6)]
	dst[10] = Encoding[(id[6]&62)>>1]
	dst[11] = Encoding[((id[6]&1)<<4)|((id[7]&240)>>4)]
	dst[12] = Encoding[((id[7]&15)<<1)|((id[8]&128)>>7)]
	dst[13] = Encoding[(id[8]&124)>>2]
	dst[14] = Encoding[((id[8]&3)<<3)|((id[9]&224)>>5)]
	dst[15] = Encoding[id[9]&31]

	return nil
}

// Decode the RLID from a string represented as a UTF encoded byte array. If strict is
// true then the decoder will validate that the string contains only valid base32
// characters, but this is slightly slower if the input is known to be valid.
func (id *RLID) Decode(src []byte, strict bool) error {
	// Check that the string is the correct length.
	if len(src) != EncodedSize {
		return ErrDataSize
	}

	// Check if all characters are part of the expected character set
	if strict && (dec[src[0]] == 0xFF ||
		dec[src[1]] == 0xFF ||
		dec[src[2]] == 0xFF ||
		dec[src[3]] == 0xFF ||
		dec[src[4]] == 0xFF ||
		dec[src[5]] == 0xFF ||
		dec[src[6]] == 0xFF ||
		dec[src[7]] == 0xFF ||
		dec[src[8]] == 0xFF ||
		dec[src[9]] == 0xFF ||
		dec[src[10]] == 0xFF ||
		dec[src[11]] == 0xFF ||
		dec[src[12]] == 0xFF ||
		dec[src[13]] == 0xFF ||
		dec[src[14]] == 0xFF ||
		dec[src[15]] == 0xFF) {
		return ErrInvalidCharacters
	}

	// Optimized unrolled loop for performance.
	// Note: the unrolled loop repeats every 8 steps and covers 5 bytes evenly
	(*id)[0] = (dec[src[0]] << 3) | dec[src[1]]>>2
	(*id)[1] = (dec[src[1]] << 6) | (dec[src[2]] << 1) | (dec[src[3]] >> 4)
	(*id)[2] = (dec[src[3]] << 4) | (dec[src[4]] >> 1)
	(*id)[3] = (dec[src[4]] << 7) | (dec[src[5]] << 2) | (dec[src[6]] >> 3)
	(*id)[4] = (dec[src[6]] << 5) | dec[src[7]]
	(*id)[5] = (dec[src[8]] << 3) | dec[src[9]]>>2
	(*id)[6] = (dec[src[9]] << 6) | (dec[src[10]] << 1) | (dec[src[11]] >> 4)
	(*id)[7] = (dec[src[11]] << 4) | (dec[src[12]] >> 1)
	(*id)[8] = (dec[src[12]] << 7) | (dec[src[13]] << 2) | (dec[src[14]] >> 3)
	(*id)[9] = (dec[src[14]] << 5) | dec[src[15]]

	return nil
}
