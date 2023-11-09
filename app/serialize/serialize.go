package serial

type Serialize interface {
	GetObjectID() int
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}
