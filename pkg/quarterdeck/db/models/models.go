package models

import "time"

// The Base model provides model audit functionality for setting created and modified
// timestamps in the database so we can track how rows are being modified over time.
type Base struct {
	Created  string
	Modified string
}

// Scanner is an interface for *sql.Rows and *sql.Row so that models can implement how
// they scan fields into their struct without having to specify every field every time.
type Scanner interface {
	Scan(dest ...any) error
}

// Return the parsed created timestamp.
func (b *Base) GetCreated() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, b.Created)
}

// Return the parsed modified timestamp.
func (b *Base) GetModified() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, b.Modified)
}

// Sets the created timestamp as the string formatted representation of the ts.
func (b *Base) SetCreated(ts time.Time) {
	b.Created = ts.Format(time.RFC3339Nano)
}

// Sets the modified timestamp as the string formatted representation of the ts.
func (b *Base) SetModified(ts time.Time) {
	b.Modified = ts.Format(time.RFC3339Nano)
}
