package playback

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// PlaybackClient replays recorded client requests to stdout
type PlaybackClient struct {
	session  *PlaybackSession
	messages []json.RawMessage
	delay    time.Duration
}

// NewPlaybackClient creates a new playback client
func NewPlaybackClient(session *PlaybackSession) *PlaybackClient {
	clientMessages := session.GetClientMessages()
	messages := make([]json.RawMessage, len(clientMessages))
	
	for i, msg := range clientMessages {
		messages[i] = msg.Message
	}
	
	return &PlaybackClient{
		session:  session,
		messages: messages,
		delay:    100 * time.Millisecond, // Small delay between messages
	}
}

// SetDelay sets the delay between messages
func (c *PlaybackClient) SetDelay(delay time.Duration) {
	c.delay = delay
}

// Run starts the playback client
func (c *PlaybackClient) Run() error {
	log.Printf("Starting playback client with %d messages", len(c.messages))
	
	// Wait for server to be ready by reading from stdin
	scanner := bufio.NewScanner(os.Stdin)
	messageIndex := 0
	
	for scanner.Scan() {
		serverResponse := scanner.Text()
		
		// Log server response (to stderr so it doesn't interfere with stdout)
		log.Printf("Server response: %s", serverResponse)
		
		// Send next client request if available
		if messageIndex < len(c.messages) {
			time.Sleep(c.delay)
			
			// Send message to stdout (which goes to server's stdin)
			fmt.Println(string(c.messages[messageIndex]))
			log.Printf("Sent client request %d/%d", messageIndex+1, len(c.messages))
			
			messageIndex++
		} else {
			log.Printf("All messages sent, exiting")
			break
		}
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading server responses: %w", err)
	}
	
	log.Printf("Playback client finished")
	return nil
}

// RunBatch sends all messages without waiting for responses (for testing)
func (c *PlaybackClient) RunBatch() error {
	log.Printf("Starting batch playback with %d messages", len(c.messages))
	
	for i, message := range c.messages {
		fmt.Println(string(message))
		log.Printf("Sent message %d/%d", i+1, len(c.messages))
		
		if i < len(c.messages)-1 {
			time.Sleep(c.delay)
		}
	}
	
	log.Printf("Batch playback finished")
	return nil
}