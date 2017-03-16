package cors

import (
	"fmt"
	"github.com/cssivision/looli"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestEmptyHost(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := looli.New()
	router.Use(New(Options{}))

	router.Get("/a", func(c *looli.Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", "")

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Empty(t, resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestWithoutOptions(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{}))

	router.Get("/a", func(c *looli.Context) {
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
	router := looli.New()
	router.Use(New(Options{
		AllowOrigins: []string{"*"},
	}))

	router.Get("/a", func(c *looli.Context) {
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
	router := looli.New()
	router.Use(New(Options{
		AllowOrigins: []string{origin},
	}))

	router.Get("/a", func(c *looli.Context) {
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

func TestAllowOriginsFunc(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, "looli")
		},
	}))

	router.Get("/a", func(c *looli.Context) {
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
	getReq.Header.Set("Origin", "github.com")

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
	assert.Equal(t, "Origin: github.com is not allowed", string(bodyBytes))
	resp.Body.Close()
}

func TestAllowMethods(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	allowedMethods := []string{http.MethodGet, http.MethodPut}
	router.Use(New(Options{
		AllowMethods: allowedMethods,
	}))

	router.Get("/a", func(c *looli.Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	// preflight request
	getReq, err := http.NewRequest(http.MethodOptions, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)
	getReq.Header.Add("Access-Control-Request-Headers", "fake-header1, fake-header2")
	getReq.Header.Set("Access-Control-Request-Method", http.MethodGet)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Equal(t, strings.Join(allowedMethods, ", "), resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, resp.Header.Get("Access-Control-Allow-Headers"), "fake-header1, fake-header2")
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Empty(t, bodyBytes)
	resp.Body.Close()

	// real request
	getReq, err = http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err = http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()
}

func TestAllowCredentials(t *testing.T) {
	t.Run("normal cors request", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		origin := "looli.xyz"
		router := looli.New()
		router.Use(New(Options{
			AllowCredentials: true,
		}))

		router.Get("/a", func(c *looli.Context) {
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
		assert.Equal(t, resp.Header.Get("Access-Control-Allow-Credentials"), "true")
		assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
		assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("preflighted New requests", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		origin := "looli.xyz"
		router := looli.New()
		router.Use(New(Options{
			AllowCredentials: true,
		}))

		router.Get("/a", func(c *looli.Context) {
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL

		// preflight request
		getReq, err := http.NewRequest(http.MethodOptions, serverURL+"/a", nil)
		assert.Nil(t, err)
		getReq.Header.Set("Origin", origin)
		getReq.Header.Add("Access-Control-Request-Headers", "fake-header1, fake-header2")
		getReq.Header.Set("Access-Control-Request-Method", http.MethodGet)

		resp, err := http.DefaultClient.Do(getReq)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Origin", resp.Header.Get("Vary"))
		assert.Equal(t, strings.Join(defaultAllowMethods, ", "), resp.Header.Get("Access-Control-Allow-Methods"))
		assert.Equal(t, resp.Header.Get("Access-Control-Allow-Headers"), "fake-header1, fake-header2")
		assert.Equal(t, resp.Header.Get("Access-Control-Allow-Credentials"), "true")
		assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
		assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Empty(t, bodyBytes)
		resp.Body.Close()

		// real request
		getReq, err = http.NewRequest(http.MethodGet, serverURL+"/a", nil)
		assert.Nil(t, err)
		getReq.Header.Set("Origin", origin)

		resp, err = http.DefaultClient.Do(getReq)
		assert.Nil(t, err)

		assert.Equal(t, statusCode, resp.StatusCode)
		assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "Origin", resp.Header.Get("Vary"))
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
		assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
		assert.Equal(t, resp.Header.Get("Access-Control-Allow-Credentials"), "true")
		assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
		assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, serverResponse, string(bodyBytes))
		resp.Body.Close()
	})
}

func TestExposeHeaders(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{
		ExposeHeaders: []string{"X-My-Custom-Header", "X-Another-Custom-Header"},
	}))

	router.Get("/a", func(c *looli.Context) {
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
	assert.Equal(t, resp.Header.Get("Access-Control-Expose-Headers"), "X-My-Custom-Header, X-Another-Custom-Header")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestMaxAge(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{
		MaxAge: time.Second * 10,
	}))

	router.Get("/a", func(c *looli.Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	// preflight request
	getReq, err := http.NewRequest(http.MethodOptions, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)
	getReq.Header.Add("Access-Control-Request-Headers", "fake-header1, fake-header2")
	getReq.Header.Set("Access-Control-Request-Method", http.MethodGet)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Equal(t, strings.Join(defaultAllowMethods, ", "), resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, resp.Header.Get("Access-Control-Allow-Headers"), "fake-header1, fake-header2")
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, resp.Header.Get("Access-Control-Max-Age"), "10")
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Empty(t, bodyBytes)
	resp.Body.Close()

	// real request
	getReq, err = http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err = http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()
}

func TestAllowHeaders(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{
		AllowHeaders: []string{"fake-header1", "fake-header2", "fake-header3"},
	}))

	router.Get("/a", func(c *looli.Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	// preflight request
	getReq, err := http.NewRequest(http.MethodOptions, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)
	getReq.Header.Add("Access-Control-Request-Headers", "fake-header1, fake-header2")
	getReq.Header.Set("Access-Control-Request-Method", http.MethodGet)

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Equal(t, strings.Join(defaultAllowMethods, ", "), resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, resp.Header.Get("Access-Control-Allow-Headers"), "fake-header1, fake-header2, fake-header3")
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Empty(t, bodyBytes)
	resp.Body.Close()

	// real request
	getReq, err = http.NewRequest(http.MethodGet, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)

	resp, err = http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, origin, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
	resp.Body.Close()
}

func TestRequestMethodEmpty(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	origin := "looli.xyz"
	router := looli.New()
	router.Use(New(Options{}))

	router.Get("/a", func(c *looli.Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	// preflight request
	getReq, err := http.NewRequest(http.MethodOptions, serverURL+"/a", nil)
	assert.Nil(t, err)
	getReq.Header.Set("Origin", origin)
	getReq.Header.Add("Access-Control-Request-Headers", "fake-header1, fake-header2")
	getReq.Header.Set("Access-Control-Request-Method", "")

	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", resp.Header.Get("Vary"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Empty(t, resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Empty(t, resp.Header.Get("Access-Control-Max-Age"))
	assert.Empty(t, resp.Header.Get("Access-Control-Expose-Headers"))
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, "invalid preflighted request, missing Access-Control-Request-Method header", string(bodyBytes))
}
