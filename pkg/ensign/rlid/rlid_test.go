package rlid_test

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	. "github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()

	// Basic parsing tests, for more advanced decoding tests see TestEncoding
	id := RLID{0x01, 0x83, 0x42, 0x5F, 0x66, 0x6F, 0x00, 0x6F, 0xEB, 0x6B}

	sid, err := Parse("061m4qv6dw06ztvb")
	require.NoError(t, err, "could not parse valid rlid")
	require.Equal(t, id, sid, "parsed rlid did not match fixture value")

	sid, err = ParseStrict("061m4qv6dw06ztvb")
	require.NoError(t, err, "could not parse valid rlid in strict mode")
	require.Equal(t, id, sid, "strictly parsed rlid did not match fixture value")

	// Parse should return invalid length but not invalid chars
	_, err = Parse("foo")
	require.ErrorIs(t, err, ErrDataSize, "should not be able to parse short rlid")
	_, err = Parse("061m4?v6!w0>ztvb")
	require.NoError(t, err, "parse should not perform strict validation")

	// ParseStrict should return an error for both invalid length and invalid chars
	_, err = ParseStrict("foo")
	require.ErrorIs(t, err, ErrDataSize, "should not be able to strict parse short rlid")
	_, err = ParseStrict("061m4?v6!w0>ztvb")
	require.ErrorIs(t, err, ErrInvalidCharacters, "strict parse should perform validation")

	// MustParse should not panic on a valid string
	require.NotPanics(t, func() { MustParse("061m4qv6dw06ztvb") })
	require.NotPanics(t, func() { MustParseStrict("061m4qv6dw06ztvb") })

	// MustParse should panic on an invalid string
	require.Panics(t, func() { MustParse("foo") })
	require.Panics(t, func() { MustParseStrict("foo") })

}

func TestBasic(t *testing.T) {
	t.Parallel()

	id := RLID{0x08, 0x12, 0xF5, 0x59, 0x12, 0xA2, 0x3B, 0xFE, 0x01, 0x98}
	require.Equal(t, "109fap8jm8xzw0cr", id.String(), "incorrect stringification of id")
	require.Equal(t, id[:], id.Bytes(), "incorect byte sliceification of id")
	require.Equal(t, 0, id.Compare(id), "id should be equal with itself")

	require.Equal(t, 1, id.Compare(RLID{0x07, 0x11, 0xF4, 0x58, 0x11, 0xA1, 0x3A, 0xFD, 0x00, 0x97}), "could not compare id to lower id")
	require.Equal(t, 0, id.Compare(RLID{0x08, 0x12, 0xF5, 0x59, 0x12, 0xA2, 0x3B, 0xFE, 0x01, 0x98}), "could not compare id to equal id")
	require.Equal(t, -1, id.Compare(RLID{0x18, 0x22, 0xF6, 0x69, 0x22, 0xB2, 0x4B, 0xFF, 0x11, 0xF0}), "could not compare id to greater id")
}

func TestTime(t *testing.T) {
	t.Parallel()

	// Ensure that maxTime is correct with respect to the encoding scheme.
	maxTime := RLID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}.Time()
	require.Equal(t, maxTime, MaxTime, "expected the max time constant to equal the max byte encoding")

	// Create a timestapm fixture for use in following tests
	ts, err := time.Parse(time.RFC3339Nano, "2022-09-15T08:05:42.562000-05:00")
	require.NoError(t, err, "could not parse timestamp fixture")

	// Ensure time.Time to Unix millisecond time conversions work correctly
	ms := Timestamp(ts)
	require.Equal(t, uint64(1663247142562), ms, "timestamp was not correctly generated")
	require.True(t, ts.Equal(Time(ms)), "the timestamp should be correctly converted back to a time.Time")

	// Now should be betweeen when the tests were created and MaxTime
	// NOTE: if this is the year 10889, I apologize for this test failing
	now := Now()
	require.Greater(t, now, ms, "expected unix milliseconds now to be between timestamp tests were created and max time")
	require.Less(t, ms, MaxTime, "expected unix milliseconds now to be between timestamp tests were created and max time")

	// Should be able to set and retrieve time from an RLID
	id := &RLID{}
	require.NoError(t, id.SetTime(ms), "could not set RLID timestamp")
	require.Equal(t, ms, id.Time(), "RLID timestamp was not set properly")
	require.True(t, ts.Equal(Time(id.Time())), "RLID time.Time was not set properly")

	// Should be able to reset a time on an RLID
	require.NoError(t, id.SetTime(now), "could not reset RLID timestamp")
	require.Equal(t, now, id.Time(), "RLID timestamp was not set properly")

	// Should not be able to set a time greater than max time
	err = (&RLID{}).SetTime(maxTime + 2000)
	require.ErrorIs(t, err, ErrOverTime, "should not be able to set a time greater than max time")
}

