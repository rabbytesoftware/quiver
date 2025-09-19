package logstream

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/domain/events"
)

// SimulatedService manages simulated log streaming with filtering and level control
type SimulatedService struct {
	mu           sync.RWMutex
	filter       *regexp.Regexp
	level        string
	paused       bool
	validLevels  map[string]int
	currentLevel int
}

// NewSimulatedService creates a new simulated log streaming service
func NewSimulatedService() *SimulatedService {
	validLevels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	return &SimulatedService{
		level:        "debug", // Show all levels by default
		validLevels:  validLevels,
		currentLevel: 0,
	}
}

// Start begins streaming log lines in simulation mode
func (s *SimulatedService) Start(ctx context.Context) (<-chan events.LogLine, func()) {
	logChan := make(chan events.LogLine, 100)
	
	// Create a context for the streaming goroutine
	streamCtx, cancel := context.WithCancel(ctx)
	
	// Start the simulation goroutine
	go s.simulate(streamCtx, logChan)
	
	// Return the channel and stop function
	stopFunc := func() {
		cancel()
		close(logChan)
	}
	
	return logChan, stopFunc
}

// SetFilter sets a regex filter for log lines
func (s *SimulatedService) SetFilter(pattern string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if pattern == "" {
		s.filter = nil
		return nil
	}
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}
	
	s.filter = regex
	return nil
}

// SetLevel sets the minimum log level to display
func (s *SimulatedService) SetLevel(level string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	level = strings.ToLower(level)
	levelInt, exists := s.validLevels[level]
	if !exists {
		return fmt.Errorf("invalid log level: %s", level)
	}
	
	s.level = level
	s.currentLevel = levelInt
	return nil
}

// Pause pauses log streaming
func (s *SimulatedService) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paused = true
}

// Resume resumes log streaming
func (s *SimulatedService) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.paused = false
}

// IsPaused returns whether streaming is paused
func (s *SimulatedService) IsPaused() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.paused
}

// GetCurrentFilter returns the current filter pattern
func (s *SimulatedService) GetCurrentFilter() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.filter == nil {
		return ""
	}
	return s.filter.String()
}

// GetCurrentLevel returns the current log level
func (s *SimulatedService) GetCurrentLevel() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.level
}

// shouldIncludeLine determines if a log line should be included based on filters
func (s *SimulatedService) shouldIncludeLine(logLine events.LogLine) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Check level filter
	logLevelInt, exists := s.validLevels[strings.ToLower(logLine.Level)]
	if !exists || logLevelInt < s.currentLevel {
		return false
	}
	
	// Check regex filter
	if s.filter != nil && !s.filter.MatchString(logLine.Text) {
		return false
	}
	
	return true
}

// simulate generates simulated log lines for demonstration
func (s *SimulatedService) simulate(ctx context.Context, logChan chan<- events.LogLine) {
	ticker := time.NewTicker(200 * time.Millisecond) // ~5 Hz
	defer ticker.Stop()
	
	levels := []string{"debug", "info", "warn", "error"}
	messages := []string{
		"Processing request from client",
		"Database connection established",
		"Cache miss for key: user_session_%d",
		"Authentication successful for user: %s",
		"Rate limit exceeded for IP: 192.168.1.%d",
		"Garbage collection completed in %dms",
		"New connection established",
		"Request completed in %dms",
		"Memory usage: %dMB",
		"CPU usage: %d%%",
		"Disk space warning: %d%% full",
		"Service health check passed",
		"Configuration reloaded",
		"Background task completed",
		"Error connecting to external service",
		"Retrying failed operation",
		"Transaction rolled back",
		"Session expired for user",
		"File upload completed: %s",
		"Email notification sent",
	}
	
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if s.IsPaused() {
				continue
			}
			
			// Generate a random log line
			level := levels[rand.Intn(len(levels))]
			messageTemplate := messages[rand.Intn(len(messages))]
			
			// Fill in message template with random data
			var message string
			switch {
			case strings.Contains(messageTemplate, "%d") && strings.Contains(messageTemplate, "%s"):
				message = fmt.Sprintf(messageTemplate, rand.Intn(1000), fmt.Sprintf("user_%d", rand.Intn(100)))
			case strings.Contains(messageTemplate, "%s"):
				message = fmt.Sprintf(messageTemplate, fmt.Sprintf("file_%d.txt", rand.Intn(100)))
			case strings.Contains(messageTemplate, "%d"):
				message = fmt.Sprintf(messageTemplate, rand.Intn(1000))
			default:
				message = messageTemplate
			}
			
			logLine := events.LogLine{
				Text:  fmt.Sprintf("[%04d] %s", counter, message),
				Level: level,
				Time:  time.Now(),
			}
			
			// Apply filters
			if s.shouldIncludeLine(logLine) {
				select {
				case logChan <- logLine:
				case <-ctx.Done():
					return
				default:
					// Channel is full, skip this log line
				}
			}
			
			counter++
		}
	}
}
