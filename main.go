package main

import (
	"log"
	"net/http"
	"os"

	"github.com/loudsmilestudios/TetraControl/core"
)

func main() {

	log.Print("Hello")

	if _, isAwsLambda := os.LookupEnv("AWS_EXECUTION_ENV"); isAwsLambda {
		core.AwsLambdaHandler()
		return
	}

	s := &http.Server{
		Addr:    ":8080",
		Handler: core.Router,
	}

	log.Fatal(s.ListenAndServe())
}
