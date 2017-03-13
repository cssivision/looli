package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

// Sessions is a global session manager
type Sessions struct {
	// Store use to store session
	Store Store
	mu    sync.Mutex
}

// generate unique session id for per request
func sessionId(length int) (string, error) {
	b := make([]byte, length)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (sm *Sessions) Get(req *http.Request, name string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, err := sm.Store.Get(req, name)
	return session, err
}

func NewSessions(secret string) *Sessions {
	if secret == "" {
		panic("secret key can not be empty")
	}

	return &Sessions{
		mu:    sync.Mutex{},
		Store: NewCookieStore(secret),
	}
}
