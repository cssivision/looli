package session

import (
	"net/http"
	"time"
)

type Values map[interface{}]interface{}

type Session struct {
	Name    string
	Values  Values
	Id      string
	store   Store
	Options *Options
}

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

var defaultOptions = &Options{
	Path:     "/",
	MaxAge:   86400 * 30,
	HttpOnly: true,
}

func NewSession(name string, store Store) (session *Session) {
	session = &Session{
		Name:   name,
		Values: make(map[interface{}]interface{}),
		store:  store,
	}

	return
}

func (session *Session) Save(rw http.ResponseWriter, req *http.Request) error {
	return session.store.Save(rw, req, session)
}

func NewCookie(name, value string, options *Options) *http.Cookie {
	if options == nil {
		options = defaultOptions
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}

	if options.MaxAge > 0 {
		d := time.Duration(options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if options.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}