func TestRLIDSequence(t *testing.T) {
	t.Parallel()

	// Generate several random pairs of sequence numbers for testing
	pairs := []struct {
		seq uint32
		ovr uint32
	}{
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
		{rand.Uint32(), rand.Uint32()},
	}

	for _, pair := range pairs {
		require.NotEqual(t, pair.seq, pair.ovr, "pair should have two different numbers, try rerunning the test")

		// Should be able to set and retrieve a sequence from an RLID
		id := &RLID{}
		require.NoError(t, id.SetSequence(pair.seq), "expected no error setting sequence")
		require.Equal(t, pair.seq, id.Sequence(), "expected sequence to be retrieved from rlid correctly")

		// Should be able to reset a sequence on an RLID
		require.NoError(t, id.SetSequence(pair.ovr), "expected no error overwriting sequence")
		require.Equal(t, pair.ovr, id.Sequence(), "sequence number was not correctly overwritten")
	}
}

func TestEncoding(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	t.Parallel()

	// Ensure the encoding set has 32 characters and that they are all unique
	// NOTE: this test is due to a bug where I accidentally made the original encoding
	// character set 0123456790abcdefghjkmnpqrstvwxyz -- missing 8 and duplicating 0.
	cset := make(map[rune]struct{})
	for _, c := range Encoding {
		cset[c] = struct{}{}
	}
	require.Len(t, cset, 32, "character set contains duplicates or is not 32 characters long")

	// Encode and decode a bunch of different RLIDs with different timestamps
	for i := uint32(0); i < 131072; i++ {
		if i != 0 && i%32768 == 0 {
			time.Sleep(time.Millisecond)
		}

		id := Make(i)
		buf := make([]byte, 16)
		require.NoError(t, id.Encode(buf), "could not encode seq %d", i)

		var jd RLID
		require.NoError(t, jd.Decode(buf, true), "could not decode seq %d", i)

		require.Equal(t, id, jd, "decoded id does not match original for seq %d", i)
		require.NotSame(t, id, jd, "decoding should not clone the pointer for seq %d", i)
	}

	// Test bad buffers for encoding
	id := Make(738123)
	for _, buf := range [][]byte{nil, make([]byte, 8), make([]byte, 64)} {
		require.ErrorIs(t, id.Encode(buf), ErrBufferSize, "expected invalid buffer size")
	}

	// Test bad buffers for decoding
	for _, s := range []string{"", "tooshort", "waytoolongforanrlidevenifitisexpanded"} {
		var id RLID
		require.ErrorIs(t, id.Decode([]byte(s), false), ErrDataSize)
	}

	// Test invalid characters are ignored when not in strict mode
	require.NoError(t, id.Decode([]byte("+!@#$%^&*()_-?><"), false), "expected no character validation when not in strict mode")

	// Test invalid character checks when in strict mode
	for i := 0; i < EncodedSize; i++ {
		buf := []byte(strings.Repeat("8", EncodedSize)) // expects that "8" is a valid character
		buf[i] = byte('?')                              // expects that "?" is an invalid character

		var id RLID
		require.ErrorIs(t, id.Decode(buf, true), ErrInvalidCharacters, "expected character validation when in strict mode")
		require.NoError(t, id.Decode(buf, false), "expected no characer validation when not in strict mode")
	}
}

