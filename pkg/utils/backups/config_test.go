package backups_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/backups"
	"github.com/stretchr/testify/require"
)

func TestArchiveName(t *testing.T) {
	conf := &backups.Config{}
	require.Regexp(t, `backup-\d{12}\.tgz`, conf.ArchiveName())

	conf.Prefix = "ensign"
	require.Regexp(t, `ensign-\d{12}\.tgz`, conf.ArchiveName())
}

func TestConfigStorage(t *testing.T) {
	testCases := []struct {
		dsn      string
		err      error
		expected interface{}
	}{
		{
			"file:///testdata", nil, &backups.FileStorage{},
		},
	}

	for i, tc := range testCases {
		conf := backups.Config{StorageDSN: tc.dsn}
		storage, err := conf.Storage()

		if tc.err != nil {
			require.ErrorIs(t, err, tc.err, "expected error from storage configuration on test case %d", i)
		} else {
			require.NoError(t, err, "expected no error on test case %d", i)
			require.IsType(t, tc.expected, storage, "expected correct type for storage object on test case %d", i)
		}
	}

}
