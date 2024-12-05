package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aimbot1526/mindaro-vsdk/models"
	"github.com/aimbot1526/mindaro-vsdk/repositories"
)

type UserHandler struct {
	UserRepo *repositories.UserRepository
}

func NewUserHandler(userRepo *repositories.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// GetUserByUsername retrieves a user by their username.
func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing username query parameter", http.StatusBadRequest)
		return
	}

	user, err := h.UserRepo.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the input (for simplicity, let's assume just username and password are required)
	if newUser.Username == "" || newUser.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	// Save the user to the database using the repository method
	createdUser, err := h.UserRepo.CreateUser(&newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	// Return a successful response with the created user details
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
