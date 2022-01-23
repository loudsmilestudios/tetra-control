package core

// ModuleSet is used to hold a module of each type
type ModuleSet struct {
	Server ServerManager
}

// ActiveModules is a ModuleSet holding all currently active modules
var ActiveModules ModuleSet = ModuleSet{}
