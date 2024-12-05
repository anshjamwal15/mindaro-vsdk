package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aimbot1526/mindaro-vsdk/db"
	"github.com/aimbot1526/mindaro-vsdk/handlers"
	"github.com/aimbot1526/mindaro-vsdk/repositories"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Server struct to hold server dependencies and configuration
type Server struct {
	Router  *mux.Router
	DB      *gorm.DB
	HTTPSrv *http.Server
}

// StartServer initializes and starts the server
func StartServer() *Server {
	log.Println("Initializing server...")

	// Initialize database
	db := db.InitializeDB()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	groupHandler := handlers.NewGroupHandler(groupRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	websocketHandler := handlers.NewWebSocketHandler(groupRepo, messageRepo)

	// Initialize router and setup routes
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/user/create", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/user/get", userHandler.GetUserByUsername).Methods("GET")

	// Group routes
	router.HandleFunc("/group/create", groupHandler.CreateGroup).Methods("POST")
	router.HandleFunc("/group/{group_id}/join", groupHandler.JoinGroup).Methods("POST")

	// Message routes
	router.HandleFunc("/group/{group_id}/message/send", messageHandler.SendMessageToGroup).Methods("POST")

	// WebSocket route for group chat
	router.HandleFunc("/ws/group", websocketHandler.GroupWebSocketHandler)

	// Create HTTP server
	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server := &Server{
		Router:  router,
		DB:      db,
		HTTPSrv: httpSrv,
	}

	// Start server in a separate goroutine
	go func() {
		log.Println("Starting server on http://localhost:8080")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return server
}

// ShutdownServer gracefully shuts down the server
func (s *Server) ShutdownServer() {
	log.Println("Shutting down server gracefully...")

	// Context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := s.HTTPSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}

	log.Println("Closing database connection...")
	sqlDB, err := s.DB.DB() // Get the underlying sql.DB instance
	if err != nil {
		log.Printf("Error accessing database instance: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	log.Println("Server stopped.")
}
