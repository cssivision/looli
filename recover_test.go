package looli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverWithWriter(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(RecoverWithWriter(buffer))
	router.Get("/", func(c *Context) {
		panic("error panic")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
	assert.Contains(t, buffer.String(), "error panic")
}

func TestPanicWithAbort(t *testing.T) {
	router := New()
	router.Use(RecoverWithWriter(nil))
	router.Get("/", func(c *Context) {
		panic("error panic")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
}
