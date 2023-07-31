package models_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	. "github.com/rotationalio/ensign/pkg/quarterdeck/db/models"
)

func (m *modelTestSuite) TestIdentifySubject() {
	require := m.Require()

	testCases := []struct {
		subID    string
		orgID    string
		expected SubjectType
		err      error
	}{
		{"01GQFQ4475V3BZDMSXFV5DK6XX", "01GQFQ14HXF2VC7C1HJECS60XX", UserSubject, nil},
		{"01GZ458FBCST1AC1M4PY4X7QZZ", "01GKHJRF01YXHZ51YMMKV3RCMK", APIKeySubject, nil},
		{"01GZ458FBCST1AC1M4PY4X7QZZ", "01GQFQ14HXF2VC7C1HJECS60XX", UnknownSubject, ErrNotFound},
		{"01GQFQ4475V3BZDMSXFV5DK6XX", "01GQZAC80RAZ1XQJKRZJ2R4KNJ", UnknownSubject, ErrNotFound},
		{"01H6PZZ9WMEAWR3YZDDPFXWKBZ", "01GQZAC80RAZ1XQJKRZJ2R4KNJ", UnknownSubject, ErrNotFound},
		{"", "01GQFQ14HXF2VC7C1HJECS60XX", UnknownSubject, ErrNotFound},
		{"01GZ458FBCST1AC1M4PY4X7QZZ", "", UnknownSubject, ErrNotFound},
		{"notaulid", "01GQFQ14HXF2VC7C1HJECS60XX", UnknownSubject, ulid.ErrDataSize},
		{"01GZ458FBCST1AC1M4PY4X7QZZ", "notaulid", UnknownSubject, ulid.ErrDataSize},
	}

	for i, tc := range testCases {
		actual, err := IdentifySubject(context.Background(), tc.subID, tc.orgID)
		if tc.err != nil {
			require.ErrorIs(err, tc.err, "expected error for test case %d", i)
		} else {
			require.NoError(err, "expected no error for test case %d", i)
			require.Equal(tc.expected, actual, "unexpected subject type returned")
		}
	}
}
