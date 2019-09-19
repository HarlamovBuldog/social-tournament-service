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

type userID struct {
	ID string `json:"userID"`
}

type winnerUserID struct {
	ID string `json:"winnerUserID"`
}
type userName struct {
	Name string `json:"name"`
}

type userPoints struct {
	Points float64 `json:"points"`
}

type tournamentInit struct {
	Name    string  `json:"name"`
	Deposit float64 `json:"deposit"`
}

type tournamentID struct {
	ID string `json:"id"`
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
	router.HandleFunc("/tournament/{id}", s.getTournamentInfo).Methods("GET")
	router.HandleFunc("/tournament/{id}/join", s.joinTournament).Methods("POST")
	router.HandleFunc("/tournament/{id}/finish", s.finishTournament).Methods("POST")
	router.HandleFunc("/tournament/{id}", s.cancelTournament).Methods("DELETE")

	return &s
}

func (s *Server) createNewUser(w http.ResponseWriter, req *http.Request) {
	var user userName
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("createNewUser: can't decode request body: %v", err)
		return
	}

	usrID, err := s.service.AddUser(req.Context(), user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error createNewUser: %v", err)
		return
	}

	err = json.NewEncoder(w).Encode(userID{
		ID: usrID,
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
	var points userPoints
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
	var points userPoints
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
	var tournament tournamentInit
	err := json.NewDecoder(req.Body).Decode(&tournament)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("createNewTournament: can't decode request body: %s", err)
		return
	}

	tourneyID, err := s.service.AddTournament(req.Context(), tournament.Name, tournament.Deposit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("createNewTournament: %s", err)
		return
	}

	err = json.NewEncoder(w).Encode(tournamentID{
		ID: tourneyID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("createNewTournament: error encoding json: %s", err)
		return
	}
}

func (s *Server) getTournamentInfo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tournamentID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("getTournamentInfo: tournament id is not provided")
		return
	}

	tournament, err := s.service.GetTournament(req.Context(), tournamentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getTournamentInfo: %s", err)
		return
	}

	if err = json.NewEncoder(w).Encode(&tournament); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getTournamentInfo: error encoding json: %s", err)
		return
	}
}

func (s *Server) joinTournament(w http.ResponseWriter, req *http.Request) {
	var usrID userID
	err := json.NewDecoder(req.Body).Decode(&usrID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("joinTournament: can't decode request body: %s", err)
		return
	}

	vars := mux.Vars(req)
	tournamentID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("joinTournament: tournament id is not provided")
		return
	}

	err = s.service.JoinTournament(req.Context(), tournamentID, usrID.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("joinTournament: %s", err)
		return
	}
}

func (s *Server) finishTournament(w http.ResponseWriter, req *http.Request) {
	var winnerUsrID winnerUserID
	err := json.NewDecoder(req.Body).Decode(&winnerUsrID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("finishTournament: can't decode request body: %s", err)
		return
	}

	vars := mux.Vars(req)
	tournamentID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("finishTournament: tournament id is not provided")
		return
	}

	err = s.service.FinishTournament(req.Context(), tournamentID, winnerUsrID.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("finishTournament: %s", err)
		return
	}
}

func (s *Server) cancelTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tournamentID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("cancelTournament: tournament id is not provided")
		return
	}

	err := s.service.DeleteTournament(req.Context(), tournamentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cancelTournament: %s", err)
		return
	}
}
