package looli

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQuery(t *testing.T) {
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Empty(t, c.Query(""))
		assert.Equal(t, c.Query("name"), "cssivision")
		assert.Equal(t, c.Query("age"), "23")
		assert.Equal(t, c.Query("bar"), "イモト")
		assert.Empty(t, c.Query("other"))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "?&name=cssivision&age=23&bar=イモト")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestDefaultQuery(t *testing.T) {
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Empty(t, c.DefaultQuery("", ""))
		assert.Equal(t, c.DefaultQuery("name", "biz"), "cssivision")
		assert.Equal(t, c.DefaultQuery("age", "24"), "23")
		assert.Equal(t, c.DefaultQuery("other", "other value"), "other value")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL + "?&name=cssivision&age=23")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestPostForm(t *testing.T) {
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Equal(t, c.PostForm("foo"), "bar")
		assert.Equal(t, c.PostForm("page"), "11")
		assert.Empty(t, c.PostForm("both"))
		assert.Equal(t, c.PostForm("foo"), "second")
		assert.Empty(t, c.PostForm("other"))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	body := bytes.NewReader([]byte("foo=bar&page=11&both=&foo=second"))
	resp, err := http.Post(serverURL+"/path", MIMEPOSTForm, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestDefaultPostForm(t *testing.T) {
	router := New()
	router.Get("/path", func(c *Context) {
		assert.Equal(t, c.DefaultPostForm("foo", "hh"), "bar")
		assert.Equal(t, c.DefaultPostForm("page", "12"), "11")
		assert.Equal(t, c.DefaultPostForm("both", "other"), "other")
		assert.Equal(t, c.DefaultPostForm("foo", "third"), "second")
		assert.Empty(t, c.DefaultPostForm("other", ""))
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL

	body := bytes.NewReader([]byte("foo=bar&page=11&both=&foo=second"))
	resp, err := http.Post(serverURL+"/path", MIMEPOSTForm, body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestPostFormMultipart(t *testing.T) {
	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	must(mw.SetBoundary(boundary))
	must(mw.WriteField("foo", "bar"))
	must(mw.WriteField("bar", "10"))
	must(mw.WriteField("array", "first"))
	must(mw.WriteField("id", "12"))
	mw.Close()

	router := New()
	router.Post("/a/b", func(c *Context) {
		assert.Equal(t, c.PostForm("bar"), "10")
		assert.Equal(t, c.PostForm("foo"), "bar")
		assert.Equal(t, c.PostForm("array"), "first")
		assert.Equal(t, c.PostForm("id"), "12")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodPost, serverURL+"/a/b", body)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	_, err = http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func TestStatus(t *testing.T) {
	statusCode := 404
	router := New()
	router.Get("/", func(c *Context) {
		c.Status(statusCode)
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
}

func TestRedirect(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/redirect", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})
	router.Get("/", func(c *Context) {
		c.Redirect("/redirect")
	})

	server := httptest.NewServer(router)
	defer server.Close()
	serverURL := server.URL
	resp, err := http.Get(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestServeFile(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.ServeFile("test/index.html")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadFile("test/index.html")
	if err != nil {
		t.Fatal(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, content, bodyBytes)
}

func TestHeader(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		assert.Equal(t, c.Header("fake-header"), "fake")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("fake-header", "fake")
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSetHeader(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.SetHeader("fake-header", "fake")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, resp.Header.Get("fake-header"), "fake")
}

func TestCookie(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		val, err := c.Cookie("fake-cookie")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, val, "fake")
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodGet, serverURL+"/a/b", nil)
	if err != nil {
		t.Fatal(err)
	}

	getReq.AddCookie(&http.Cookie{
		Name:  "fake-cookie",
		Value: "fake",
	})
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestSetCookie(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.SetCookie(&http.Cookie{
			Name:  "fake-cookie",
			Value: "fake",
		})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, "fake", resp.Cookies()[0].Value)
}

func TestAbort(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	middleware1 := func(c *Context) {
		c.Status(statusCode)
		c.Abort()
		assert.True(t, c.IsAborted())
	}
	middleware2 := func(c *Context) {
		c.String(serverResponse)
	}

	router := New()
	router.Use(middleware1, middleware2)
	router.Get("/a/b", func(c *Context) {})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Empty(t, bodyBytes)
}

func TestAbortWithStatus(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	middleware1 := func(c *Context) {
		c.AbortWithStatus(statusCode)
		assert.True(t, c.IsAborted())
	}
	middleware2 := func(c *Context) {
		c.String(serverResponse)
	}

	router := New()
	router.Use(middleware1, middleware2)
	router.Get("/a/b", func(c *Context) {})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Empty(t, bodyBytes)
}

func TestParam(t *testing.T) {
	router := New()
	router.Get("/a/:name", func(c *Context) {
		assert.Equal(t, "cssivision", c.Param("name"))
		assert.Empty(t, c.Param("other"))
	})

	router.Get("/b/:filepath", func(c *Context) {
		assert.Equal(t, "c/cssivision", c.Param("filepath"))
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/cssivision")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	resp, err = http.Get(serverURL + "/b/c/cssivision")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestContentType(t *testing.T) {
	router := New()
	router.Get("/a/b", func(c *Context) {
		assert.Equal(t, "text/plain", c.ContentType())
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodPost, serverURL+"/a/b", nil)
	getReq.Header.Set("Content-Type", "text/plain")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(getReq)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
}

func TestString(t *testing.T) {
	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Get("/a/b", func(c *Context) {
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	resp, err := http.Get(serverURL + "/a/b")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}
