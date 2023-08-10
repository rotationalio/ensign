package units_test

import (
	"testing"

	"github.com/rotationalio/ensign/pkg/utils/units"
	"github.com/stretchr/testify/require"
)

func TestByteUnits(t *testing.T) {
	// Test that bytes are correctly converted to the expected units.
	tests := []struct {
		bytes uint64
		units string
		value float64
	}{
		{bytes: 0, units: "B", value: 0},
		{bytes: 1, units: "B", value: 1},
		{bytes: 1023, units: "B", value: 1023},
		{bytes: 1024, units: "KB", value: 1},
		{bytes: (1024 * 1024) + (1024 * 512), units: "MB", value: 1.5},
		{bytes: 1024 * 1024 * 1024, units: "GB", value: 1},
		{bytes: 1024 * 1024 * 1024 * 1024, units: "TB", value: 1},
	}

	for _, test := range tests {
		units, value := units.FromBytes(test.bytes)
		require.Equal(t, test.units, units, "wrong units")
		require.Equal(t, test.value, value, "wrong value")
	}
}
