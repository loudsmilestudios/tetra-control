package core

type ServerManager interface {
	CreateServer(string) (Server, error)
	GetServer(string) (Server, error)
	DeleteServerByIdentifier(string) error
	DeleteServer(Server) error
}
