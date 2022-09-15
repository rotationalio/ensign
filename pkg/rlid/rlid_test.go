package rlid_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/rotationalio/ensign/pkg/rlid"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
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

func TestSequence(t *testing.T) {
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
