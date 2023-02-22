package incident

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/vmihailenco/msgpack/v5"
)

var Prefix = []byte{0, 105, 110, 99, 105, 100, 101, 110, 116, 32, 104, 105, 115, 116, 111, 114, 121}

// Returns the last week's worth of incidents or an empty incident group if none exists.
func LastWeek() (_ []*Group, err error) {
	history := make([]*Group, 0, 7)

	date := Today()
	for i := 0; i < 7; i++ {
		group := &Group{
			Date: date.AddDate(0, 0, -1*i),
		}

		var key []byte
		if key, err = group.Key(); err != nil {
			return nil, err
		}

		if err = db.Get(key, group); err != nil {
			if !errors.Is(err, db.ErrNotFound) {
				return nil, err
			}
		}

		history = append(history, group)
	}

	return history, nil
}

// Create the incident in the correct incident day based on the start time.
func Create(incident *Incident) (err error) {
	if incident.ServiceID == uuid.Nil {
		return errors.New("cannot create an incident without a service ID")
	}

	if incident.StartTime.IsZero() {
		return errors.New("cannot create an incident without a start time")
	}

	// Attempt to fetch the group if it already exists
	group := &Group{
		Date: DateFromTime(incident.StartTime),
	}

	var key []byte
	if key, err = group.Key(); err != nil {
		return err
	}

	// Ensure that creating the incident occurs in a transaction
	var tx *db.Transaction
	if tx, err = db.BeginTx(); err != nil {
		return err
	}
	defer tx.Discard()

	if err = tx.Get(key, group); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			group.Incidents = make([]*Incident, 0)
		} else {
			return err
		}
	}

	// Add the incident to the group
	// TODO: should we insort this to make sure they're in time order?
	group.Incidents = append(group.Incidents, incident)

	// Save the group back to disk
	if err = tx.Put(group); err != nil {
		return err
	}

	return tx.Commit()
}

// Find the incident with the specified serviceID and startTime and update its endTime.
func Conclude(serviceID uuid.UUID, startTime, endTime time.Time) (err error) {
	return Update(serviceID, startTime, func(i *Incident) error {
		i.EndTime = endTime
		return nil
	})
}

func Update(serviceID uuid.UUID, startTime time.Time, update func(*Incident) error) (err error) {
	// Attempt to fetch the group if it already exists
	group := &Group{
		Date: DateFromTime(startTime),
	}

	var key []byte
	if key, err = group.Key(); err != nil {
		return err
	}

	// Ensure that creating the incident occurs in a transaction
	var tx *db.Transaction
	if tx, err = db.BeginTx(); err != nil {
		return err
	}
	defer tx.Discard()

	// Get the group
	if err = tx.Get(key, group); err != nil {
		return err
	}

	// Locate the incident
	for _, incident := range group.Incidents {
		if incident.ServiceID == serviceID && incident.StartTime.Equal(startTime) {
			// Perform the update
			if err = update(incident); err != nil {
				return err
			}
		}
	}

	// Save the group back to disk
	if err = tx.Put(group); err != nil {
		return err
	}

	return tx.Commit()
}

// An incident group key is by date, ensuring the date is in the UTC timezone and that
// the timestamp is truncated to midnight (e.g. contains only date components).
func (g *Group) Key() (key []byte, err error) {
	if g.Date.IsZero() || g.Date.Location() != time.UTC || NotADate(g.Date) {
		return nil, errors.New("invalid incident date")
	}

	key = make([]byte, 32)
	copy(key[0:17], Prefix)

	var date []byte
	if date, err = g.Date.MarshalBinary(); err != nil {
		return nil, err
	}

	copy(key[17:], date)
	return key, nil
}

func (g *Group) Marshal() ([]byte, error) {
	return msgpack.Marshal(g)
}

func (g *Group) Unmarshal(data []byte) error {
	if err := msgpack.Unmarshal(data, g); err != nil {
		return err
	}

	// This msgpack library does not have timezone support and returns local times
	g.Date = g.Date.In(time.UTC)
	return nil
}

func Today() time.Time {
	now := time.Now().In(time.UTC)
	return DateFromTime(now)
}

func DateFromTime(ts time.Time) time.Time {
	y, m, d := ts.In(time.UTC).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func NotADate(t time.Time) bool {
	return t.Hour() > 0 || t.Minute() > 0 || t.Second() > 0 || t.Nanosecond() > 0
}
