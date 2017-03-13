package session

import (
	"net/http"
	"net/url"
	"sync"
)

type Store interface {
	New(string, string) *Session
	Get(*http.Request, string) (*Session, error)
	Save(http.ResponseWriter, *http.Request, *Session) error
}

type CookieStore struct {
	secret string
	mu     sync.Mutex
}

func NewCookieStore(secret string) *CookieStore {
	return &CookieStore{
		secret: secret,
		mu:     sync.Mutex{},
	}
}

func (store *CookieStore) Get(req *http.Request, name string) (*Session, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var session *Session
	cookie, err := req.Cookie(name)
	if err != nil || cookie.Value == "" {
		sid, err := sessionId(16)
		if err != nil {
			return nil, err
		}
		session = store.New(sid, name)
	} else {
		value, _ := url.QueryUnescape(cookie.Value)
		var values Values
		session = NewSession(name, store)
		if err := decodeCookie([]byte(store.secret), name, value, &values); err == nil {
			session.Values = values
		}
	}

	return session, nil
}

func (store *CookieStore) New(sid, name string) *Session {
	return NewSession(name, store)
}

func (store *CookieStore) Save(rw http.ResponseWriter, req *http.Request, session *Session) error {
	encoded, err := encodeCookie([]byte(store.secret), session.Name, session.Values)
	if err != nil {
		return err
	}

	http.SetCookie(rw, NewCookie(session.Name, url.QueryEscape(encoded), session.Options))
	return nil
}

type MemoryStore struct {
	mu       sync.Mutex
	sessions map[string]*Session
	secret   string
}

func NewMemoryStore(secret string) Store {
	return &MemoryStore{
		secret:   secret,
		mu:       sync.Mutex{},
		sessions: make(map[string]*Session),
	}
}

func (store *MemoryStore) Get(req *http.Request, name string) (*Session, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var session *Session
	cookie, err := req.Cookie(name)
	if err != nil || cookie.Value == "" {
		sid, err := sessionId(16)
		if err != nil {
			return nil, err
		}
		session = store.New(sid, name)
	} else {
		value, _ := url.QueryUnescape(cookie.Value)
		var sid string
		if err := decodeCookie([]byte(store.secret), name, value, &sid); err != nil {
			sid, err := sessionId(16)
			if err != nil {
				return nil, err
			}
			session = store.New(sid, name)
		} else {
			var ok bool
			if session, ok = store.sessions[sid]; !ok {
				session = store.New(sid, name)
			}
		}
	}

	return session, nil
}

func (store *MemoryStore) New(sid, name string) *Session {
	session := NewSession(name, store)
	session.Id = sid
	store.sessions[sid] = session

	return session
}

func (store *MemoryStore) Save(rw http.ResponseWriter, req *http.Request, session *Session) error {
	encoded, err := encodeCookie([]byte(store.secret), session.Name, session.Id)
	if err != nil {
		return err
	}

	http.SetCookie(rw, NewCookie(session.Name, url.QueryEscape(encoded), session.Options))
	return nil
}
