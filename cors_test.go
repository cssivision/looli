package looli

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithoutOptions(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := New()
	router.Use(Cors(CorsOption{}))

	router.Get("/a", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestAllOrigins(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := New()
	router.Use(Cors(CorsOption{
		AllowOrigins: []string{"*"},
	}))

	router.Get("/a", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestAllowedOrigins(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	notAllowedOrigin := "looli.com"
	router := New()
	router.Use(Cors(CorsOption{
		AllowOrigins: []string{origin},
	}))

	router.Get("/a", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	// test allowed origin
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()

	// test not allowed origin
	getReq, err = http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", notAllowedOrigin)

	resp, err = http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Empty(t, resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("Origin: %v is not allowed", notAllowedOrigin), string(bodyBytes))
	resp.Body.Close()
}
