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

type joinLobbyResponse struct {
	Name  string      `json:"name"`
	State ServerState `json:"state"`
	Host  string      `json:"host"`
	Port  uint16      `json:"port"`
}

// DeleteLobbyHandler get server data and sends a delete request
func DeleteLobbyHandler(w http.ResponseWriter, r *http.Request) {
	// Eventually add auth system here
	if true {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	lobby, ok := vars["lobby"]
	if !ok || !lobbyIsValid(lobby) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := ActiveModules.Server.GetServer(fmt.Sprintf("lobby:%v", lobby))
	if err != nil {
		log.Printf("error occured when getting server for delection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = ActiveModules.Server.DeleteServer(server)
	if err != nil {
		log.Printf("error occured when deleting server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// GetLobbyHandler grabs the lobby status and returns it to the client
func GetLobbyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lobby, ok := vars["lobby"]
	if !ok || !lobbyIsValid(lobby) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ActiveModules.Server == nil {
		log.Printf("server manager is not initialized!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server, err := ActiveModules.Server.GetServer(fmt.Sprintf("lobby:%v", lobby))
	if err != nil {
		log.Printf("error occured getting lobby server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

// JoinLobbyHandler returns lobby connection information
func JoinLobbyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lobby, ok := vars["lobby"]
	if !ok || !lobbyIsValid(lobby) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := ActiveModules.Server.GetServer(fmt.Sprintf("lobby:%v", lobby))
	if err != nil {
		log.Printf("error occured getting lobby server: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if server == nil {
		server, err = ActiveModules.Server.CreateServer(fmt.Sprintf("lobby:%v", lobby))
		if err != nil {
			log.Printf("error occured creating lobby server: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if server != nil {
		serverState, err := server.GetState()
		if err != nil {
			log.Printf("error occured getting lobby state: %v", err)
		}

		IP, err := server.GetIP()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		Port, err := server.GetPort()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(joinLobbyResponse{
			Name:  lobby,
			State: serverState,
			Host:  IP,
			Port:  Port,
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

// AddLobbyRoutes updaates the mux router will all lobby routes
func AddLobbyRoutes(router *mux.Router) {
	r := router.PathPrefix("/lobby").Subrouter()
	r.HandleFunc("/get/{lobby}", GetLobbyHandler).Methods("GET")
	r.HandleFunc("/join/{lobby}", JoinLobbyHandler).Methods("POST")
	r.HandleFunc("/delete/{lobby}", DeleteLobbyHandler).Methods("DELETE")
}
