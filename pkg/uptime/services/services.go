/*
Package services provides an Info and Service data structure that is used for two
purposes: loading configuration information from disk and storing status information in
the database.
*/
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/ensign/pkg"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/vmihailenco/msgpack/v5"
)

type ServiceType string

const (
	UnknownServiceType ServiceType = "unknown"
	HTTPServiceType    ServiceType = "http"
	APIServiceType     ServiceType = "api"
)

type Info struct {
	Version  string     `json:"version" msgpack:"version"`
	Services []*Service `json:"services" msgpack:"services"`
	Updated  time.Time  `json:"-" msgpack:"updated"`
}

type Service struct {
	ID          uuid.UUID     `json:"id" msgpack:"id"`
	Title       string        `json:"title" msgpack:"title"`
	Description string        `json:"description,omitempty" msgpack:"description"`
	Type        ServiceType   `json:"type" msgpack:"-"`
	Endpoint    string        `json:"endpoint" msgpack:"-"`
	Status      health.Status `json:"-" msgpack:"status"`
	LastUpdate  time.Time     `json:"-" msgpack:"last_update"`
}

var versre *regexp.Regexp = regexp.MustCompile(`^v?\d+$`)

// Load services info from a path on disk.
func Load(path string) (info *Info, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&info); err != nil {
		return nil, err
	}

	if err = info.Validate(); err != nil {
		return nil, err
	}

	return info, nil
}

func (i *Info) Dump(path string) (err error) {
	// Set the version info
	i.Version = fmt.Sprintf("v%d", pkg.VersionMajor)

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(i); err != nil {
		return err
	}
	return nil
}

func (i *Info) Validate() (err error) {
	if !versre.MatchString(i.Version) {
		return fmt.Errorf("could not parse version info %q", i.Version)
	}

	vers, _ := strconv.Atoi(i.Version[1:])
	if vers != pkg.VersionMajor {
		return fmt.Errorf("could not load status version %d in package version %s", vers, pkg.Version())
	}

	for _, service := range i.Services {
		if err = service.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Validate() (err error) {
	if s.ID == uuid.Nil {
		return errors.New("invalid service configuration: missing service id")
	}

	if s.Title == "" {
		return errors.New("invalid service configuration: missing service title")
	}

	if s.Type == "" {
		return errors.New("invalid service configuration: missing service type")
	}

	if s.Endpoint == "" {
		return errors.New("invalid service configuration: missing service endpoint")
	}

	return nil
}

func (i *Info) Key() ([]byte, error) {
	return db.KeyCurrentStatus, nil
}

func (i *Info) Marshal() ([]byte, error) {
	return msgpack.Marshal(i)
}

func (i *Info) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, i)
}
