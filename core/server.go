package core

// ServerState is an int that represents the servers state
type ServerState uint8

// ServerStates increment from 0, where 0 is Unknown
const (
	Unknown ServerState = iota
	Starting
	Active
	Error
	Exiting
)

// Server is an interface used to represent a running server
type Server interface {
	GetIP() (string, error)
	GetPort() (uint16, error)
	GetIdentifier() (string, error)
	GetMetadata() (map[string]string, error)
	GetState() (ServerState, error)
}
