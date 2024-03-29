package buffer_test

import (
	"context"
	"sync"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/buffer"
	mimetype "github.com/rotationalio/ensign/pkg/ensign/mimetype/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestChannel(t *testing.T) {
	event := &api.EventWrapper{
		Id:      rlid.Make(42).Bytes(),
		TopicId: rlid.Make(24).Bytes(),
	}
	event.Wrap(&api.Event{
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "Test Topic",
			MajorVersion: 1,
		},
		Data: []byte(`{"color": "red", "count": 74}`),
	})

	// Make a channel with a buffer size of 1
	buf := make(buffer.Channel, 1)

	// Should be able to read and write from the buffer
	err := buf.Write(context.Background(), event)
	require.NoError(t, err, "could not write to the buffer")

	other, err := buf.Read(context.Background())
	require.NoError(t, err, "could not read from the buffer")
	require.Same(t, event, other, "expected the pointer to the same value sent on the channel")
}

func BenchmarkChannelRead(b *testing.B) {
	done := make(chan bool, 2)
	ctx := context.Background()
	native := make(chan *api.EventWrapper, 1024)
	buffer := make(buffer.Channel, 1024)

	event := &api.EventWrapper{
		Id:      rlid.Make(42).Bytes(),
		TopicId: rlid.Make(24).Bytes(),
	}
	event.Wrap(&api.Event{
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "Test Topic",
			MajorVersion: 1,
		},
		Data: []byte(`{"color": "red", "count": 74}`),
	})
	data, _ := proto.Marshal(event)

	var wg sync.WaitGroup
	wg.Add(2)

	b.Cleanup(func() {
		done <- true
		done <- true
		wg.Wait()
		close(native)
		close(buffer)
	})

	go func(done <-chan bool) {
		defer wg.Done()
		var seq uint32
		for {
			seq++
			e := &api.EventWrapper{
				Id:      rlid.Make(seq).Bytes(),
				TopicId: event.TopicId,
			}
			e.Wrap(&api.Event{
				Mimetype: mimetype.ApplicationJSON,
				Type: &api.Type{
					Name:         "Test Topic",
					MajorVersion: 1,
				},
				Data: []byte(`{"color": "red", "count": 74}`),
			})

			select {
			case <-done:
				return
			case native <- e:
			}

		}
	}(done)

	go func(done <-chan bool) {
		defer wg.Done()
		var seq uint32
		for {
			seq++
			e := &api.EventWrapper{
				Id:      rlid.Make(seq).Bytes(),
				TopicId: event.TopicId,
			}
			e.Wrap(&api.Event{
				Mimetype: mimetype.ApplicationJSON,
				Type: &api.Type{
					Name:         "Test Topic",
					MajorVersion: 1,
				},
				Data: []byte(`{"color": "red", "count": 74}`),
			})

			select {
			case <-done:
				return
			case buffer <- e:
			}
		}
	}(done)

	b.Run("Buffer", func(b *testing.B) {
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buffer.Read(ctx)
		}
	})

	b.Run("Native", func(b *testing.B) {
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			<-native
		}
	})
}

func BenchmarkChannelWrite(b *testing.B) {
	ctx := context.Background()
	native := make(chan *api.EventWrapper, 1024)
	buffer := make(buffer.Channel, 1024)

	event := &api.EventWrapper{
		Id:      rlid.Make(0).Bytes(),
		TopicId: rlid.Make(0).Bytes(),
	}
	event.Wrap(&api.Event{
		Mimetype: mimetype.ApplicationJSON,
		Type: &api.Type{
			Name:         "Test Topic",
			MajorVersion: 1,
		},
		Data: []byte(`{"color": "red", "count": 74}`),
	})
	data, _ := proto.Marshal(event)

	b.Cleanup(func() {
		close(native)
		close(buffer)
	})

	go func() {
		for range native {
		}
	}()

	go func() {
		for {
			buffer.Read(ctx)
		}
	}()

	b.Run("Buffer", func(b *testing.B) {
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buffer.Write(ctx, event)
		}
	})

	b.Run("Native", func(b *testing.B) {
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			native <- event
		}
	})
}
