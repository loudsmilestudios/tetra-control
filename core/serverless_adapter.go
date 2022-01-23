package core

// ServerManager is an interface that is used to run serverlessly
type ServerlessAdapter interface {
	IsServerless() bool
	Start()
}
