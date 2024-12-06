package webrtc

import (
	"log"
	"sync"

	"github.com/pion/webrtc/v4"
)

// SFU is the main structure for our Selective Forwarding Unit (SFU).
type SFU struct {
	mu              sync.Mutex
	peerConnections map[string]*webrtc.PeerConnection
	videoTracks     map[string]*webrtc.TrackLocalStaticSample
	audioTracks     map[string]*webrtc.TrackLocalStaticSample
}

// NewSFU creates a new instance of the SFU.
func NewSFU() *SFU {
	return &SFU{
		peerConnections: make(map[string]*webrtc.PeerConnection),
		videoTracks:     make(map[string]*webrtc.TrackLocalStaticSample),
		audioTracks:     make(map[string]*webrtc.TrackLocalStaticSample),
	}
}

// AddTrack adds a track to the SFU and forwards it to the right clients.
func (s *SFU) AddTrack(clientID string, pc *webrtc.PeerConnection, track *webrtc.TrackLocalStaticSample) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Add track to the map
	if track.Kind() == webrtc.RTPCodecTypeVideo {
		s.videoTracks[clientID] = track
	} else if track.Kind() == webrtc.RTPCodecTypeAudio {
		s.audioTracks[clientID] = track
	}

	// Forward video and audio tracks to other participants
	for otherClientID, otherPC := range s.peerConnections {
		if otherClientID != clientID {
			if track.Kind() == webrtc.RTPCodecTypeVideo {
				_, err := otherPC.AddTrack(track)
				if err != nil {
					log.Println("Error forwarding video track:", err)
					return err
				}
			} else if track.Kind() == webrtc.RTPCodecTypeAudio {
				_, err := otherPC.AddTrack(track)
				if err != nil {
					log.Println("Error forwarding audio track:", err)
					return err
				}
			}
		}
	}
	return nil
}

// AddPeerConnection adds a new peer connection to the SFU.
func (s *SFU) AddPeerConnection(clientID string, pc *webrtc.PeerConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peerConnections[clientID] = pc
}

// RemovePeerConnection removes a peer connection from the SFU.
func (s *SFU) RemovePeerConnection(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.peerConnections, clientID)
}

// BroadcastICECandidates broadcasts an ICE candidate to all connected peers.
func (s *SFU) BroadcastICECandidates(candidate webrtc.ICECandidateInit) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Loop over all peer connections and add the ICE candidate
	for _, pc := range s.peerConnections {
		if candidate.Candidate != "" {
			// Adding the candidate to the peer connection
			err := pc.AddICECandidate(candidate)
			if err != nil {
				log.Printf("Error adding ICE candidate to peer connection %v: %v", pc, err)
			}
		}
	}
}
