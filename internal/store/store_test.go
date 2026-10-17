package store

import (
	"sync"
	"testing"
	"time"
)

// Chack basic Set, Get, and Delete behavior
func TestStoreBasicOperations(t *testing.T) {
	s := NewSessionStore()
	token := "Secret-token-123"
	session := Session{
		Username:  "admin",
		Role:      "Administrator",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Test Set and Get
	s.Set(token, session)
	retrieved, err := s.Get(token)
	if err != nil {
		t.Fatalf("expected no error getting session, got %v", err)
	}
	if retrieved.Username != session.Username {
		t.Fatalf("expected username %s, got %s", session.Username, retrieved.Username)
	}

	// Test delete
	s.Delete(token)
	_, err = s.Get(token)
	if err == nil {
		t.Fatal("expected error getting delete session, got nil")
	}
}

// Stresses the store with parallel reads and writes to detect race conditions
func TestConcurrentAccess(t *testing.T) {
	s := NewSessionStore()
	var wg sync.WaitGroup
	iterations := 500

	// Concurrent writers
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			token := "token-shared"
			s.Set(token, Session{
				Username:  "concurrent-user",
				Role:      "Operator",
				ExpiresAt: time.Now().Add(time.Hour),
			})
		}(i)
	}

	// Concurrent readers
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.Get("token-shared")
		}()
	}
	wg.Wait()
}
