package buffer_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/rotationalio/ensign/pkg/ensign/buffer"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
)

func BenchmarkBuffer(b *testing.B) {
	// Setup the benchmark by creating an array of events to enqueue
	ctx := context.Background()
	mkevent := func(seq uint32) *api.Event {
		e := &api.Event{
			Id:       rlid.Make(seq).String(),
			TopicId:  rlid.Make(42).String(),
			Mimetype: mimetype.ApplicationJSON,
			Type: &api.Type{
				Name:    "Test Topic",
				Version: 1,
			},
			Data: make([]byte, 64),
		}

		rand.Read(e.Data)
		return e
	}

	events := make([]*api.Event, 0, 128)
	for i := uint32(0); i < 128; i++ {
		events = append(events, mkevent(i+1))
	}

	// Benchmark reading and writing to a buffer continuously
	b.Run("Single", func(b *testing.B) {
		// Test with a buffer size of 1 (worst case for ring)
		cbuf := make(buffer.Channel, 1)
		rbuf := buffer.NewRing(1)
		lbuf := buffer.NewLockingRing(1)

		b.Run("Channel", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cbuf.Write(ctx, events[i%128])
				cbuf.Read(ctx)
			}
		})

		b.Run("Ring", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rbuf.Write(ctx, events[i%128])
				rbuf.Read(ctx)
			}
		})

		b.Run("Locking", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				lbuf.Write(ctx, events[i%128])
				lbuf.Read(ctx)
			}
		})
	})

	b.Run("Alternating", func(b *testing.B) {
		// Test with a large buffer size, alternating between reads and writes
		cbuf := make(buffer.Channel, 128)
		rbuf := buffer.NewRing(128)
		lbuf := buffer.NewLockingRing(128)

		b.Run("Channel", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					cbuf.Write(ctx, events[j])
					cbuf.Read(ctx)
				}
			}
		})

		b.Run("Ring", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					rbuf.Write(ctx, events[j])
					rbuf.Read(ctx)
				}
			}
		})

		b.Run("Locking", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					lbuf.Write(ctx, events[j])
					lbuf.Read(ctx)
				}
			}
		})
	})

	b.Run("FillAndEmpty", func(b *testing.B) {
		// Test with a large buffer size, filling the buffer, then emptying it
		cbuf := make(buffer.Channel, 128)
		rbuf := buffer.NewRing(128)
		lbuf := buffer.NewLockingRing(128)

		b.Run("Channel", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					cbuf.Write(ctx, events[j])
				}

				for j := 0; j < 128; j++ {
					cbuf.Read(ctx)
				}
			}
		})

		b.Run("Ring", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					rbuf.Write(ctx, events[j])
				}

				for j := 0; j < 128; j++ {
					rbuf.Read(ctx)
				}
			}
		})

		b.Run("Locking", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for j := 0; j < 128; j++ {
					lbuf.Write(ctx, events[j])
				}

				for j := 0; j < 128; j++ {
					lbuf.Read(ctx)
				}
			}
		})
	})
}
