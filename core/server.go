package core

type Server interface {
	GetIP() (string, error)
	GetPort() (uint16, error)
	GetIdentifier() (string, error)
	GetMetadata() (map[string]string, error)
}
