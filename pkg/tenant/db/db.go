package db

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/rotationalio/ensign/pkg/tenant/config"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	mu     sync.RWMutex
	cc     *grpc.ClientConn
	client trtl.TrtlClient
)

// Connect to the trtl database, this function must be called at least once before any
// database interaction can occur. Multiple calls to Connect will not error (e.g. if the
// database is already connected then nothing will happen).
func Connect(conf config.DatabaseConfig) (err error) {
	if conf.Testing {
		// TODO: setup mock trtl connection for testing
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	// Check if we're already connected and don't try to reconnect if we are.
	if connected() {
		return nil
	}

	var endpoint string
	if endpoint, err = conf.Endpoint(); err != nil {
		return err
	}

	// Otherwise connect to the trtl database
	opts := make([]grpc.DialOption, 0, 1)
	if conf.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// TODO: connect with mtls
		return errors.New("not implemented: mtls currently not implemented")
	}

	if cc, err = grpc.Dial(endpoint, opts...); err != nil {
		return err
	}

	client = trtl.NewTrtlClient(cc)
	return nil
}

// Close the connection to the database, once closed the package must be reconnected
// otherwise database operations will not succeed. If the database is already closed
// then no error will occur.
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if !connected() {
		return nil
	}

	err := cc.Close()
	cc = nil
	client = nil
	return err
}

// IsConnected returns true if the database has been connected to without error and the
// db module is ready to interact with the trtl database.
func IsConnected() bool {
	mu.RLock()
	defer mu.RUnlock()
	return connected()
}

// Internal check without locks to determine connection state.
func connected() bool {
	return cc == nil || client == nil
}

// Models are structs that have key-value properties that can used for Get, Put, and
// Delete operations to the database. The Model interface allows us to unify common
// interaction patterns (for example checking connections) and returning specific errors
// as well as ensuring that serialization and deserialization occur correctly.
type Model interface {
	// Handle database storage semantics
	Key() ([]byte, error)
	Namespace() string

	// Handle serialization and deserialization of a single Model
	MarshalValue() ([]byte, error)
	UnmarshalValue([]byte) error
}

// Get retrieves a model value based on its key and namespace.
func Get(ctx context.Context, model Model) (err error) {
	// mu is the connection lock, so it ensures that the connection cannot be closed
	// while we're performing this operation. All database calls should have an rlock
	// so that each db call can be concurrent.
	mu.RLock()
	defer mu.RUnlock()

	// Ensure we're connected so that we can do the Get.
	if !connected() {
		return ErrNotConnected
	}

	// Prepare the Get request
	req := &trtl.GetRequest{
		Namespace: model.Namespace(),
	}

	// Compute the key from the model
	if req.Key, err = model.Key(); err != nil {
		return err
	}

	// Execute the Get request
	var rep *trtl.GetReply
	if rep, err = client.Get(ctx, req); err != nil {
		// TODO: transform this error into a more meaningful error
		// E.g. if it's NotFound or Unavailable, etc. return a db.Error instead
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return ErrNotFound
			}
		}
		return err
	}

	// Unmarshal the data into the model
	if err = model.UnmarshalValue(rep.Value); err != nil {
		return err
	}
	return nil
}

func Put(ctx context.Context, model Model) (err error) {
	mu.RLock()
	defer mu.Unlock()

	if !connected() {
		return ErrNotConnected
	}

	req := &trtl.PutRequest{
		Namespace: model.Namespace(),
	}

	if req.Key, err = model.Key(); err != nil {
		return err
	}

	if req.Value, err = model.MarshalValue(); err != nil {
		return err
	}

	if _, err = client.Put(ctx, req); err != nil {
		// TODO: transform this error into a more meaningful error
		// E.g. if it's NotFound or Unavailable, etc. return a db.Error instead
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return ErrNotFound
			}
		}
		return err
	}
	return nil
}

func Delete(ctx context.Context, model Model) (err error) {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return ErrNotConnected
	}

	req := &trtl.DeleteRequest{
		Namespace: model.Namespace(),
	}

	if req.Key, err = model.Key(); err != nil {
		return err
	}

	if _, err = client.Delete(ctx, req); err != nil {
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return ErrNotFound
			}
		}
		return err
	}
	return nil
}

func List(ctx context.Context, prefix []byte, namespace string) (values [][]byte, err error) {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return nil, ErrNotConnected
	}

	req := &trtl.CursorRequest{
		Prefix:    prefix,
		Namespace: namespace,
	}

	// If pagination is required, use Iter instead of Cursor
	var stream trtl.Trtl_CursorClient
	if stream, err = client.Cursor(ctx, req); err != nil {
		return nil, err
	}

	values = make([][]byte, 0)

	// Keep looping over stream until done
	for {
		var item *trtl.KVPair
		if item, err = stream.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		values = append(values, item.Value)
	}

	return values, nil
}
