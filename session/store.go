package session

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
)

type Store interface {
	New(string, string) *Session
	Get(*http.Request, string) (*Session, error)
	Save(http.ResponseWriter, *http.Request, *Session) error
}

// ------------------------------- CookieStore -------------------------------

type CookieStore struct {
	secret string
	aesKey string
	mu     sync.Mutex
}

func NewCookieStore(secret, aesKey string) *CookieStore {
	return &CookieStore{
		secret: secret,
		aesKey: aesKey,
		mu:     sync.Mutex{},
	}
}

func (store *CookieStore) Get(req *http.Request, name string) (*Session, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var session *Session
	cookie, err := req.Cookie(name)
	if err != nil || cookie.Value == "" {
		sid, err := generateRandomKey(16)
		if err != nil {
			return nil, err
		}
		session = store.New(string(sid), name)
	} else {
		value, _ := url.QueryUnescape(cookie.Value)
		var values Values
		session = NewSession(name, store)
		err := DecodeCookie([]byte(store.secret), []byte(store.aesKey), name, value, &values)
		if err != nil {
			return nil, errors.New("decode session error")
		}
		session.Values = values
	}

	return session, nil
}

func (store *CookieStore) New(sid, name string) *Session {
	session := NewSession(name, store)
	session.Id = sid
	return session
}

func (store *CookieStore) Save(rw http.ResponseWriter, req *http.Request, session *Session) error {
	encoded, err := EncodeCookie([]byte(store.secret), []byte(store.aesKey), session.Name, session.Values)
	if err != nil {
		return err
	}

	http.SetCookie(rw, NewCookie(session.Name, url.QueryEscape(encoded), session.Options))
	return nil
}

// -------------------------------  MemoryStore  -------------------------------

// type MemoryStore struct {
// 	mu       sync.Mutex
// 	sessions map[string]*Session
// 	secret   string
// 	aesKey   string
// }

// func NewMemoryStore(secret, aesKey string) Store {
// 	return &MemoryStore{
// 		secret:   secret,
// 		aesKey:   aesKey,
// 		mu:       sync.Mutex{},
// 		sessions: make(map[string]*Session),
// 	}
// }

// func (store *MemoryStore) Get(req *http.Request, name string) (*Session, error) {
// 	store.mu.Lock()
// 	defer store.mu.Unlock()

// 	var session *Session
// 	cookie, err := req.Cookie(name)
// 	if err != nil || cookie.Value == "" {
// 		sid, err := generateRandomKey(16)
// 		if err != nil {
// 			return nil, err
// 		}
// 		session = store.New(string(sid), name)
// 	} else {
// 		value, _ := url.QueryUnescape(cookie.Value)
// 		var sid string
// 		if err := DecodeCookie([]byte(store.secret), []byte(store.aesKey), name, value, &sid); err != nil {
// 			sid, err := generateRandomKey(16)
// 			if err != nil {
// 				return nil, err
// 			}
// 			session = store.New(string(sid), name)
// 		} else {
// 			var ok bool
// 			if session, ok = store.sessions[sid]; !ok {
// 				session = store.New(sid, name)
// 			}
// 		}
// 	}

// 	return session, nil
// }

// func (store *MemoryStore) New(sid, name string) *Session {
// 	session := NewSession(name, store)
// 	session.Id = sid
// 	store.sessions[sid] = session

// 	return session
// }

// func (store *MemoryStore) Save(rw http.ResponseWriter, req *http.Request, session *Session) error {
// 	encoded, err := EncodeCookie([]byte(store.secret), []byte(store.aesKey), session.Name, session.Id)
// 	if err != nil {
// 		return err
// 	}

// 	http.SetCookie(rw, NewCookie(session.Name, url.QueryEscape(encoded), session.Options))
// 	return nil
// }
