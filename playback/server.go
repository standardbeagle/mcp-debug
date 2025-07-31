package playback

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// PlaybackServer replays recorded server responses
type PlaybackServer struct {
	session   *PlaybackSession
	responses []json.RawMessage
	delay     time.Duration
}

// NewPlaybackServer creates a new playback server
func NewPlaybackServer(session *PlaybackSession) *PlaybackServer {
	serverMessages := session.GetServerMessages()
	responses := make([]json.RawMessage, len(serverMessages))
	
	for i, msg := range serverMessages {
		responses[i] = msg.Message
	}
	
	return &PlaybackServer{
		session:   session,
		responses: responses,
		delay:     50 * time.Millisecond, // Small delay before responding
	}
}

// SetDelay sets the delay before sending responses
func (s *PlaybackServer) SetDelay(delay time.Duration) {
	s.delay = delay
}

// Run starts the playback server
func (s *PlaybackServer) Run() error {
	log.Printf("Starting playback server with %d responses", len(s.responses))
	
	scanner := bufio.NewScanner(os.Stdin)
	responseIndex := 0
	
	for scanner.Scan() {
		clientRequest := scanner.Text()
		
		// Log client request (to stderr)
		log.Printf("Client request: %s", clientRequest)
		
		// Send corresponding server response if available
		if responseIndex < len(s.responses) {
			time.Sleep(s.delay)
			
			// Send response to stdout (which goes to client)
			fmt.Println(string(s.responses[responseIndex]))
			log.Printf("Sent server response %d/%d", responseIndex+1, len(s.responses))
			
			responseIndex++
		} else {
			// If no more responses, send a generic error
			errorResponse := map[string]interface{}{
				"jsonrpc": "2.0",
				"error": map[string]interface{}{
					"code":    -32000,
					"message": "No more recorded responses available",
				},
				"id": nil,
			}
			
			errorBytes, _ := json.Marshal(errorResponse)
			fmt.Println(string(errorBytes))
			log.Printf("Sent generic error response (no more recorded responses)")
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading client requests: %w", err)
	}
	
	log.Printf("Playback server finished")
	return nil
}

// RunStateless starts the server without maintaining request-response pairing
// Useful for testing where request order might differ
func (s *PlaybackServer) RunStateless() error {
	log.Printf("Starting stateless playback server")
	
	scanner := bufio.NewScanner(os.Stdin)
	responseIndex := 0
	
	for scanner.Scan() {
		clientRequest := scanner.Text()
		log.Printf("Client request: %s", clientRequest)
		
		// Always cycle through responses
		if len(s.responses) > 0 {
			time.Sleep(s.delay)
			
			response := s.responses[responseIndex%len(s.responses)]
			fmt.Println(string(response))
			log.Printf("Sent cycled response %d (index %d)", responseIndex+1, responseIndex%len(s.responses))
			
			responseIndex++
		}
	}
	
	return nil
}