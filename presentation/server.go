package presentation

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type CustomServer struct {
	Addr    string
	Handler http.Handler
}

func NewCustomServer(handler http.Handler, addr string) *CustomServer {
	return &CustomServer{
		Addr:    addr,
		Handler: handler,
	}
}

func (srv *CustomServer) Start() {
	server := &http.Server{Addr: srv.Addr, Handler: srv.Handler}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("Shutting down")
	server.Shutdown(ctx)
}
