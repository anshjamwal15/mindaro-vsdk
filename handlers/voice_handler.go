package handlers

import (
	"net/http"

	"github.com/aimbot1526/mindaro-vsdk/utils"
	"github.com/aimbot1526/mindaro-vsdk/webrtc"
	"gorm.io/gorm"
)

type VoiceHandler struct {
	SFU *webrtc.SFU
	db  *gorm.DB
}

// NewVoiceHandler creates a new instance of VoiceHandler.
func NewVoiceHandler(sfu *webrtc.SFU, db *gorm.DB) *VoiceHandler {
	return &VoiceHandler{
		SFU: sfu,
		db:  db,
	}
}

// HandleVoiceCall initializes WebSocket signaling for voice calls.
func (vh *VoiceHandler) HandleVoiceCall(w http.ResponseWriter, r *http.Request) {
	conn, err := utils.AcceptWebSocket(w, r)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		http.Error(w, "Client ID is required", http.StatusBadRequest)
		return
	}

	// Start the signaling process:
	err = webrtc.HandleSignaling(conn, vh.db, vh.SFU, clientID)
	if err != nil {
		http.Error(w, "Signaling error: "+err.Error(), http.StatusInternalServerError)
	}
}
