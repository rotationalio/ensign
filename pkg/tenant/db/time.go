package db

import "time"

// Helper to convert a time.Time to an RFC3339Nano string for JSON serialization.
func TimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339Nano)
}
