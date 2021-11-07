package core

type ModuleSet struct {
	server ServerManager
}

var ActiveModules ModuleSet = ModuleSet{}
