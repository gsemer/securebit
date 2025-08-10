package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"securebit/persistence"
	"securebit/presentation"
	"time"

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
