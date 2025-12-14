package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahadahamedakash/students-api-with-go/internal/config"
	"github.com/ahadahamedakash/students-api-with-go/internal/http/handlers/student"
	"github.com/ahadahamedakash/students-api-with-go/internal/storage/sqlite"
)

func main() {
	fmt.Println("Welcome to students api")
	// load config
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server started %s", slog.String("Address", cfg.Address))

	// fmt.Println("Server is running!")
	fmt.Printf("Server is running at %s", cfg.HTTPServer.Address)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Failed to start server!")
		}
	}()

	<-done

	slog.Info("Shutting down the server!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("Error", err.Error()))
	}

	slog.Info("Server shutdown successfully!")

}
