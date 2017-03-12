package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

var defaultSessionCookieName = "sess"

type SessionManager struct {
	CookieName string
	Store      Store
	mu         sync.Mutex
}

func (sm *SessionManager) sessionId() (string, error) {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (sm *SessionManager) Get(req *http.Request) (session *Session, err error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	cookie, err := req.Cookie(sm.CookieName)
	if err != nil || cookie.Value == "" {
		_, err := sm.sessionId()
		if err != nil {
			return nil, err
		}
	} else {
		// sm.Store.Get()
	}
	return
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
	// mu:         sync.Mutex{},
	// CookieName: defaultSessionCookieName,
	// Store:      NewMemoryStore(),
	}
}
