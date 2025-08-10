package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"securebit/config"
	"securebit/persistence"
	"securebit/presentation"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	postgresUsername := config.GetEnv("POSTGRES_USERNAME", "")
	postgresPassword := config.GetEnv("POSTGRES_PASSWORD", "")
	postgresDB := config.GetEnv("POSTGRES_DB", "")
	postgresIP := config.GetEnv("POSTGRES_IP", "")
	postgresPort := config.GetEnv("POSTGRES_PORT", "")

	db, err := persistence.InitDB(postgresUsername, postgresPassword, postgresIP, postgresPort, postgresDB)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	userRepo := persistence.NewUserRepository(db)
	handler := presentation.NewAuthHandler(userRepo)

	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()

	s.HandleFunc("/register", handler.Register).Methods("POST")
	s.HandleFunc("/login", nil).Methods("POST")
	s.HandleFunc("/logout", nil).Methods("POST")
	s.HandleFunc("/refresh", nil).Methods("POST")
	s.HandleFunc("/validate", nil).Methods("POST")

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down")
	server.Shutdown(ctx)
	os.Exit(0)
}
