package meta_test

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	api "github.com/rotationalio/ensign/pkg/ensign/api/v1beta1"
	"github.com/rotationalio/ensign/pkg/ensign/config"
	"github.com/rotationalio/ensign/pkg/ensign/store/meta"
	"github.com/rotationalio/ensign/pkg/ensign/store/mock"
	"github.com/stretchr/testify/suite"
)

const (
	topicsFixturePath     = "testdata/topics.json"
	topicInfosFixturePath = "testdata/topic_infos.json"
	groupsFixturePath     = "testdata/groups.json"
)

type metaTestSuite struct {
	suite.Suite
	store  *meta.Store
	dbPath string
}

func (s *metaTestSuite) SetupSuite() {
	// Note use assert instead of require so that go routines are properly handled in
	// tests; assert uses t.Error while require uses t.FailNow and multiple go routines
	// might lead to incorrect testing behavior.
	assert := s.Assert()

	// Create a temporary test database for the tests
	err := s.OpenStore()
	assert.NoError(err, "could not open test suite database")
}

func (s *metaTestSuite) TearDownSuite() {
	assert := s.Assert()

	// Close the open store being tested against and remove the temporary directory.
	err := s.CloseStore()
	assert.NoError(err, "could not close meta store")
}

func (s *metaTestSuite) OpenStore() (err error) {
	if s.dbPath, err = os.MkdirTemp("", "ensignmeta-*"); err != nil {
		return err
	}

	// Open a store for testing
	conf := config.StorageConfig{
		ReadOnly: false,
		DataPath: s.dbPath,
	}

	if s.store, err = meta.Open(conf); err != nil {
		return err
	}
	return nil
}

func (s *metaTestSuite) CloseStore() (err error) {
	if err = s.store.Close(); err != nil {
		return err
	}
	return os.RemoveAll(s.dbPath)
}

func (s *metaTestSuite) ResetDatabase() {
	assert := s.Assert()
	assert.NoError(s.CloseStore(), "could not close store during reset")
	assert.NoError(s.OpenStore(), "could not open store during reset")
}

func (s *metaTestSuite) LoadAllFixtures() (nFixtures uint64, err error) {
	var n uint64
	if n, err = s.LoadTopicFixtures(); err != nil {
		return nFixtures, err
	}
	nFixtures += n

	if n, err = s.LoadTopicInfoFixtures(); err != nil {
		return nFixtures, err
	}
	nFixtures += n

	if n, err = s.LoadGroupFixtures(); err != nil {
		return nFixtures, err
	}
	nFixtures += n

	return nFixtures, nil
}

func (s *metaTestSuite) LoadTopicFixtures() (nFixtures uint64, err error) {
	var topics []*api.Topic
	if topics, err = mock.TopicListFixture(topicsFixturePath); err != nil {
		return 0, err
	}

	for _, topic := range topics {
		if err = s.store.CreateTopic(topic); err != nil {
			return nFixtures, err
		}
		nFixtures += 3
	}
	return nFixtures, nil
}

func (s *metaTestSuite) LoadTopicInfoFixtures() (nFixtures uint64, err error) {
	var infos map[string]*api.TopicInfo
	if infos, err = mock.TopicInfoListFixture(topicInfosFixturePath); err != nil {
		return 0, err
	}

	for _, info := range infos {
		if err = s.store.UpdateTopicInfo(info); err != nil {
			return nFixtures, err
		}
		nFixtures++
	}

	return nFixtures, nil
}

func (s *metaTestSuite) LoadGroupFixtures() (nFixtures uint64, err error) {
	var groups []*api.ConsumerGroup
	if groups, err = mock.GroupListFixture(groupsFixturePath); err != nil {
		return 0, err
	}

	for _, group := range groups {
		if err = s.store.CreateGroup(group); err != nil {
			return nFixtures, err
		}
		nFixtures += 1
	}
	return nFixtures, nil
}

func TestMetaStore(t *testing.T) {
	suite.Run(t, &metaTestSuite{})
}

const readonlyDataPath = "testdata/readonly.zip"

type readonlyMetaTestSuite struct {
	suite.Suite
	store  *meta.Store
	dbPath string
}

func (s *readonlyMetaTestSuite) SetupSuite() {
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
	s.dbPath, err = os.MkdirTemp("", "ensignmeta-readonly-*")
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
	s.store, err = meta.Open(config.StorageConfig{
		ReadOnly: true,
		DataPath: s.dbPath,
	})
	assert.NoError(err, "could not open test suite readonly database")
}

func (s *readonlyMetaTestSuite) TearDownSuite() {
	assert := s.Assert()

	// Close the open store being tested against
	err := s.store.Close()
	assert.NoError(err, "could not close readonly meta store")

	// Delete the temporary read only database
	os.RemoveAll(s.dbPath)
}

func (s *readonlyMetaTestSuite) GenerateFixture() (err error) {
	var dir string
	if dir, err = os.MkdirTemp("", "ensignmeta-readonly-gen-fixture-*"); err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	var db *meta.Store
	if db, err = meta.Open(config.StorageConfig{DataPath: dir}); err != nil {
		return err
	}
	defer db.Close()

	// Create topics to read in the database
	var topics []*api.Topic
	if topics, err = mock.TopicListFixture(topicsFixturePath); err != nil {
		return err
	}

	for _, topic := range topics {
		if err = db.CreateTopic(topic); err != nil {
			return err
		}
	}

	// Create topic infos to read in to the database
	var infos map[string]*api.TopicInfo
	if infos, err = mock.TopicInfoListFixture(topicInfosFixturePath); err != nil {
		return err
	}

	for _, info := range infos {
		if err = db.UpdateTopicInfo(info); err != nil {
			return err
		}
	}

	// Create groups to read in the database
	var groups []*api.ConsumerGroup
	if groups, err = mock.GroupListFixture(groupsFixturePath); err != nil {
		return err
	}

	for _, group := range groups {
		if err = db.CreateGroup(group); err != nil {
			return err
		}
	}

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

func TestReadOnlyMetaStore(t *testing.T) {
	suite.Run(t, &readonlyMetaTestSuite{})
}
