package session

import (
	"crypto/rand"
	"encoding/gob"
	"github.com/cssivision/looli"
	"sync"
)

func init() {
	gob.Register([]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[int]string{})
	gob.Register(map[int]int{})
	gob.Register(map[int]int64{})
}

// Sessions is a global session manager
type Sessions struct {
	// Store use to store session, it can be any store that implement the Store interface.
	Store Store

	// use mu to avoid data race.
	mu sync.Mutex
}

// generate unique session id for per request
func generateRandomKey(length int) ([]byte, error) {
	b := make([]byte, length)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return nil, err
	}
	return b, nil
}

// NewSessions create a global sessions manager, default store is Cookie store.
//
// secret used to sign the cookie value.
//
// The aesKey argument should be the AES key, either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256,
// if aesKey is empty string, the cookie value will not Encrypted.
func NewSessions(secret string, aesKey string) *Sessions {
	if secret == "" {
		panic("secret key can not be empty")
	}

	return &Sessions{
		mu:    sync.Mutex{},
		Store: NewCookieStore(secret, aesKey),
	}
}

// Get return a existed session or create a new session
func (sm *Sessions) Get(ctx *looli.Context, name string) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	req := ctx.Request

	session, err := sm.Store.Get(req, name)
	return session, err
}
