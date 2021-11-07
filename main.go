package main

import (
	"log"
	"net/http"
	"os"

	"github.com/loudsmilestudios/TetraControl/aws"
	"github.com/loudsmilestudios/TetraControl/core"
)

func main() {

	log.Print("Starting TetraControl...")

	if _, isAwsLambda := os.LookupEnv("AWS_EXECUTION_ENV"); isAwsLambda {
		adapter := aws.ServerlessAdapter{}
		adapter.Start()
		return
	}

	s := &http.Server{
		Addr:    ":8080",
		Handler: core.Router,
	}

	log.Fatal(s.ListenAndServe())
}
