package health

type Versioned interface {
	Version() string
}

// Checks if the service statuses are versioned, and if so, if the version has changed.
func VersionChanged(first, second ServiceStatus) bool {
	var (
		ok      bool
		firstv  Versioned
		secondv Versioned
	)

	if firstv, ok = first.(Versioned); !ok {
		return false
	}

	if secondv, ok = second.(Versioned); !ok {
		return false
	}

	return firstv.Version() != secondv.Version()
}
