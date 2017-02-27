package looli

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestBindJSON(t *testing.T) {
	type Info struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Other string `json:"other"`
	}

	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		err := c.Bind(form)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "cssivision", form.Name)
		assert.Equal(t, 21, form.Age)
		assert.Empty(t, form.Other)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	body, err := json.Marshal(Info{
		Name: "cssivision",
		Age:  21,
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(serverURL, MIMEJSON, bytes.NewReader(body))
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

func TestBindXML(t *testing.T) {
	type Info struct {
		XMLName xml.Name `xml:"person"`
		Name    string   `xml:"name"`
		Age     int      `xml:"age"`
		Other   string   `xml:"other"`
	}

	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		err := c.Bind(form)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "cssivision", form.Name)
		assert.Equal(t, 21, form.Age)
		assert.Empty(t, form.Other)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	body, err := xml.Marshal(Info{
		Name: "cssivision",
		Age:  21,
	})

	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(serverURL, MIMEXML, bytes.NewReader(body))
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

func TestBindForm(t *testing.T) {
	type Info struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Other string `json:"other"`
	}

	t.Run("Get query", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		router := New()
		router.Get("/", func(c *Context) {
			form := new(Info)
			err := c.Bind(form)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "cssivision", form.Name)
			assert.Equal(t, 21, form.Age)
			assert.Empty(t, form.Other)
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		resp, err := http.Get(serverURL + "?name=cssivision&age=21")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
	})

	t.Run("Post form", func(t *testing.T) {
		statusCode := 404
		serverResponse := "server response"
		router := New()
		router.Post("/", func(c *Context) {
			form := new(Info)
			err := c.Bind(form)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, "cssivision", form.Name)
			assert.Equal(t, 21, form.Age)
			assert.Empty(t, form.Other)
			c.Status(statusCode)
			c.String(serverResponse)
		})

		server := httptest.NewServer(router)
		defer server.Close()

		serverURL := server.URL
		data := url.Values{}
		data.Add("name", "cssivision")
		data.Add("age", "21")
		resp, err := http.Post(serverURL, MIMEPOSTForm, bytes.NewBufferString(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, serverResponse, string(bodyBytes))
	})
}

func TestBindPostForm(t *testing.T) {
	type Info struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Other string `json:"other"`
	}

	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		err := c.Bind(form)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "cssivision", form.Name)
		assert.Equal(t, 21, form.Age)
		assert.Empty(t, form.Other)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	data := url.Values{}
	data.Add("name", "cssivision")
	data.Add("age", "21")
	resp, err := http.Post(serverURL, MIMEPOSTForm, bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestBindMultiPart(t *testing.T) {
	type Info struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Other string `json:"other"`
	}

	boundary := "--testboundary"
	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)
	must(mw.SetBoundary(boundary))
	must(mw.WriteField("name", "cssivision"))
	must(mw.WriteField("age", "21"))
	mw.Close()

	statusCode := 200
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		err := c.Bind(form)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "cssivision", form.Name)
		assert.Equal(t, 21, form.Age)
		assert.Empty(t, form.Other)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	getReq, err := http.NewRequest(http.MethodPost, serverURL, body)
	if err != nil {
		t.Fatal(err)
	}
	getReq.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	resp, err := http.DefaultClient.Do(getReq)
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
