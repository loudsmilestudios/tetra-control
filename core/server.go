package core

type ServerState uint8

const (
	Unknown ServerState = iota
	Starting
	Active
	Error
	Exiting
)

type Server interface {
	GetIP() (string, error)
	GetPort() (uint16, error)
	GetIdentifier() (string, error)
	GetMetadata() (map[string]string, error)
	GetState() (ServerState, error)
}
