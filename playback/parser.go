package playback

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"mcp-debug/integration"
)

// PlaybackSession represents a parsed recording session
type PlaybackSession struct {
	StartTime  time.Time                        `json:"start_time"`
	ServerInfo string                           `json:"server_info"`
	Messages   []integration.RecordedMessage    `json:"messages"`
}

// ParseRecordingFile parses a recorded session file
func ParseRecordingFile(filename string) (*PlaybackSession, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open recording file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var session *PlaybackSession
	var messages []integration.RecordedMessage

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Try to parse as session header first
		if session == nil {
			var tempSession PlaybackSession
			if err := json.Unmarshal([]byte(line), &tempSession); err == nil {
				session = &tempSession
				continue
			}
		}

		// Parse as recorded message
		var message integration.RecordedMessage
		if err := json.Unmarshal([]byte(line), &message); err != nil {
			// Skip invalid lines but continue parsing
			continue
		}

		messages = append(messages, message)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if session == nil {
		// Create a default session if header not found
		session = &PlaybackSession{
			StartTime:  time.Now(),
			ServerInfo: "Unknown",
		}
	}

	session.Messages = messages
	return session, nil
}

// GetClientMessages returns only the client request messages
func (s *PlaybackSession) GetClientMessages() []integration.RecordedMessage {
	var clientMessages []integration.RecordedMessage
	for _, message := range s.Messages {
		if message.Direction == "request" {
			clientMessages = append(clientMessages, message)
		}
	}
	return clientMessages
}

// GetServerMessages returns only the server response messages
func (s *PlaybackSession) GetServerMessages() []integration.RecordedMessage {
	var serverMessages []integration.RecordedMessage
	for _, message := range s.Messages {
		if message.Direction == "response" {
			serverMessages = append(serverMessages, message)
		}
	}
	return serverMessages
}

// GetMessagePairs returns request-response pairs
func (s *PlaybackSession) GetMessagePairs() []MessagePair {
	var pairs []MessagePair
	var currentRequest *integration.RecordedMessage

	for _, message := range s.Messages {
		if message.Direction == "request" {
			currentRequest = &message
		} else if message.Direction == "response" && currentRequest != nil {
			pairs = append(pairs, MessagePair{
				Request:  *currentRequest,
				Response: message,
			})
			currentRequest = nil
		}
	}

	return pairs
}

// MessagePair represents a request-response pair
type MessagePair struct {
	Request  integration.RecordedMessage
	Response integration.RecordedMessage
}