package events_test

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oklog/ulid/v2"
	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/rlid"
	"github.com/rotationalio/ensign/pkg/ensign/store/events"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/suite"
)

type eventsTestSuite struct {
	suite.Suite
	store  *events.Store
	dbPath string
}

func (s *eventsTestSuite) SetupSuite() {
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	assert := s.Assert()

	// Create a temporary test database for the tests
	err := s.OpenStore()
	assert.NoError(err, "could not open test suite database")
}

func (s *eventsTestSuite) TearDownSuite() {
	assert := s.Assert()

	// Close the open store being tested against and remove the temporary directory.
	err := s.CloseStore()
	assert.NoError(err, "could not close events store")
}

func (s *eventsTestSuite) OpenStore() (err error) {
	if s.dbPath, err = os.MkdirTemp("", "ensignevents-*"); err != nil {
		return err
	}

	// Open a store for testing
	conf := config.StorageConfig{
		ReadOnly: false,
		DataPath: s.dbPath,
	}

	if s.store, err = events.Open(conf); err != nil {
		return err
	}
	return nil
}

func (s *eventsTestSuite) CloseStore() (err error) {
	if err = s.store.Close(); err != nil {
		return err
	}
	return os.RemoveAll(s.dbPath)
}

func (s *eventsTestSuite) ResetDatabase() {
	assert := s.Assert()
	assert.NoError(s.CloseStore(), "could not close store during reset")
	assert.NoError(s.OpenStore(), "could not open store during reset")
}

func (s *eventsTestSuite) LoadAllFixtures() (nFixtures uint64, err error) {
	var n uint64
	if n, err = s.LoadEventFixtures(); err != nil {
		return nFixtures, err
	}
	nFixtures += n

	// TODO: load meta-event fixtures
	return nFixtures, nil
}

func (s *eventsTestSuite) LoadEventFixtures() (nFixtures uint64, err error) {
	var events []*api.EventWrapper
	if events, err = mock.EventListFixture("testdata/events.pb.json"); err != nil {
		return 0, err
	}

	for _, event := range events {
		if err = s.store.Insert(event); err != nil {
			return nFixtures, err
		}
		nFixtures += 1

		// Compute the event hash in strict mode
		var hash []byte
		if hash, err = event.HashStrict(); err != nil {
			return nFixtures, err
		}

		var topicID ulid.ULID
		if topicID, err = event.ParseTopicID(); err != nil {
			return nFixtures, err
		}

		var eventID rlid.RLID
		if eventID, err = event.ParseEventID(); err != nil {
			return nFixtures, err
		}

		// Insert the hash into the database
		if err = s.store.Indash(topicID, hash, eventID); err != nil {
			return nFixtures, err
		}
		nFixtures++
	}
	return nFixtures, nil
}

func TestEventsStore(t *testing.T) {
	suite.Run(t, &eventsTestSuite{})
}

const readonlyDataPath = "testdata/readonly.zip"

type readonlyEventsTestSuite struct {
	suite.Suite
	store  *events.Store
	dbPath string
}

func (s *readonlyEventsTestSuite) SetupSuite() {
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	assert := s.Assert()

	// Check if the testdata fixtures need to be generated
	// NOTE: to regenerate the fixtures simply delete testdata/readonly.zip and rerun the tests
	var err error
	if _, err = os.Stat(readonlyDataPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = s.GenerateFixture()
		}
		assert.NoError(err, "could not stat or generate the readonly.ldb fixture")
	}

	// Unzip the readonly database into a temporary directory
	s.dbPath, err = os.MkdirTemp("", "ensignevents-readonly-*")
	assert.NoError(err, "could not create temporary directory for database")

	z, err := zip.OpenReader(readonlyDataPath)
	assert.NoError(err, "could not unzip readonly database")

	extract := func(f *zip.File) (err error) {
		var rc io.ReadCloser
		if rc, err = f.Open(); err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(s.dbPath, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(s.dbPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		// Handle directories
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			return nil
		}

		// Handle files
		os.MkdirAll(filepath.Dir(path), 0755)
		var fp *os.File
		if fp, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
			return err
		}
		defer fp.Close()

		if _, err = io.Copy(fp, rc); err != nil {
			return err
		}
		return nil
	}

	for _, f := range z.File {
		err = extract(f)
		assert.NoError(err, "could not extract zip file")
	}

	// Open a read-only data store to the testdata fixtures
	s.store, err = events.Open(config.StorageConfig{
		ReadOnly: true,
		DataPath: s.dbPath,
	})
	assert.NoError(err, "could not open test suite readonly database")
}

func (s *readonlyEventsTestSuite) TearDownSuite() {
	assert := s.Assert()

	// Close the open store being tested against
	err := s.store.Close()
	assert.NoError(err, "could not close readonly events store")

	// Delete the temporary read only database
	os.RemoveAll(s.dbPath)
}

func (s *readonlyEventsTestSuite) GenerateFixture() (err error) {
	var dir string
	if dir, err = os.MkdirTemp("", "ensignevents-readonly-gen-fixture-*"); err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	var db *events.Store
	if db, err = events.Open(config.StorageConfig{DataPath: dir}); err != nil {
		return err
	}
	defer db.Close()

	// Create events to read in the database
	var events []*api.EventWrapper
	if events, err = mock.EventListFixture("testdata/events.pb.json"); err != nil {
		return err
	}

	for _, event := range events {
		// Insert the event into the database
		if err = db.Insert(event); err != nil {
			return err
		}

		// Compute the event hash in strict mode
		var hash []byte
		if hash, err = event.HashStrict(); err != nil {
			return err
		}

		var topicID ulid.ULID
		if topicID, err = event.ParseTopicID(); err != nil {
			return err
		}

		var eventID rlid.RLID
		if eventID, err = event.ParseEventID(); err != nil {
			return err
		}

		// Insert the hash into the database
		if err = db.Indash(topicID, hash, eventID); err != nil {
			return err
		}
	}

	// TODO: create meta-events to read in the database

	// Create the fixture
	var f *os.File
	if f, err = os.Create(readonlyDataPath); err != nil {
		return err
	}
	defer f.Close()

	// Zip the files in the directory
	z := zip.NewWriter(f)
	defer z.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var file *os.File
		if file, err = os.Open(path); err != nil {
			return err
		}
		defer file.Close()

		var f io.Writer
		if f, err = z.Create(strings.TrimPrefix(path, dir)); err != nil {
			return err
		}

		if _, err = io.Copy(f, file); err != nil {
			return err
		}

		return nil
	}

	if err = filepath.Walk(dir, walker); err != nil {
		return err
	}
	return nil
}

func TestReadOnlyEventsStore(t *testing.T) {
	suite.Run(t, &readonlyEventsTestSuite{})
}
