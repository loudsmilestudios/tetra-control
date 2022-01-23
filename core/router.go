package core

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Router is the primary router utilized by TetraControl
var Router *mux.Router

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(Response{
		Success: true,
		Message: "OK",
	})
	w.Write(data)
}

func init() {
	Router = mux.NewRouter()
	Router.HandleFunc("/status", statusHandler).Methods("GET")
	AddLobbyRoutes(Router)
}
