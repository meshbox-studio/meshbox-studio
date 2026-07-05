package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]User
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]User),
	}
}

func (s *SessionStore) Create(user User) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := generateToken()
	s.sessions[token] = user
	return token
}

func (s *SessionStore) Get(token string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.sessions[token]
	return user, ok
}

func generateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}