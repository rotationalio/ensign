package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/utils/mtls"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/trisacrypto/directory/pkg/trtl/mock"
	trtl "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	mu      sync.RWMutex
	cc      *grpc.ClientConn
	client  trtl.TrtlClient
	mockdb  *mock.RemoteTrtl
	testing bool
)

type OnListItem func(*trtl.KVPair) error

// Connect to the trtl database, this function must be called at least once before any
// database interaction can occur. Multiple calls to Connect will not error (e.g. if the
// database is already connected then nothing will happen).
func Connect(conf config.DatabaseConfig) (err error) {
	mu.Lock()
	defer mu.Unlock()

	// Check if we're already connected and don't try to reconnect if we are.
	if connected() {
		return nil
	}

	// Setup a mock remote trtl for in-memory testing of trtl interactions.
	if conf.Testing {
		// Create the mock database and connect to the bufconn client
		mockdb = mock.New(nil)
		if client, err = mockdb.DBClient(); err != nil {
			// Set mock connection to nil to ensure that we can retry the connection.
			client = nil
			mockdb = nil
			return fmt.Errorf("could not connect to mock remote trtl bufconn: %w", err)
		}

		testing = true
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
		// If we're in secure mode, we expect to have certificates
		// We expect that the configuration has been validated prior to this point
		var certs *mtls.Provider
		if certs, err = mtls.Load(conf.CertPath); err != nil {
			return err
		}

		// Load the trusted pool from the provider if it has been specified.
		var trusted []*mtls.Provider
		if conf.PoolPath != "" {
			var trust *mtls.Provider
			if trust, err = mtls.Load(conf.PoolPath); err != nil {
				return err
			}
			trusted = append(trusted, trust)
		}

		// Create client credentials
		var creds grpc.DialOption
		if creds, err = mtls.ClientCreds(endpoint, certs, trusted...); err != nil {
			return err
		}
		opts = append(opts, creds)
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

	if testing {
		mockdb.CloseClient()
		mockdb.Shutdown()
		mockdb = nil
		client = nil
		testing = false
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

func IsTesting() bool {
	mu.RLock()
	defer mu.RUnlock()
	return testing
}

// Internal check without locks to determine connection state.
func connected() bool {
	if testing {
		return mockdb != nil && client != nil
	}
	return cc != nil && client != nil
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
	// Compute the key from the model
	var key []byte
	if key, err = model.Key(); err != nil {
		return err
	}

	// Execute the Get request
	var value []byte
	if value, err = getRequest(ctx, model.Namespace(), key); err != nil {
		return err
	}

	// Unmarshal the data into the model
	if err = model.UnmarshalValue(value); err != nil {
		return err
	}
	return nil
}

func getRequest(ctx context.Context, namespace string, key []byte) (value []byte, err error) {
	// mu is the connection lock, so it ensures that the connection cannot be closed
	// while we're performing this operation. All database calls should have an rlock
	// so that each db call can be concurrent.
	mu.RLock()
	defer mu.RUnlock()

	// Ensure we're connected so that we can do the Get.
	if !connected() {
		return nil, ErrNotConnected
	}

	req := &trtl.GetRequest{
		Namespace: namespace,
		Key:       key,
	}

	// Execute the Get request
	var rep *trtl.GetReply
	if rep, err = client.Get(ctx, req); err != nil {
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return nil, ErrNotFound
			case codes.Unavailable:
				return nil, ErrUnavailable
			}
		}
		return nil, err
	}

	return rep.Value, nil
}

func Put(ctx context.Context, model Model) (err error) {
	var key, value []byte
	if key, err = model.Key(); err != nil {
		return err
	}

	if value, err = model.MarshalValue(); err != nil {
		return err
	}

	if err = putRequest(ctx, model.Namespace(), key, value); err != nil {
		return err
	}

	return nil
}

func putRequest(ctx context.Context, namespace string, key []byte, value []byte) (err error) {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return ErrNotConnected
	}

	req := &trtl.PutRequest{
		Namespace: namespace,
		Key:       key,
		Value:     value,
	}

	if _, err = client.Put(ctx, req); err != nil {
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return ErrNotFound
			case codes.Unavailable:
				return ErrUnavailable
			}
		}
		return err
	}

	return nil
}

func Delete(ctx context.Context, model Model) (err error) {
	var key []byte
	if key, err = model.Key(); err != nil {
		return err
	}

	if err = deleteRequest(ctx, model.Namespace(), key); err != nil {
		return err
	}

	return nil
}

func deleteRequest(ctx context.Context, namespace string, key []byte) (err error) {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return ErrNotConnected
	}

	req := &trtl.DeleteRequest{
		Namespace: namespace,
		Key:       key,
	}

	if _, err = client.Delete(ctx, req); err != nil {
		if serr, ok := status.FromError(err); ok {
			switch serr.Code() {
			case codes.NotFound:
				return ErrNotFound
			case codes.Unavailable:
				return ErrUnavailable
			}
		}
		return err
	}
	return nil
}

// List retrieves a pagination cursor.
func List(ctx context.Context, prefix, seekKey []byte, namespace string, onListItem OnListItem, c *pg.Cursor) (cursor *pg.Cursor, err error) {
	mu.RLock()
	defer mu.RUnlock()

	if !connected() {
		return nil, ErrNotConnected
	}

	// Check to see if a default cursor exists and create one if it does not.
	if c == nil {
		c = pg.New("", "", 0)
	}

	// Set a default page size if one does not exist.
	if c.PageSize <= 0 {
		c.PageSize = 100
	}

	req := &trtl.CursorRequest{
		Prefix:    prefix,
		SeekKey:   seekKey,
		Namespace: namespace,
	}

	var stream trtl.Trtl_CursorClient
	if stream, err = client.Cursor(ctx, req); err != nil {
		return nil, err
	}

	// Keep looping over stream until done
	var startKey, endKey []byte
	nItems := int32(0)
	for {
		var item *trtl.KVPair
		if item, err = stream.Recv(); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		endKey = item.Key
		if startKey == nil {
			startKey = item.Key
		}
		// Check if the number of items is greater than the page size.
		nItems++
		if nItems > c.PageSize {
			break
		}
		if err = onListItem(item); err != nil {
			if errors.Is(err, ErrListBreak) {
				return nil, nil
			}
			return nil, err
		}
	}

	if startKey != nil && nItems > c.PageSize {
		var startID, endID ulid.ULID
		if err = startID.UnmarshalBinary(startKey); err != nil {
			return nil, err
		}
		if err = endID.UnmarshalBinary(endKey); err != nil {
			return nil, err
		}
		cursor = pg.New(startID.String(), endID.String(), c.PageSize)
	}
	return cursor, nil
}

func GetMock() *mock.RemoteTrtl {
	return mockdb
}
