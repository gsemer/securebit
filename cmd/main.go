package main

import (
	"log"
	"securebit/persistence"
	"securebit/presentation"

	"github.com/gorilla/mux"
)

func main() {
	config := persistence.LoadPostgresConfig()
	db, err := config.Init()
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	userRepo := persistence.NewUserRepository(db)
	handler := presentation.NewAuthHandler(userRepo)

	r := mux.NewRouter()
	api := r.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/register", handler.Register).Methods("POST")
	api.HandleFunc("/login", nil).Methods("POST")
	api.HandleFunc("/logout", nil).Methods("POST")
	api.HandleFunc("/refresh", nil).Methods("POST")
	api.HandleFunc("/validate", nil).Methods("POST")

	server := presentation.NewCustomServer(r, "localhost:8080")
	server.Start()
}
