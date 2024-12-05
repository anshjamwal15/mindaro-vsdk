package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aimbot1526/mindaro-vsdk/models"
	"github.com/aimbot1526/mindaro-vsdk/repositories"
	"github.com/gorilla/mux"
)

type GroupHandler struct {
	GroupRepo *repositories.GroupRepository
}

func NewGroupHandler(groupRepo *repositories.GroupRepository) *GroupHandler {
	return &GroupHandler{GroupRepo: groupRepo}
}

// CreateGroup allows users to create a group.
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var groupRequest struct {
		Name      string `json:"name"`
		IsPrivate bool   `json:"is_private"`
		CreatorID uint   `json:"creator_id"`
		Members   []uint `json:"members"`
	}

	err := json.NewDecoder(r.Body).Decode(&groupRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group := &models.Group{
		Name:      groupRequest.Name,
		IsPrivate: groupRequest.IsPrivate,
		CreatorID: groupRequest.CreatorID,
	}
	if err := h.GroupRepo.CreateGroup(group); err != nil {
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	for _, memberID := range groupRequest.Members {
		if err := h.GroupRepo.AddMember(group.ID, memberID); err != nil {
			http.Error(w, "Failed to add members to group", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Group created successfully",
		"group_id":   group.ID,
		"group_name": group.Name,
	})
}

// JoinGroup allows a user to join an existing group.
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	groupIDStr := mux.Vars(r)["group_id"]
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	var joinRequest struct {
		UserID uint `json:"user_id"`
	}

	err = json.NewDecoder(r.Body).Decode(&joinRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.GroupRepo.AddMember(uint(groupID), joinRequest.UserID); err != nil {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User added to group successfully",
	})
}
