package session

import (
	"net/http"
	"net/url"
	"sync"
)

type Store interface {
	New(string, string) *Session
	Get(string, string) *Session
	Save(http.ResponseWriter, *http.Request, *Session)
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
	return &MemoryStore{
		mu:       sync.Mutex{},
		sessions: make(map[string]*Session),
	}
}

func (store *MemoryStore) Get(sid, name string) *Session {
	store.mu.Lock()
	defer store.mu.Unlock()

	session, ok := store.sessions[sid]
	if !ok {
		session = store.New(sid, name)
		session.store = store
	}

	return session
}

func (store *MemoryStore) New(sid, name string) *Session {
	session := NewSession(sid, name, store)
	store.sessions[sid] = session

	return session
}

func (store *MemoryStore) Save(rw http.ResponseWriter, req *http.Request, session *Session) {
	http.SetCookie(rw, NewCookie(session.Name, url.QueryEscape(session.Id), session.Options))
}
