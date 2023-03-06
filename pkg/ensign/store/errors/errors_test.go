package errors_test

import (
	"errors"
	"testing"

	. "github.com/rotationalio/ensign/pkg/ensign/store/errors"
	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestErrors(t *testing.T) {
	testCases := []struct {
		Err  error
		LDBE error
	}{
		{ErrReadOnly, leveldb.ErrReadOnly},
		{ErrNotFound, leveldb.ErrNotFound},
		{ErrClosed, leveldb.ErrClosed},
		{ErrIterReleased, leveldb.ErrIterReleased},
		{ErrSnapshotReleased, leveldb.ErrSnapshotReleased},
	}

	for i, tc := range testCases {
		require.ErrorIs(t, tc.Err, tc.LDBE, "test case %d failed", i)
		require.ErrorIs(t, Wrap(tc.LDBE), tc.Err, "test case %d failed", i)
		require.Equal(t, errors.Unwrap(tc.Err), tc.LDBE, "test case %d unwrap failed", i)
		require.NotEqual(t, tc.LDBE.Error(), tc.Err.Error(), "test case %d error failed", i)
	}
}

func TestDefaultWrap(t *testing.T) {
	err := Wrap(errors.New("something bad happened"))

	var target *Error
	require.ErrorAs(t, err, &target)
	require.EqualError(t, err, "unhandled store exception occurred: something bad happened")
}
