package core

// ServerManager is an interface that is used by TetraControl to orchestrate servers
type ServerManager interface {
	CreateServer(string) (Server, error)
	GetServer(string) (Server, error)
	GetServerCount() (uint, error)
	DeleteServerByIdentifier(string) error
	DeleteServer(Server) error
}
