package buffer_test

import (
	"context"
	"fmt"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/buffer"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
)

func TestRingSingle(t *testing.T) {
	ctx := context.Background()
	event := &api.Event{
		Id:       rlid.Make(42).String(),
		TopicId:  rlid.Make(24).String(),
		Mimetype: mimetype.TextPlain,
		Type: &api.Type{
			Name:    "Test Topic",
			Version: 1,
		},
		Data: []byte("A"),
	}

	// Make a ring with a buffer size of 1
	buf := buffer.NewRing(1)

	// Should not be able to read from an empty buffer
	_, err := buf.Read(ctx)
	require.ErrorIs(t, err, buffer.ErrBufferEmpty)

	// Should be able to read and write from the buffer
	err = buf.Write(ctx, event)
	require.NoError(t, err, "could not write to the buffer")

	// Should not be able to write to a full buffer
	err = buf.Write(ctx, event)
	require.ErrorIs(t, err, buffer.ErrBufferFull)

	other, err := buf.Read(ctx)
	require.NoError(t, err, "could not read from the buffer")
	require.Same(t, event, other, "expected the point to the same value written on the buffer")

	// Should not be able to read from an empty buffer
	_, err = buf.Read(ctx)
	require.ErrorIs(t, err, buffer.ErrBufferEmpty)

	// Should be able to continuously read and write to the buffer
	for i := uint32(0); i < 1e4; i++ {
		e := &api.Event{
			Id:       rlid.Make(i).String(),
			TopicId:  event.TopicId,
			Mimetype: event.Mimetype,
			Type:     event.Type,
			Data:     []byte(fmt.Sprintf("%06X", i)),
		}

		err := buf.Write(ctx, e)
		require.NoError(t, err, "could not write to the buffer")

		other, err := buf.Read(ctx)
		require.NoError(t, err, "could not read from the buffer")
		require.Same(t, e, other, "expected the point to the same value written on the buffer")
		require.NotSame(t, event, other, "expected previous value to be overwritten")
		event = e
	}
}

func TestRingBuffer(t *testing.T) {
	ctx := context.Background()
	event := &api.Event{
		Id:       rlid.Make(42).String(),
		TopicId:  rlid.Make(24).String(),
		Mimetype: mimetype.TextPlain,
		Type: &api.Type{
			Name:    "Test Topic",
			Version: 1,
		},
		Data: []byte("A"),
	}

	// Make a ring with a buffer size of 10
	buf := buffer.NewRing(10)

	// Should not be able to read from an empty buffer
	_, err := buf.Read(ctx)
	require.ErrorIs(t, err, buffer.ErrBufferEmpty)

	// Should be able to write 10 events into the buffer
	for i := uint32(0); i < 10; i++ {
		e := &api.Event{
			Id:       rlid.Make(i).String(),
			TopicId:  event.TopicId,
			Mimetype: event.Mimetype,
			Type:     event.Type,
			Data:     []byte(fmt.Sprintf("%06X", i)),
		}

		err := buf.Write(ctx, e)
		require.NoError(t, err, "could not write to the buffer on iter %d", i)
	}

	// Should not be able to write to a full buffer
	err = buf.Write(ctx, event)
	require.ErrorIs(t, err, buffer.ErrBufferFull)

	// Should be able to read 10 events off the buffer
	for i := uint32(0); i < 10; i++ {
		e, err := buf.Read(ctx)
		require.NoError(t, err, "could not read to the buffer on iter %d", i)
		require.Equal(t, []byte(fmt.Sprintf("%06x", i)), e.Data, "did not read correct event off the buffer on iter %d", i)
	}

	_, err = buf.Read(ctx)
	require.ErrorIs(t, err, buffer.ErrBufferEmpty)
}

func BenchmarkRing(b *testing.B) {
	ctx := context.Background()
	e := &api.Event{
		Id:       rlid.Make(uint32(42)).String(),
		TopicId:  rlid.Make(42).String(),
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:    "Test Topic",
			Version: 1,
		},
		Data: []byte(`{"color": "red", "count": 74}`),
	}

	b.Run("Single", func(b *testing.B) {
		buf := buffer.NewRing(1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Write(ctx, e)
			buf.Read(ctx)
		}
	})

	b.Run("Small", func(b *testing.B) {
		buf := buffer.NewRing(100)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Write(ctx, e)
			buf.Read(ctx)
		}
	})

	b.Run("Medium", func(b *testing.B) {
		buf := buffer.NewRing(1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Write(ctx, e)
			buf.Read(ctx)
		}
	})

	b.Run("Large", func(b *testing.B) {
		buf := buffer.NewRing(1e6)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Write(ctx, e)
			buf.Read(ctx)
		}
	})
}
