package looli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethod(t *testing.T) {
	handleMethod(http.MethodGet, t)
	handleMethod(http.MethodOptions, t)
	handleMethod(http.MethodPatch, t)
	handleMethod(http.MethodDelete, t)
	handleMethod(http.MethodTrace, t)
	handlePostPutMethod(http.MethodPost, t)
	handlePostPutMethod(http.MethodPut, t)
}

func TestHeadMethod(t *testing.T) {
	serverResponse := "server response"
	statusCode := 404
	router := New()
	router.Head("/a/b", func(c *Context) {
		c.Status(statusCode)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodHead, serverURL+"/a/b", bytes.NewReader(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NotEqual(t, string(bodyBytes), serverResponse)
	assert.Equal(t, string(bodyBytes), "")
}

func handlePostPutMethod(method string, t *testing.T) {
	requestBody := bytes.Repeat([]byte("a"), 1<<20)
	statusCode := 404
	serverResponse := "serverResponse"

	router := New()
	router.Handle(method, "/a/b", func(c *Context) {
		requestData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			t.Error(err)
		}

		if !bytes.Equal(requestData, requestBody) {
			t.Errorf("Server read %d request body bytes; want %d", len(requestData), len(requestBody))
		}
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	serverURL := server.URL
	defer server.Close()

	getReq, err := http.NewRequest(method, serverURL+"/a/b", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
}

func handleMethod(method string, t *testing.T) {
	serverResponse := "server response"
	statusCode := 404
	router := New()
	router.Handle(method, "/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	getReq, err := http.NewRequest(method, serverURL+"/a/b", bytes.NewReader(nil))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(bodyBytes), serverResponse)
}