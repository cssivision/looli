package looli

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

var (
	plainContentType = []string{"text/plain; charset=utf-8"}
	jsonContentType  = []string{"application/json; charset=utf-8"}
	htmlContentType  = []string{"text/html; charset=utf-8"}
)

func setContentType(rw http.ResponseWriter, value []string) {
	header := rw.Header()
	header["Content-Type"] = value
}

func renderString(rw http.ResponseWriter, format string, values ...interface{}) (err error) {
	setContentType(rw, plainContentType)

	if len(values) > 0 {
		_, err = fmt.Fprintf(rw, format, values...)
	} else {
		_, err = io.WriteString(rw, format)
	}
	return
}

func renderJSON(rw http.ResponseWriter, data interface{}) error {
	setContentType(rw, jsonContentType)
	return json.NewEncoder(rw).Encode(data)
}

func renderHTML(rw http.ResponseWriter, templ *template.Template, name string, data interface{}) error {
	setContentType(rw, htmlContentType)
	if name == "" {
		return templ.Execute(rw, data)
	}

	return templ.ExecuteTemplate(rw, name, data)
}
