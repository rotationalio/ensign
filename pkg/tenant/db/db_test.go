package db_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/rotationalio/ensign/pkg/tenant/config"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/logger"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/trisacrypto/directory/pkg/trtl"
	trtlconfig "github.com/trisacrypto/directory/pkg/trtl/config"
	"github.com/trisacrypto/directory/pkg/trtl/mock"
	pb "github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//===========================================================================
// Database Test Suite
//===========================================================================

type dbTestSuite struct {
	suite.Suite
	mock *mock.RemoteTrtl
}

func (s *dbTestSuite) SetupSuite() {
	require := s.Require()

	// Reduce logging clutter for tests
	logger.Discard()

	require.NoError(db.Connect(config.DatabaseConfig{Testing: true}), "unable to connect to db in testing mode")
	require.True(db.IsTesting(), "expected database to be in testing mode")
	require.True(db.IsConnected(), "expected database to be connected in testing mode")

	// Add the mock to the suite for ease of use
	s.mock = db.GetMock()
}

func (s *dbTestSuite) TearDownSuite() {
	require := s.Require()
	require.NoError(db.Close(), "could not close and cleanup connection to test database")
}

func (s *dbTestSuite) AfterTest(suiteName, testName string) {
	s.mock.Reset()
}

func TestDB(t *testing.T) {
	suite.Run(t, new(dbTestSuite))
}

//===========================================================================
// Mock Model
//===========================================================================

// MockModel implements the Model interface and records calls to its methods.
type MockModel struct {
	Calls            map[string]int
	OnKey            func() ([]byte, error)
	OnNamespace      func() string
	OnMarshalValue   func() ([]byte, error)
	OnUnmarshalValue func([]byte) error
}

var _ db.Model = &MockModel{}

func (m *MockModel) Key() ([]byte, error) {
	m.Incr("Key")
	if m.OnKey == nil {
		return []byte("testkey"), nil
	}
	return m.OnKey()
}

func (m *MockModel) Namespace() string {
	m.Incr("Namespace")
	if m.OnNamespace == nil {
		return "testing"
	}
	return m.OnNamespace()
}

func (m *MockModel) MarshalValue() ([]byte, error) {
	m.Incr("MarshalValue")
	if m.OnMarshalValue == nil {
		return []byte("testvalue"), nil
	}
	return m.OnMarshalValue()
}

func (m *MockModel) UnmarshalValue(data []byte) error {
	m.Incr("UnmarshalValue")
	if m.OnUnmarshalValue == nil {
		return nil
	}
	return m.OnUnmarshalValue(data)
}

func (m *MockModel) Incr(name string) {
	if m.Calls == nil {
		m.Calls = make(map[string]int)
	}
	m.Calls[name]++
}

//===========================================================================
// DB Interaction Tests
//===========================================================================

func (s *dbTestSuite) TestGet() {
	require := s.Require()
	model := &MockModel{}
	ctx := context.Background()

	// Test NotFound path
	s.mock.UseError(mock.GetRPC, codes.NotFound, "document not found")
	err := db.Get(ctx, model)
	require.ErrorIs(err, db.ErrNotFound)

	// Test Happy Path and handling
	s.mock.OnGet = func(ctx context.Context, in *pb.GetRequest) (out *pb.GetReply, err error) {
		if in.Namespace != "testing" {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown namespace %q", in.Namespace)
		}

		if !bytes.Equal(in.Key, []byte("testkey")) {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown key %q", string(in.Key))
		}

		return &pb.GetReply{
			Value: []byte("pontoons"),
		}, nil
	}

	// Execute happy path Get request
	err = db.Get(ctx, model)
	require.NoError(err, "could not execute happy path Get request")

	// Ensure the mock model has been called correctly
	require.Equal(2, model.Calls["Namespace"])
	require.Equal(2, model.Calls["Key"])
	require.Equal(1, model.Calls["UnmarshalValue"])
	require.Equal(0, model.Calls["MarshalValue"])
}

func (s *dbTestSuite) TestPut() {
	require := s.Require()
	model := &MockModel{}
	ctx := context.Background()

	// Test NotFound path
	s.mock.UseError(mock.PutRPC, codes.NotFound, "document not found")
	err := db.Put(ctx, model)
	require.ErrorIs(err, db.ErrNotFound)

	// Test Happy Path and handling
	s.mock.OnPut = func(ctx context.Context, in *pb.PutRequest) (out *pb.PutReply, err error) {
		if in.Namespace != "testing" {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown namespace %q", in.Namespace)
		}

		if !bytes.Equal(in.Key, []byte("testkey")) {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown key %q", string(in.Key))
		}

		if !bytes.Equal(in.Value, []byte("testvalue")) {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown value %q", string(in.Value))
		}

		return &pb.PutReply{
			Success: true,
		}, nil
	}

	// Execute happy path Put request
	err = db.Put(ctx, model)
	require.NoError(err, "could not execute happy path Put request")

	// Ensure the mock model has been called correctly
	require.Equal(2, model.Calls["Namespace"])
	require.Equal(2, model.Calls["Key"])
	require.Equal(0, model.Calls["UnmarshalValue"])
	require.Equal(2, model.Calls["MarshalValue"])
}

