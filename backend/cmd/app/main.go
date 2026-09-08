package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"example.com/german/backend/internal/handlers"
	"example.com/german/backend/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	ciTestFailure()

	_ = godotenv.Load(".env", "../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	healthService := service.NewHealthService()
	healthHandler := handlers.NewHealthHandler(healthService)

	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("backend started", "address", "http://localhost:"+port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("backend stopped", "error", err)
		os.Exit(1)
	}
}
