package core

import (
	"errors"
)

// ModuleSet is used to hold a module of each type
type ModuleSet struct {
	Server     ServerManager
	Serverless ServerlessAdapter
}

// ActiveModules is a ModuleSet holding all currently active modules
var ActiveModules ModuleSet = ModuleSet{}
var serverManagerModules map[string]ServerManager = map[string]ServerManager{}
var severlessModules map[string]ServerlessAdapter = map[string]ServerlessAdapter{}

// RegisterServerModule
func RegisterServerModule(id string, module ServerManager) {
	serverManagerModules[id] = module
}

// SetServerModule
func SetServerModule(id string) error {
	_, ok := serverManagerModules[id]
	if !ok {
		return errors.New("ID is not a valid server module")
	}

	ActiveModules.Server = serverManagerModules[id]
	return nil
}

// RegisterServerlessModule
func RegisterServerlessModule(id string, module ServerlessAdapter) {
	severlessModules[id] = module
}

// SetServerlessModule
func SetServerlessModule(id string) error {
	_, ok := severlessModules[id]
	if !ok {
		return errors.New("ID is not a valid serverless module")
	}

	ActiveModules.Serverless = severlessModules[id]
	return nil
}
