package webrtc

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aimbot1526/mindaro-vsdk/models"
	"github.com/aimbot1526/mindaro-vsdk/utils"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type OfferData struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

type Offer struct {
	Type string    `json:"type"`
	Data OfferData `json:"data"`
}

// CreatePeerConnection initializes and returns a new WebRTC peer connection.
func CreatePeerConnection() (*webrtc.PeerConnection, error) {
	// WebRTC configuration and peer connection initialization
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun.l.google.com:19302",
					"stun:stun1.l.google.com:19302",
					"stun:stun2.l.google.com:19302",
					"stun:stun.1und1.de:3478",
					"stun:stun.gmx.net:3478",
					"stun:stun3.l.google.com:19302",
					"stun:stun4.l.google.com:19302",
					"stun:23.21.150.121:3478",
					"stun:stun.12connect.com:3478",
					"stun:stun.12voip.com:3478",
				},
			},
		},
	}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}
	return peerConnection, nil
}

// AddMediaTracks adds media tracks (audio/video) to the peer connection.
func AddMediaTracks(sfu *SFU, pc *webrtc.PeerConnection, clientID string) error {
	// Add Video Track (replace with actual media source)
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/VP8"}, "video", "video")
	if err != nil {
		log.Println("Error creating video track:", err)
		return err
	}
	err = sfu.AddTrack(clientID, pc, videoTrack)
	if err != nil {
		log.Println("Error adding video track:", err)
		return err
	}

	// Add Audio Track (replace with actual media source)
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "audio")
	if err != nil {
		log.Println("Error creating audio track:", err)
		return err
	}
	err = sfu.AddTrack(clientID, pc, audioTrack)
	if err != nil {
		log.Println("Error adding audio track:", err)
		return err
	}

	return nil
}

// HandleSignaling manages WebSocket signaling for WebRTC with SFU integration.
func HandleSignaling(conn *websocket.Conn, db *gorm.DB, sfu *SFU, clientID string) error {
	peerConnection, err := CreatePeerConnection()
	if err != nil {
		return err
	}
	defer peerConnection.Close()

	// Add media tracks to the peer connection
	if err := AddMediaTracks(sfu, peerConnection, clientID); err != nil {
		return err
	}

	// Add the peer connection to the SFU
	sfu.AddPeerConnection(clientID, peerConnection)

	// Set up ICE candidate handling
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {

			offerData := OfferData{
				Type: "offer",
				Sdp:  candidate.ToJSON().Candidate,
			}
			// Create the signal message
			message := Offer{
				Type: "candidate",
				Data: offerData,
			}

			// Send the candidate to the connection
			if err := conn.WriteJSON(message); err != nil {
				log.Println("Failed to send ICE candidate:", err)
			}

			// Broadcast ICE candidate to all peers via SFU
			sfu.BroadcastICECandidates(candidate.ToJSON())
		}
	})

	// WebSocket signaling message loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket error:", err)
			break
		}

		var signal Offer
		if err := json.Unmarshal(message, &signal); err != nil {
			return errors.New("invalid signaling message")
		}

		switch signal.Type {
		case "offer":
			err := HandleOffer(peerConnection, signal.Data, conn, db, clientID)
			if err != nil {
				return err
			}
		case "answer":
			err := HandleAnswer(peerConnection, signal.Data.Sdp)
			if err != nil {
				return err
			}
		case "candidate":
			err := HandleCandidate(peerConnection, signal.Data.Sdp)
			if err != nil {
				return err
			}
		default:
			log.Println("Unknown signaling type:", signal.Type)
		}
	}

	// Clean up and remove peer connection when done
	sfu.RemovePeerConnection(clientID)

	return nil
}

// HandleOffer processes an offer from a client and sends an answer back.
func HandleOffer(pc *webrtc.PeerConnection, offerData OfferData, conn *websocket.Conn, db *gorm.DB, clientID string) error {
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerData.Sdp,
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		return err
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return err
	}

	if err := pc.SetLocalDescription(answer); err != nil {
		return err
	}

	// Generate session ID dynamically (based on clientID or another identifier)
	session := models.Session{
		SessionID: generateRandomID(20),              // Generating random session ID
		UserID:    utils.ParseStringToUint(clientID), // Ideally, use the actual user ID associated with the WebSocket connection
	}

	// Save session to DB
	if err := db.Create(&session).Error; err != nil {
		return err
	}

	log.Println("Session created successfully")

	answerJSON, err := json.Marshal(answer)

	if err != nil {
		return err
	}

	// Send the answer back
	answerMessage := Offer{
		Type: "answer",
		Data: OfferData{
			Type: "answer",
			Sdp:  string(answerJSON),
		},
	}

	// Send answer to WebSocket connection
	if err := conn.WriteJSON(answerMessage); err != nil {
		return err
	}

	return nil
}

// HandleAnswer processes the answer from a client.
func HandleAnswer(pc *webrtc.PeerConnection, answerData string) error {
	answer := webrtc.SessionDescription{}
	if err := json.Unmarshal([]byte(answerData), &answer); err != nil {
		return err
	}

	return pc.SetRemoteDescription(answer)
}

// HandleCandidate processes an ICE candidate from a client.
func HandleCandidate(pc *webrtc.PeerConnection, candidateData string) error {
	candidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal([]byte(candidateData), &candidate); err != nil {
		return err
	}

	return pc.AddICECandidate(candidate)
}

// generateRandomID creates a random alphanumeric string of given length.
func generateRandomID(length int) string {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))) // Create a local random generator
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[r.Intn(len(charset))]
	}
	return string(id)
}
