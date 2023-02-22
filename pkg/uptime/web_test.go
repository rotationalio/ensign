package uptime_test

import (
	"testing"
	"time"

	"github.com/rotationalio/ensign/pkg/uptime"
	"github.com/stretchr/testify/require"
)

func TestIncidentContext(t *testing.T) {
	testCases := []struct {
		ctx      uptime.IncidentContext
		expected string
	}{
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC),
			},
			expected: "Feb 14, 2023 at 12:32 UTC",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 12, 32, 12, 0, time.FixedZone("BST", -2*3600)),
			},
			expected: "Feb 14, 2023 at 12:32 BST",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC),
			},
			expected: "Feb 14, 2023 at 12:32 UTC",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 14, 6, 32, 12, 0, time.FixedZone("UTC-6", -6*3600)),
			},
			expected: "Feb 14, 2023 at 12:32 UTC",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 14, 14, 49, 1, 201221, time.UTC),
			},
			expected: "Feb 14, 2023 from 12:32 - 14:49 UTC",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 6, 32, 12, 0, time.FixedZone("UTC-6", -6*3600)),
				EndTime:   time.Date(2023, 2, 14, 12, 41, 12, 0, time.UTC),
			},
			expected: "Feb 14, 2023 from 06:32 - 06:41 UTC-6",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2023, 2, 14, 22, 58, 12, 0, time.UTC),
				EndTime:   time.Date(2023, 2, 15, 02, 12, 1, 201221, time.UTC),
			},
			expected: "from Feb 14, 22:58 - Feb 15, 02:12 2023 UTC",
		},
		{
			ctx: uptime.IncidentContext{
				StartTime: time.Date(2022, 12, 25, 22, 58, 12, 0, time.UTC),
				EndTime:   time.Date(2023, 1, 5, 02, 12, 1, 201221, time.UTC),
			},
			expected: "Dec 25 2022, 22:58 UTC - Jan 05 2023, 02:12 UTC",
		},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.ctx.TimeFormat(), "test case %d failed", i)
	}
}

func TestDateEqual(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		date1  time.Time
		date2  time.Time
		assert require.BoolAssertionFunc
	}{
		{now, now, require.True},
		{time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC), time.Date(2023, 2, 14, 9, 14, 58, 0, time.UTC), require.True},
		{time.Date(2021, 11, 5, 8, 21, 46, 0, time.UTC), time.Date(2021, 11, 5, 23, 14, 8, 0, time.UTC), require.True},
		{time.Date(2023, 2, 14, 12, 32, 12, 0, time.UTC), time.Date(2023, 2, 13, 20, 48, 12, 0, time.FixedZone("UTC-6", -6*3600)), require.True},
		{time.Date(2022, 11, 5, 8, 21, 46, 0, time.UTC), time.Date(2021, 11, 5, 23, 14, 8, 0, time.UTC), require.False},
		{time.Date(2021, 11, 5, 8, 21, 46, 0, time.UTC), time.Date(2021, 3, 5, 23, 14, 8, 0, time.UTC), require.False},
		{time.Date(2021, 11, 15, 8, 21, 46, 0, time.UTC), time.Date(2021, 11, 5, 23, 14, 8, 0, time.UTC), require.False},
		{time.Date(2023, 2, 14, 22, 32, 12, 0, time.UTC), time.Date(2023, 2, 14, 22, 32, 12, 0, time.FixedZone("UTC-12", -12*3600)), require.False},
		{now.AddDate(1, 0, 0), now, require.False},
		{now.AddDate(0, 1, 0), now, require.False},
		{now.AddDate(0, 0, 1), now, require.False},
		{now.AddDate(-1, 0, 0), now, require.False},
		{now.AddDate(0, -1, 0), now, require.False},
		{now.AddDate(0, 0, -1), now, require.False},
	}

	for i, tc := range testCases {
		tc.assert(t, uptime.DateEqual(tc.date1, tc.date2), "test case %d failed comparing %s and %s", i, tc.date1, tc.date2)
	}
}