func TestLexicographicalOrderLowToHigh(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	t.Parallel()

	ms := uint64(0)
	seq := uint32(0)
	orig := RLID{}
	comp := RLID{}

	for i := 0; i < 1e3; i++ {
		ms += uint64(rand.Int63n(1e3)) + 1

		// Test lexicographic ordering with same timestamp
		orig.SetTime(ms)
		comp.SetTime(ms)

		seq = 0
		orig.SetSequence(seq)

		for j := 0; j < 1e3; j++ {
			seq += uint32(rand.Int31n(1e3)) + 1

			comp.SetSequence(seq)
			require.Equal(t, 1, comp.Compare(orig), "comp should be greater than orig on %d-%d", i, j)
			require.Greater(t, comp.String(), orig.String(), "comp string should be lexicographically greater than orig string on %d-%d", i, j)
			orig.SetSequence(seq)
		}

		ms2 := ms
		seq = 0

		// Test lexicographic ordering with a different timestamp and sequence
		for j := 0; j < 1e3; j++ {
			ms2 += uint64(rand.Int63n(1e3)) + 1
			seq += uint32(rand.Int31n(1e3))

			comp.SetTime(ms2)
			comp.SetSequence(seq)

			require.Equal(t, 1, comp.Compare(orig), "comp should be greater than orig on %d-%d", i, j)
			require.Greater(t, comp.String(), orig.String(), "comp string should be lexicographically greater than orig string on %d-%d", i, j)
			orig.SetTime(ms2)
			orig.SetSequence(seq)
		}
	}
}

func TestLexicographicalOrderHighToLow(t *testing.T) {
	// This is a long running test, skip if in short mode
	if testing.Short() {
		t.Skip("skipping long running test in short mode")
	}

	t.Parallel()

	ms := MaxTime
	seq := uint32(math.MaxUint32)

	orig, comp := RLID{}, RLID{}
	orig.SetTime(ms)
	orig.SetSequence(seq)
	comp.SetTime(ms)
	comp.SetSequence(seq - 1)

	require.Equal(t, -1, comp.Compare(orig), "comp should be less than orig")
	require.Less(t, comp.String(), orig.String(), "comp string should be lexicographically less than orig string")

	for i := 0; i < 1e3; i++ {
		ms -= uint64(rand.Int63n(1e3)) + 1

		// Test lexicographic ordering with same timestamp
		orig.SetTime(ms)
		comp.SetTime(ms)

		seq = uint32(math.MaxUint32)
		orig.SetSequence(seq)

		for j := 0; j < 1e3; j++ {
			seq -= uint32(rand.Int31n(1e3)) + 1

			comp.SetSequence(seq)
			require.Equal(t, -1, comp.Compare(orig), "comp should be less than orig on %d-%d", i, j)
			require.Less(t, comp.String(), orig.String(), "comp string should be lexicographically less than orig string on %d-%d", i, j)
			orig.SetSequence(seq)
		}

		ms2 := ms
		seq = uint32(math.MaxUint32)

		// Test lexicographic ordering with a different timestamp and sequence
		for j := 0; j < 1e3; j++ {
			ms2 -= uint64(rand.Int63n(1e3)) + 1
			seq -= uint32(rand.Int31n(1e3))

			comp.SetTime(ms2)
			comp.SetSequence(seq)

			require.Equal(t, -1, comp.Compare(orig), "comp should be less than orig on %d-%d", i, j)
			require.Less(t, comp.String(), orig.String(), "comp string should be lexicographically less than orig string on %d-%d", i, j)
			orig.SetTime(ms2)
			orig.SetSequence(seq)
		}
	}
}

func TestCaseInsensitivity(t *testing.T) {
	var seq uint32
	for i := 0; i < 100; i++ {
		seq += uint32(rand.Int31n(1e3)) + 1
		orig := Make(seq)

		upper := strings.ToUpper(orig.String())
		lower := strings.ToLower(orig.String())

		uid, err := Parse(upper)
		require.NoError(t, err, "could not parse upper case string")
		lid, err := Parse(lower)
		require.NoError(t, err, "Could not parse lower case string")
		require.Equal(t, 0, uid.Compare(lid), "uid and lid should be equal")
	}
}

