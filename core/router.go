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

func lobbyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action, ok := vars["action"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Lobby: %s", action)

}

func init() {
	Router = mux.NewRouter()
	Router.HandleFunc("/", exampleHandler)
	Router.HandleFunc("/lobby/{action}", lobbyHandler)
}
