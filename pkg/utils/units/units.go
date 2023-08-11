package units

type ByteUnit uint64

const (
	Bytes ByteUnit = 1 << (10 * iota)
	Kilobytes
	Megabytes
	Gigabytes
	Terabytes
)

// String returns the string representation of the byte unit.
func (b ByteUnit) String() string {
	switch b {
	case Kilobytes:
		return "KiB"
	case Megabytes:
		return "MiB"
	case Gigabytes:
		return "GiB"
	case Terabytes:
		return "TiB"
	default:
		return "B"
	}
}

// FromBytes returns the appropriate byte unit and value for the provided number
// of bytes for human readability (e.g. 1024 bytes = 1 KB).
func FromBytes(bytes uint64) (units string, value float64) {
	switch {
	case bytes >= uint64(Terabytes):
		return Terabytes.String(), float64(bytes) / float64(Terabytes)
	case bytes >= uint64(Gigabytes):
		return Gigabytes.String(), float64(bytes) / float64(Gigabytes)
	case bytes >= uint64(Megabytes):
		return Megabytes.String(), float64(bytes) / float64(Megabytes)
	case bytes >= uint64(Kilobytes):
		return Kilobytes.String(), float64(bytes) / float64(Kilobytes)
	default:
		return Bytes.String(), float64(bytes)
	}
}
