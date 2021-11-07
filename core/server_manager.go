package core

type ServerManager interface {
	CreateServer(string) (Server, error)
	GetServer(string) (Server, error)
	GetServerCount() (uint, error)
	DeleteServerByIdentifier(string) error
	DeleteServer(Server) error
}
