package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/loudsmilestudios/TetraControl/aws"
	"github.com/loudsmilestudios/TetraControl/core"
)

func main() {

	// Register built in modules
	// future plugins should have similar
	// login in their init function
	aws.RegisterModules()

	if _, err := os.Stat("config.yaml"); err == nil {
		// Load config from file
		if err := cleanenv.ReadConfig("config.yaml", &core.Config); err != nil {
			log.Printf("could not load config from file: %v", err)
		}

	} else if os.IsNotExist(err) {
		// Load config from environment
		if err := cleanenv.ReadEnv(&core.Config); err != nil {
			log.Printf("could not load config from environment: %v", err)
		}
	}
	if core.Config.ServerModuleID == "" {
		log.Fatal("A server module is not set!")
	}

	log.Print("Loading modules")
	if err := core.SetServerModule(core.Config.ServerModuleID); err != nil {
		log.Fatalf("could not load module '%v': %v", core.Config.ServerModuleID, err)
	}
	if err := core.SetServerlessModule(core.Config.ServerlessModuleID); err != nil {
		log.Printf("could not load serverless module '%v': %v", core.Config.ServerlessModuleID, err)
	}

	log.Print("Starting TetraControl")
	if err := core.ActiveModules.Server.Initialize(); err != nil {
		log.Fatalf("failed to initialize server manager: %v", err)
	}

	if core.ActiveModules.Serverless != nil {
		if core.ActiveModules.Serverless.IsServerless() {
			core.ActiveModules.Serverless.Start()
			return
		}
	}

	s := &http.Server{
		Addr:    ":8080",
		Handler: core.Router,
	}

	log.Fatal(s.ListenAndServe())
}
