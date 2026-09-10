package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"example.com/german/backend/internal/database"
	"example.com/german/backend/internal/handlers"
	"example.com/german/backend/internal/repos"
	"example.com/german/backend/internal/service"
	"example.com/german/backend/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env", "../.env")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := database.New()
	if err != nil {
		log.Fatal("❌ ", err)
	}
	defer db.Close()

	healthService := service.NewHealthService()
	healthHandler := handlers.NewHealthHandler(healthService)
	readinessHandler := handlers.NewReadinessHandler(db)

	attachmentStorage, err := storage.NewLocal(getEnv("ATTACHMENT_STORAGE_PATH", "./data/attachments"))
	if err != nil {
		log.Fatal("❌ initialize attachment storage: ", err)
	}
	imageProcessor, err := service.NewImageProcessor(service.DefaultImageProcessorConfig())
	if err != nil {
		log.Fatal("❌ initialize image processor: ", err)
	}
	attachmentPreparer, err := service.NewAttachmentPreparer(imageProcessor, attachmentStorage)
	if err != nil {
		log.Fatal("❌ initialize attachment preparer: ", err)
	}

	ticketRepository := repos.NewTicketRepository(db.DB)
	messageRepository := repos.NewMessageRepository(db.DB)
	attachmentRepository := repos.NewAttachmentRepository(db.DB)
	userRepository := repos.NewUserRepository(db.DB)
	operatorRepository := repos.NewOperatorRepository(db.DB)
	expertRepository := repos.NewExpertRepository(db.DB)
	ticketService, err := service.NewTicketService(ticketRepository)
	if err != nil {
		log.Fatal("❌ initialize ticket service: ", err)
	}
	messageService, err := service.NewTicketMessageService(messageRepository, attachmentPreparer, attachmentStorage)
	if err != nil {
		log.Fatal("❌ initialize ticket message service: ", err)
	}
	attachmentService, err := service.NewTicketAttachmentService(attachmentRepository, attachmentStorage)
	if err != nil {
		log.Fatal("❌ initialize ticket attachment service: ", err)
	}
	ticketHandler, err := handlers.NewTicketHandler(ticketService, messageService, attachmentService)
	if err != nil {
		log.Fatal("❌ initialize ticket handler: ", err)
	}
	authService, err := service.NewAuthService(userRepository, os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatal("❌ initialize auth service: ", err)
	}
	operatorService, err := service.NewOperatorService(operatorRepository)
	if err != nil {
		log.Fatal("❌ initialize operator service: ", err)
	}
	expertService, err := service.NewExpertService(expertRepository)
	if err != nil {
		log.Fatal("❌ initialize expert service: ", err)
	}
	authHandler, err := handlers.NewAuthHandler(authService)
	if err != nil {
		log.Fatal("❌ initialize auth handler: ", err)
	}
	operatorHandler, err := handlers.NewOperatorHandler(operatorService, authService)
	if err != nil {
		log.Fatal("❌ initialize operator handler: ", err)
	}
	expertHandler, err := handlers.NewExpertHandler(expertService, authService, attachmentService)
	if err != nil {
		log.Fatal("❌ initialize expert handler: ", err)
	}

	mux := http.NewServeMux()
	healthHandler.RegisterRoutes(mux)
	readinessHandler.RegisterRoutes(mux)
	ticketHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)
	operatorHandler.RegisterRoutes(mux)
	expertHandler.RegisterRoutes(mux)

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

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
