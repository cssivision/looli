package session

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"sync"
)

type SessionManager struct {
	Store Store
	mu    sync.Mutex
}

func (sm *SessionManager) sessionId() (string, error) {
	b := make([]byte, 16)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (sm *SessionManager) Get(req *http.Request, name string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var session *Session
	cookie, err := req.Cookie(name)
	if err != nil || cookie.Value == "" {
		sid, err := sm.sessionId()
		if err != nil {
			return nil, err
		}
		session = sm.Store.New(sid, name)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session = sm.Store.Get(sid, name)
	}
	return session, nil
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		mu:    sync.Mutex{},
		Store: NewMemoryStore(),
	}
}
