package session

import (
	"github.com/cssivision/looli"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {
	length := 32
	b, err := generateRandomKey(length)
	assert.Nil(t, err)
	assert.Equal(t, length, len(b))
}

func TestNewSessions(t *testing.T) {
	assert.Panics(t, func() {
		NewSessions("", "")
	})

	secret := "secret"
	aesKey := "secret"
	sessions := NewSessions("secret", "secret")
	assert.NotNil(t, &sessions.mu)
	assert.NotNil(t, sessions.Store)
	store := sessions.Store.(*CookieStore)
	assert.Equal(t, store.secret, secret)
	assert.Equal(t, store.aesKey, aesKey)
}

func TestSessionsGet(t *testing.T) {
	sessions := NewSessions("secret", "")
	ctx := &looli.Context{}
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	ctx.Request = req
	session, err := sessions.Get(ctx, "sess")
	assert.Nil(t, err)
	assert.NotNil(t, session)
}