func (s *dbTestSuite) TestDelete() {
	require := s.Require()
	model := &MockModel{}
	ctx := context.Background()

	// Test NotFound path
	s.mock.UseError(mock.DeleteRPC, codes.NotFound, "document not found")
	err := db.Delete(ctx, model)
	require.ErrorIs(err, db.ErrNotFound)

	// Test Happy Path and handling
	s.mock.OnDelete = func(ctx context.Context, in *pb.DeleteRequest) (out *pb.DeleteReply, err error) {
		if in.Namespace != "testing" {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown namespace %q", in.Namespace)
		}

		if !bytes.Equal(in.Key, []byte("testkey")) {
			return nil, status.Errorf(codes.FailedPrecondition, "unknown key %q", string(in.Key))
		}

		return &pb.DeleteReply{
			Success: true,
		}, nil
	}

	// Execute happy path Delete request
	err = db.Delete(ctx, model)
	require.NoError(err, "could not execute happy path Delete request")

	// Ensure the mock model has been called correctly
	require.Equal(2, model.Calls["Namespace"])
	require.Equal(2, model.Calls["Key"])
	require.Equal(0, model.Calls["UnmarshalValue"])
	require.Equal(0, model.Calls["MarshalValue"])
}

func (s *dbTestSuite) TestList() {
	require := s.Require()
	ctx := context.Background()

	type record struct {
		key   []byte
		value []byte
	}

	namespace := "testing"
	prefix := []byte("test")
	values := make([]record, 0, 101)
	for i := 0; i < 100; i++ {
		key, err := db.CreateKey(ulids.New(), ulids.New())
		require.NoError(err, "could not create key")
		keyData, err := key.MarshalValue()
		require.NoError(err, "could not marshal key")
		values = append(values, record{key: keyData, value: []byte(fmt.Sprintf("value %d", i))})
	}
	lastKey := ulids.New()
	values = append(values, record{key: lastKey[:], value: []byte("last value")})

	// Configure trtl to return the records on cursor
	s.mock.OnCursor = func(in *pb.CursorRequest, stream pb.Trtl_CursorServer) error {
		if !bytes.Equal(in.Prefix, prefix) || in.Namespace != namespace {
			return status.Error(codes.FailedPrecondition, "unexpected cursor request")
		}

		// Send back some data and terminate
		for _, record := range values {
			stream.Send(&pb.KVPair{
				Key:       record.key,
				Value:     record.value,
				Namespace: in.Namespace,
			})
		}

		return nil
	}

	onListItem := func(k *pb.KVPair) error {
		return nil
	}

	s.Run("Single Page", func() {
		// If all the results are on a single page then the next cursor is nil
		cursor := &pg.Cursor{
			PageSize: 101,
		}

		cursor, err := db.List(ctx, prefix, namespace, onListItem, cursor)
		require.NoError(err, "error returned from list request")
		require.Nil(cursor, "cursor should be nil since the results are on one page")
	})

	s.Run("Multiple Pages", func() {
		// If the results are on multiple pages then the next cursor is not nil
		cursor, err := db.List(ctx, prefix, namespace, onListItem, nil)
		require.NoError(err, "error returned from list")
		require.NotNil(cursor, "cursor should not be be nil because there is a next page")
		require.Equal(int32(100), cursor.PageSize, "expected default page size to be 100")

		// Ensure the start index for the next cursor is set correctly
		require.NoError(err, "could not convert first key to string")
		require.Equal(lastKey.String(), cursor.StartIndex, "unexpected start index")
		require.Empty(cursor.EndIndex, "expected end index to be empty")

		// Configure trtl to return the rest of the results
		values = values[100:]
		cursor, err = db.List(ctx, prefix, namespace, onListItem, cursor)
		require.NoError(err, "error returned from list")
		require.Nil(cursor, "cursor should be nil because there is no next page")
	})
}

//===========================================================================
// Connection Tests
//===========================================================================

func TestDBTestingMode(t *testing.T) {
	// Should be able to connect and close when db is in testing mode
	conf := config.DatabaseConfig{Testing: true}
	require.NoError(t, db.Connect(conf), "could not connect to database in testing mode")

	require.True(t, db.IsTesting(), "expected to be in testing mode")
	require.True(t, db.IsConnected(), "expected database to be connected in testing mode")

	require.NoError(t, db.Close(), "could not close database in testing mode")
	require.False(t, db.IsConnected(), "expected database to be not connected after close")
}

func TestDBLiveConnection(t *testing.T) {
	// Should be able to connect to a live trtl database
	// Start by running a live trtl on a tmp directory
	dbdir := t.TempDir()

	tconf, err := trtlconfig.Config{
		Maintenance: false,
		BindAddr:    "127.0.0.1:4436",
		Metrics: trtlconfig.MetricsConfig{
			Enabled: false,
		},
		Database: trtlconfig.DatabaseConfig{
			URL:           "leveldb:////" + dbdir,
			ReindexOnBoot: false,
		},
		Replica: trtlconfig.ReplicaConfig{
			Enabled: false,
			PID:     8,
			Region:  "localhost",
		},
		MTLS: trtlconfig.MTLSConfig{
			Insecure: true,
		},
		Backup: trtlconfig.BackupConfig{
			Enabled: false,
		},
	}.Mark()
	require.NoError(t, err, "trtl config invalid")

	tdb, err := trtl.New(tconf)
	require.NoError(t, err, "could not create a new trtl database server")

	go tdb.Serve()
	defer tdb.Shutdown()

	// Should be able to connect and close when db is not in testing mode
	conf := config.DatabaseConfig{URL: "trtl://127.0.0.1:4436", Insecure: true}
	require.NoError(t, db.Connect(conf), "could not connect to database in testing mode")

	require.False(t, db.IsTesting(), "expected not to be in testing mode")
	require.True(t, db.IsConnected(), "expected database to be connected to live server")

	require.NoError(t, db.Close(), "could not close connection to live database")
	require.False(t, db.IsConnected(), "expected database to be not connected after close")
}
