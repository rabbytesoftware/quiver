package pool

import (
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewPool(t *testing.T) {
	pool := NewPool()

	if pool == nil {
		t.Fatal("NewPool() returned nil")
	}

	if pool.messages == nil {
		t.Error("messages channel should be initialized")
	}

	if pool.subscribers == nil {
		t.Error("subscribers slice should be initialized")
	}

	if len(pool.subscribers) != 0 {
		t.Errorf("Expected empty subscribers slice, got length %d", len(pool.subscribers))
	}
}

func TestPool_GetSubscriberCount(t *testing.T) {
	pool := NewPool()

	// Test initial count
	if count := pool.GetSubscriberCount(); count != 0 {
		t.Errorf("Expected initial subscriber count to be 0, got %d", count)
	}

	// Add subscribers
	subscriber1 := func(level logrus.Level, message string) {}
	subscriber2 := func(level logrus.Level, message string) {}

	pool.Subscribe(subscriber1)
	if count := pool.GetSubscriberCount(); count != 1 {
		t.Errorf("Expected subscriber count to be 1, got %d", count)
	}

	pool.Subscribe(subscriber2)
	if count := pool.GetSubscriberCount(); count != 2 {
		t.Errorf("Expected subscriber count to be 2, got %d", count)
	}
}

func TestPool_Subscribe(t *testing.T) {
	pool := NewPool()

	var called bool
	var mu sync.Mutex
	subscriber := func(level logrus.Level, message string) {
		mu.Lock()
		called = true
		mu.Unlock()
	}

	// Test subscribing
	pool.Subscribe(subscriber)

	if pool.GetSubscriberCount() != 1 {
		t.Error("Expected subscriber to be added")
	}

	// Test that subscriber receives messages
	message := Message{
		Level:   logrus.InfoLevel,
		Message: "test message",
	}

	pool.AddMessage(message)

	// Give goroutine time to execute
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	wasCalled := called
	mu.Unlock()

	if !wasCalled {
		t.Error("Expected subscriber to be called")
	}
}

func TestPool_AddMessage(t *testing.T) {
	pool := NewPool()

	var receivedLevel logrus.Level
	var receivedMessage string
	var callCount int
	var mu sync.Mutex

	subscriber := func(level logrus.Level, message string) {
		mu.Lock()
		defer mu.Unlock()
		receivedLevel = level
		receivedMessage = message
		callCount++
	}

	pool.Subscribe(subscriber)

	// Test adding a message
	testMessage := Message{
		Level:   logrus.WarnLevel,
		Message: "warning message",
	}

	pool.AddMessage(testMessage)

	// Give goroutine time to execute
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if callCount != 1 {
		t.Errorf("Expected subscriber to be called once, got %d calls", callCount)
	}

	if receivedLevel != logrus.WarnLevel {
		t.Errorf("Expected level %v, got %v", logrus.WarnLevel, receivedLevel)
	}

	if receivedMessage != "warning message" {
		t.Errorf("Expected message 'warning message', got %q", receivedMessage)
	}
	mu.Unlock()
}

func TestPool_MultipleSubscribers(t *testing.T) {
	pool := NewPool()

	var calls1, calls2 int
	var mu sync.Mutex

	subscriber1 := func(level logrus.Level, message string) {
		mu.Lock()
		calls1++
		mu.Unlock()
	}

	subscriber2 := func(level logrus.Level, message string) {
		mu.Lock()
		calls2++
		mu.Unlock()
	}

	pool.Subscribe(subscriber1)
	pool.Subscribe(subscriber2)

	message := Message{
		Level:   logrus.ErrorLevel,
		Message: "error message",
	}

	pool.AddMessage(message)

	// Give goroutines time to execute
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	if calls1 != 1 {
		t.Errorf("Expected subscriber1 to be called once, got %d calls", calls1)
	}

	if calls2 != 1 {
		t.Errorf("Expected subscriber2 to be called once, got %d calls", calls2)
	}
	mu.Unlock()
}

func TestPool_ChannelBuffering(t *testing.T) {
	pool := NewPool()

	// Test that channel has buffer (should not block)
	for i := 0; i < 128; i++ {
		message := Message{
			Level:   logrus.InfoLevel,
			Message: "buffered message",
		}
		pool.AddMessage(message)
	}

	// This should not block since we have a buffered channel
	// If it blocks, the test will timeout
}

func TestPool_ChannelOverflow(t *testing.T) {
	pool := NewPool()

	// Fill the channel buffer
	for i := 0; i < 128; i++ {
		message := Message{
			Level:   logrus.InfoLevel,
			Message: "buffer message",
		}
		pool.AddMessage(message)
	}

	// This should not block due to the default case in select
	overflowMessage := Message{
		Level:   logrus.ErrorLevel,
		Message: "overflow message",
	}
	pool.AddMessage(overflowMessage)

	// Test passes if we reach here without blocking
}

func TestPool_ConcurrentAccess(t *testing.T) {
	pool := NewPool()

	var wg sync.WaitGroup
	var totalCalls int
	var mu sync.Mutex

	// Add multiple subscribers concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subscriber := func(level logrus.Level, message string) {
				mu.Lock()
				totalCalls++
				mu.Unlock()
			}
			pool.Subscribe(subscriber)
		}()
	}

	wg.Wait()

	if pool.GetSubscriberCount() != 10 {
		t.Errorf("Expected 10 subscribers, got %d", pool.GetSubscriberCount())
	}

	// Send messages concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			message := Message{
				Level:   logrus.InfoLevel,
				Message: "concurrent message",
			}
			pool.AddMessage(message)
		}(i)
	}

	wg.Wait()

	// Give goroutines time to execute
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	expectedCalls := 10 * 5 // 10 subscribers * 5 messages
	if totalCalls != expectedCalls {
		t.Errorf("Expected %d total calls, got %d", expectedCalls, totalCalls)
	}
	mu.Unlock()
}

