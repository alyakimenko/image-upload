package main

import (
	"context"
	"github.com/alyakimenko/image-upload/internal/server"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, loading from default.")
	}
}

func main() {
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	config := server.NewConfig()
	if err := env.Parse(config); err != nil {
		log.Fatalf("%+v\n", err)
	}
	if err := os.MkdirAll(config.DownloadedPath, 0755); err != nil {
		log.Fatalf("Cannot create %s directory.", config.DownloadedPath)
	}

	con := server.NewController(config)

	go gracefullyShutdown(con.Server, quit, done)

	log.Println("Server is started on ", config.BindAddr)
	if err := con.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", config.BindAddr, err)
	}

	<-done
}

func gracefullyShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}
