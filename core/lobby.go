package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

func lobbyIsValid(lobby string) bool {
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Print("failed to compile lobby regex")
		return false
	}
	return len(lobby) > 1 && !strings.Contains(lobby, " ") && !re.Match([]byte(lobby))
}

type getLobbyResponse struct {
	Name  string      `json:"name"`
	State ServerState `json:"state"`
}

func GetLobbyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lobby, ok := vars["lobby"]
	if !ok || !lobbyIsValid(lobby) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ERROR OCCURS BETWEEN HERE
	server, err := ActiveModules.server.GetServer(fmt.Sprintf("lobby:%v", lobby))
	if err != nil {
		log.Printf("error occured getting lobby server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	// AND HERE

	if server != nil {
		serverState, err := server.GetState()
		if err != nil {
			log.Printf("error occured getting lobby state: %v", err)
		}
		data, err := json.Marshal(getLobbyResponse{
			Name:  lobby,
			State: serverState,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("error occured writing response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}
func JoinLobbyHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func AddLobbyRoutes(router *mux.Router) {
	r := router.PathPrefix("/lobby").Subrouter()
	r.HandleFunc("/get/{lobby}", GetLobbyHandler).Methods("GET")
	r.HandleFunc("/join/{lobby}", JoinLobbyHandler).Methods("POST")
}