func TestPool_MessageTypes(t *testing.T) {
	pool := NewPool()

	var receivedMessages []Message
	var mu sync.Mutex

	subscriber := func(level logrus.Level, message string) {
		mu.Lock()
		receivedMessages = append(receivedMessages, Message{Level: level, Message: message})
		mu.Unlock()
	}

	pool.Subscribe(subscriber)

	// Test different log levels
	testMessages := []Message{
		{Level: logrus.DebugLevel, Message: "debug message"},
		{Level: logrus.InfoLevel, Message: "info message"},
		{Level: logrus.WarnLevel, Message: "warn message"},
		{Level: logrus.ErrorLevel, Message: "error message"},
		{Level: logrus.FatalLevel, Message: "fatal message"},
		{Level: logrus.PanicLevel, Message: "panic message"},
	}

	for _, msg := range testMessages {
		pool.AddMessage(msg)
	}

	// Give goroutines time to execute
	time.Sleep(30 * time.Millisecond)

	mu.Lock()
	if len(receivedMessages) != len(testMessages) {
		t.Errorf("Expected %d messages, got %d", len(testMessages), len(receivedMessages))
	}

	// Note: Order is not guaranteed due to goroutines, so we just check that all messages were received
	for _, expected := range testMessages {
		found := false
		for _, received := range receivedMessages {
			if received.Level == expected.Level && received.Message == expected.Message {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected message not found: %v", expected)
		}
	}
	mu.Unlock()
}

func TestPool_EmptyMessage(t *testing.T) {
	pool := NewPool()

	var receivedLevel logrus.Level
	var receivedMessage string
	var called bool
	var mu sync.Mutex

	subscriber := func(level logrus.Level, message string) {
		mu.Lock()
		receivedLevel = level
		receivedMessage = message
		called = true
		mu.Unlock()
	}

	pool.Subscribe(subscriber)

	// Test empty message
	emptyMessage := Message{
		Level:   logrus.InfoLevel,
		Message: "",
	}

	pool.AddMessage(emptyMessage)

	// Give goroutine time to execute
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	wasCalled := called
	level := receivedLevel
	message := receivedMessage
	mu.Unlock()

	if !wasCalled {
		t.Error("Expected subscriber to be called for empty message")
	}

	if level != logrus.InfoLevel {
		t.Errorf("Expected level %v, got %v", logrus.InfoLevel, level)
	}

	if message != "" {
		t.Errorf("Expected empty message, got %q", message)
	}
}

func TestMessage_Struct(t *testing.T) {
	// Test Message struct creation and field access
	message := Message{
		Level:   logrus.WarnLevel,
		Message: "test message",
	}

	if message.Level != logrus.WarnLevel {
		t.Errorf("Expected Level to be %v, got %v", logrus.WarnLevel, message.Level)
	}

	if message.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %q", message.Message)
	}
}

func TestSubscriber_Type(t *testing.T) {
	// Test that Subscriber type works as expected
	var subscriber Subscriber = func(level logrus.Level, message string) {
		// Test function signature
		_ = level
		_ = message
	}

	// Test that we can call the subscriber
	subscriber(logrus.InfoLevel, "test")

	// Test passes if no compilation errors and no panics
}
