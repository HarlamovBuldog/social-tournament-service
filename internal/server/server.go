package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
	"github.com/gorilla/mux"
)

type Server struct {
	http.Handler
	service storage.Service
}

type userIDJSON struct {
	ID string `json:"id"`
}

type userNameJSON struct {
	Name string `json:"name"`
}

type userPointsJSON struct {
	Points float64 `json:"points"`
}

// NewServer initializes router and entrypoints
func NewServer(db storage.Service) *Server {
	router := mux.NewRouter()

	s := Server{
		service: db,
		Handler: router,
	}
	router.HandleFunc("/user", s.createNewUser).Methods("POST")
	router.HandleFunc("/user/{id}", s.getUserInfo).Methods("GET")
	router.HandleFunc("/user/{id}", s.removeUser).Methods("DELETE")
	router.HandleFunc("/user/{id}/take", s.takeUserBonusPoints).Methods("POST")
	router.HandleFunc("/user/{id}/fund", s.addUserBonusPoints).Methods("POST")

	router.HandleFunc("/tournament", s.createNewTournament).Methods("POST")
	router.HandleFunc("/tournament/{id}", s.getcreateNewTournamentInfo).Methods("GET")
	router.HandleFunc("/tournament/{id}/join", s.joincreateNewTournament).Methods("POST")
	router.HandleFunc("/tournament/{id}/finish", s.finishTournament).Methods("POST")
	router.HandleFunc("/tournament/{id}", s.deleteTournament).Methods("DELETE")

	return &s
}

func (s *Server) createNewUser(w http.ResponseWriter, req *http.Request) {
	var user userNameJSON
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("createNewUser: can't decode request body: %v", err)
		return
	}

	userID, err := s.service.AddUser(req.Context(), user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error createNewUser: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(userIDJSON{
		ID: userID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("createNewUser: error encoding json: %v", err)
		return
	}
}

func (s *Server) getUserInfo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("getUserInfo: user id is not provided")
		return
	}

	userData, err := s.service.GetUser(req.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getUserInfo: %v", err)
		return
	}
	enc := json.NewEncoder(w)
	if err = enc.Encode(&userData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error getting user info: error encoding json: %v\n", err)
		return
	}
}

func (s *Server) removeUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("removeUser: user id is not provided")
		return
	}
	err := s.service.DeleteUser(req.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("removeUser: %v", err)
		return
	}
}

func (s *Server) takeUserBonusPoints(w http.ResponseWriter, req *http.Request) {
	var points userPointsJSON
	err := json.NewDecoder(req.Body).Decode(&points)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("takeUserBonusPoints: can't decode request body: %v", err)
		return
	}

	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("takeUserBonusPoints: user id is not provided")
		return
	}

	err = s.service.TakeUserBalance(req.Context(), userID, points.Points)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("takeUserBonusPoints: %v", err)
		return
	}
}

func (s *Server) addUserBonusPoints(w http.ResponseWriter, req *http.Request) {
	var points userPointsJSON
	err := json.NewDecoder(req.Body).Decode(&points)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("addUserBonusPoints: can't decode request body: %v", err)
		return
	}

	vars := mux.Vars(req)
	userID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("addUserBonusPoints: user id is not provided")
		return
	}

	err = s.service.FundUserBalance(req.Context(), userID, points.Points)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("addUserBonusPoints: %v", err)
		return
	}
}

func (s *Server) createNewTournament(w http.ResponseWriter, req *http.Request) {

}
func (s *Server) getcreateNewTournamentInfo(w http.ResponseWriter, req *http.Request) {
}
func (s *Server) joincreateNewTournament(w http.ResponseWriter, req *http.Request) {
}
func (s *Server) finishTournament(w http.ResponseWriter, req *http.Request) {
}
func (s *Server) deleteTournament(w http.ResponseWriter, req *http.Request) {
}
