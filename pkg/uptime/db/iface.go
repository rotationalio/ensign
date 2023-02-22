package db

type Model interface {
	Key() ([]byte, error)
	Unmarshal([]byte) error
	Marshal() ([]byte, error)
}
