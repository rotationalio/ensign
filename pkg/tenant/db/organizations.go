package db

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/vmihailenco/msgpack/v5"
)

const OrganizationNamespace = "organizations"

type Organization struct {
	ID       ulid.ULID `msgpack:"id"`
	Name     string    `msgpack:"name"`
	Created  time.Time `msgpack:"created"`
	Modified time.Time `msgpack:"modified"`
}

var _ Model = &Organization{}

func (o *Organization) Key() (key []byte, err error) {
	if ulids.IsZero(o.ID) {
		return nil, ErrMissingOrgID
	}

	return o.ID.MarshalBinary()
}

func (o *Organization) Namespace() string {
	return OrganizationNamespace
}

func (o *Organization) MarshalValue() ([]byte, error) {
	return msgpack.Marshal(o)
}

func (o *Organization) UnmarshalValue(data []byte) error {
	return msgpack.Unmarshal(data, o)
}

func VerifyOrg(orgID ulid.ULID, modelOrgID ulid.ULID) (bool, error) {
	if ulids.IsZero(orgID) {
		return false, ErrMissingOrgID
	}

	if ulids.IsZero(modelOrgID) {
		return false, nil
	}

	if orgID.Compare(modelOrgID) == 0 {
		return true, nil
	} else {
		return false, ErrOrgNotVerified
	}
}
