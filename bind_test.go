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
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Other   string `json:"other"`
		Payload struct {
			A string `json:"a"`
			B int    `json:"b"`
		} `json:"payload"`
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
		assert.Equal(t, "aaa", form.Payload.A)
		assert.Equal(t, 222, form.Payload.B)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL
	body, err := json.Marshal(JSON{
		"name": "cssivision",
		"age":  21,
		"payload": map[string]interface{}{
			"a": "aaa",
			"b": 222,
		},
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

	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		err := c.Bind(form)
		assert.Nil(t, err)

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
	assert.Nil(t, err)
	getReq.Header.Set("Content-Type", MIMEMultipartPOSTForm+"; boundary="+boundary)
	resp, err := http.DefaultClient.Do(getReq)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, statusCode, resp.StatusCode)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, serverResponse, string(bodyBytes))
}

func TestMutliDataType(t *testing.T) {
	type Info struct {
		Integer8   int8    `json:"integer8"`
		Integer16  int16   `json:"integer16"`
		Integer32  int32   `json:"integer32"`
		Integer64  int64   `json:"integer64"`
		Uinteger8  uint8   `json:"uinteger8"`
		Uinteger16 uint16  `json:"uinteger16"`
		Uinteger32 uint32  `json:"uinteger32"`
		Uinteger64 uint64  `json:"uinteger64"`
		Boolean    bool    `json:"boolean"`
		Float32    float32 `json:"float32"`
		Float64    float64 `json:"float64"`
	}

	statusCode := 404
	serverResponse := "server response"
	router := New()
	router.Post("/", func(c *Context) {
		form := new(Info)
		assert.Nil(t, c.Bind(form))

		assert.Equal(t, int8(7), form.Integer8)
		assert.Equal(t, int16(7), form.Integer16)
		assert.Equal(t, int32(7), form.Integer32)
		assert.Equal(t, int64(7), form.Integer64)
		assert.Equal(t, uint8(7), form.Uinteger8)
		assert.Equal(t, uint16(7), form.Uinteger16)
		assert.Equal(t, uint32(7), form.Uinteger32)
		assert.Equal(t, uint64(7), form.Uinteger64)
		assert.True(t, form.Boolean)
		assert.Equal(t, float32(7.7), form.Float32)
		assert.Equal(t, float64(7.7), form.Float64)
		c.Status(statusCode)
		c.String(serverResponse)
	})

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL

	data := url.Values{}
	data.Add("integer8", "7")
	data.Add("integer16", "7")
	data.Add("integer32", "7")
	data.Add("integer64", "7")
	data.Add("uinteger8", "7")
	data.Add("uinteger16", "7")
	data.Add("uinteger32", "7")
	data.Add("uinteger64", "7")
	data.Add("boolean", "true")
	data.Add("float32", "7.7")
	data.Add("float64", "7.7")
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