func TestInterfaces(t *testing.T) {
	id := Make(rand.Uint32())
	cmp := RLID{}

	// Marshal Binary
	data, err := id.MarshalBinary()
	require.NoError(t, err, "could not marshal id as binary data")
	err = cmp.UnmarshalBinary(data)
	require.NoError(t, err, "could not unmarshal from binary data")
	require.Equal(t, 0, id.Compare(cmp), "unmarshaled id should equal marshaled id")

	// Marshal Test
	cmp = RLID{}
	sdata, err := id.MarshalText()
	require.NoError(t, err, "could not marshal id as binary data")
	err = cmp.UnmarshalText(sdata)
	require.NoError(t, err, "could not unmarshal from binary data")
	require.Equal(t, 0, id.Compare(cmp), "unmarshaled id should equal marshaled id")

	// Cannot unmarshal bad data
	require.ErrorIs(t, cmp.UnmarshalBinary([]byte{0x14, 0x22}), ErrDataSize)

	// Cannot decode bad text
	require.ErrorIs(t, cmp.UnmarshalText([]byte("f11")), ErrDataSize)

	// Should be able to Scan a value
	cmp = RLID{}
	require.NoError(t, cmp.Scan(nil), "should be able to scan nil")
	require.NoError(t, cmp.Scan(data), "should be able to scan a byte slice")
	require.Equal(t, 0, id.Compare(cmp), "scanned id should equal marshaled id")

	cmp = RLID{}
	require.NoError(t, cmp.Scan(id.String()), "should be able to scan a string")
	require.Equal(t, 0, id.Compare(cmp), "scanned id should equal marshaled id")

	// Should return an error when scanning a byte string
	require.ErrorIs(t, cmp.Scan(sdata), ErrDataSize, "should not be able to scan a byte string")

	// Should not be able to scan an arbitrary value
	require.ErrorIs(t, cmp.Scan(uint32(44421)), ErrScanValue, "should not be able to scan a uint32")

	// Should be able to return a driver.VAlue as a byte slice.
	val, err := id.Value()
	require.NoError(t, err, "could not extract value from byte slice")
	require.Equal(t, data, val, "expected value to be a binary representation")
}

func BenchmarkMake(b *testing.B) {
	// Benchmark the performance of creating RLIDs using the default method.
	b.ReportAllocs()
	b.SetBytes(int64(len(RLID{})))
	b.ResetTimer()
	var num = rand.Uint32()

	for i := 0; i < b.N; i++ {
		Make(num)
	}
}

func BenchmarkNow(b *testing.B) {
	b.SetBytes(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Now()
	}
}

func BenchmarkTimestamp(b *testing.B) {
	now := time.Now()
	b.SetBytes(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Timestamp(now)
	}
}

func BenchmarkTime(b *testing.B) {
	var num = rand.Uint32()
	id := Make(num)
	b.SetBytes(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = id.Time()
	}
}

func BenchmarkSetTime(b *testing.B) {
	var id RLID
	now := Now()
	b.SetBytes(8)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = id.SetTime(now)
	}
}

func BenchmarkSequence(b *testing.B) {
	var num = rand.Uint32()
	id := Make(num)
	b.SetBytes(4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = id.Sequence()
	}
}

func BenchmarkSetSequence(b *testing.B) {
	var id RLID
	var num = rand.Uint32()
	b.SetBytes(4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = id.SetSequence(num)
	}
}

func BenchmarkEncode(b *testing.B) {
	buf := make([]byte, EncodedSize)
	id := Make(32192212)

	b.SetBytes(int64(len(id)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = id.Encode(buf)
	}
}

func BenchmarkDecode(b *testing.B) {
	var id RLID
	buf := []byte("061m4qv6dw06ztvb")

	b.Run("Quick", func(b *testing.B) {
		b.SetBytes(int64(len(id)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = id.Decode(buf, false)
		}
	})

	b.Run("Strict", func(b *testing.B) {
		b.SetBytes(int64(len(id)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = id.Decode(buf, true)
		}
	})
}
