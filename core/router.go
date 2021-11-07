package core

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Router is the primary router utilized by TetraControl
var Router *mux.Router

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello world")
}

func init() {
	Router = mux.NewRouter()
	Router.HandleFunc("/", exampleHandler)
	AddLobbyRoutes(Router)
}
