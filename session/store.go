package session

import (
	"net/http"
	"sync"
)

type Store interface {
	New() *Session
	Get(sessionID string) (*Session, error)
	Save(http.ResponseWriter, *http.Request) error
}

type CookieStore struct {
	mu sync.Mutex
}

type MemoryStore struct {
	mu       sync.Mutex
	sessions map[string]*Session
}

type FileSystemStore struct {
}

type RedisStore struct {
}

func NewMemoryStore() Store {
	return &MemoryStore{}
}

func (ms *MemoryStore) Get(sessionId string) (session *Session, err error) {
	return
}

func (ms *MemoryStore) New() (session *Session) {
	return
}

func (ms *MemoryStore) Save(rw http.ResponseWriter, req *http.Request) (err error) {
	return
}
