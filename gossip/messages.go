package gossip

type MessageType uint8

const (
	Reply   MessageType = iota
	Request MessageType = iota
)

type PViewMessage struct {
	Type     MessageType
	Capacity int
	View     []Descriptor
	Coords   []RemoteCoord
}
