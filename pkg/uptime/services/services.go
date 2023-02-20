package services

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/rotationalio/ensign/pkg"
)

type ServiceType string

const (
	UnknownServiceType ServiceType = "unknown"
	HTTPServiceType    ServiceType = "http"
	APIServiceType     ServiceType = "api"
)

type Info struct {
	Version  string     `json:"version"`
	Services []*Service `json:"services"`
}

type Service struct {
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Type        ServiceType `json:"type"`
	Endpoint    string      `json:"endpoint"`
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

	if !versre.MatchString(info.Version) {
		return nil, fmt.Errorf("could not parse version info %q", info.Version)
	}

	vers, _ := strconv.Atoi(info.Version[1:])
	if vers != pkg.VersionMajor {
		return nil, fmt.Errorf("could not load status version %d in package version %s", vers, pkg.Version())
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
