package store

import (
	"errors"
	"sync"
	"time"
)

// The interface
type UserStore interface {
	SessionStore(sessionID string, userID int) error
	GetSession(SessionID string) (int, error)
}

// Session represents an active user session
type Session struct {
	Username  string
	Role      string
	ExpiresAt time.Time
}

// SessionStore provides concurrency-safe storage for active sessions
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

// NewSessionStore initializes a new thread-safe session store
func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]Session),
	}
}

// Set stores or update a session safely using a write lock
func (s *SessionStore) Set(token string, session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = session
}

// Get retrieves a session safely using a read lock
func (s *SessionStore) Get(token string) (Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return Session{}, errors.New("session not found")
	}
	return session, nil
}

// Delete removes a session safely using a write lock
func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}
