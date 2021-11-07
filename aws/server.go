package aws

// Server represents a game server running on AWS
type Server struct {
	Identifier string            `json:"identifier"`
	TaskArn    string            `json:"task"`
	Metadata   map[string]string `json:"metadata"`
}

// GetIP returns the IP address of a server
func (server *Server) GetIP() (string, error) {
	return "127.0.0.1", nil
}

// GetPort returns the port of a server
func (server *Server) GetPort() (uint16, error) {
	return 7777, nil
}

// GetIdentifier returns the identifier of a server
func (server *Server) GetIdentifier() (string, error) {
	return "null", nil
}

// GetMetadata returns the metadata of a server
func (server *Server) GetMetadata() (map[string]string, error) {
	return server.Metadata, nil
}
