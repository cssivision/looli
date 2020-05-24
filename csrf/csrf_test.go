package csrf

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cssivision/looli"
	"github.com/stretchr/testify/assert"
)

func TestCsrfMiddleWare(t *testing.T) {
	t.Run("without csrf token", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		router := looli.New()

		router.Use(Default())
		router.Get("/", func(ctx *looli.Context) {
			ctx.Status(statusCode)
			ctx.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL

		resp, err := http.Get(serverURL)
		assert.Nil(t, err)
		assert.Equal(t, resp.Header["Vary"][0], "Cookie")

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, invalidCsrfTokenResponse, string(bodyBytes))
	})

	t.Run("with token", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		router := looli.New()

		router.Use(New(Options{
			Skip: func(ctx *looli.Context) bool {
				if ctx.Method == http.MethodGet {
					return true
				}

				return false
			},
		}))

		router.Get("/", func(ctx *looli.Context) {
			token := NewToken(ctx)
			ctx.String(token)
		})

		router.Post("/", func(ctx *looli.Context) {
			ctx.Status(statusCode)
			ctx.String(serverResponse)
		})

		server := httptest.NewServer(router)

		serverURL := server.URL
		resp, err := http.Get(serverURL)
		assert.Nil(t, err)

		secret := resp.Cookies()[0].Value
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		token := string(bodyBytes)
		assert.NotEqual(t, strings.Index(token, "."), -1)

		postReq, err := http.NewRequest(http.MethodPost, serverURL, nil)
		assert.Nil(t, err)

		postReq.Header.Set(headerKey, token)
		postReq.AddCookie(&http.Cookie{
			Name:  cookieName,
			Value: secret,
		})

		resp, err = http.DefaultClient.Do(postReq)
		assert.Nil(t, err)
		assert.Equal(t, statusCode, resp.StatusCode)

		bodyBytes, err = ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		assert.Equal(t, serverResponse, string(bodyBytes))

		resp.Body.Close()
		server.Close()
	})
}
