package session

import (
	"github.com/cssivision/looli"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCookieStore(t *testing.T) {
	secret := "secret"
	aesKey := "secret"
	store := NewCookieStore(secret, aesKey)
	assert.Equal(t, secret, store.secret)
	assert.Equal(t, aesKey, store.aesKey)
	assert.NotNil(t, &store.mu)
}

func TestCookieStoreNew(t *testing.T) {
	secret := "secret"
	aesKey := "secret"
	store := NewCookieStore(secret, aesKey)
	sid, _ := generateRandomKey(16)
	name := "sess"
	session := store.New(string(sid), name)
	assert.NotNil(t, session.Values)
	assert.Equal(t, session.Name, name)
	assert.Equal(t, session.Id, string(sid))
}

func TestCookieStoreGetSave(t *testing.T) {
	t.Run("with cookie", func(t *testing.T) {
		router := looli.New()
		secret := "secret"
		aesKey := ""
		serverResponse := "sever response"
		sessions := NewSessions(secret, aesKey)
		router.Get("/", func(ctx *looli.Context) {
			sess, err := sessions.Get(ctx, "sess")
			assert.Nil(t, err)
			assert.NotNil(t, sess)

			sess.Values["name"] = "cssivision"
			assert.Nil(t, sess.Save(ctx))
			ctx.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL

		req, err := http.NewRequest(http.MethodGet, serverURL, nil)
		assert.Nil(t, err)
		req.Header.Set("Cookie", "sess=cookie")
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("without cookie", func(t *testing.T) {
		router := looli.New()
		secret := "secret"
		aesKey := ""
		serverResponse := "sever response"
		sessions := NewSessions(secret, aesKey)
		router.Get("/", func(ctx *looli.Context) {
			sess, err := sessions.Get(ctx, "sess")
			assert.Nil(t, err)
			assert.NotNil(t, sess)

			sess.Values["name"] = "cssivision"
			assert.Nil(t, sess.Save(ctx))
			ctx.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL

		req, err := http.NewRequest(http.MethodGet, serverURL, nil)
		assert.Nil(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("with aes key", func(t *testing.T) {
		router := looli.New()
		secret := "secret"
		aesKey := "1111111111111111"
		serverResponse := "sever response"
		sessions := NewSessions(secret, aesKey)
		router.Get("/", func(ctx *looli.Context) {
			sess, err := sessions.Get(ctx, "sess")
			assert.Nil(t, err)
			assert.NotNil(t, sess)

			sess.Values["name"] = "cssivision"
			assert.Nil(t, sess.Save(ctx))
			ctx.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		req, err := http.NewRequest(http.MethodGet, serverURL, nil)
		assert.Nil(t, err)
		req.Header.Set("Cookie", "sess=cookie")

		resp, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})
}
