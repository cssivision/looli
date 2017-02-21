package looli

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

const (
	MIMEJSON              = "application/json"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
)

type (
	Binding interface {
		Bind(*http.Request, interface{}) error
	}

	jsonBinding          struct{}
	formBinding          struct{}
	formPostBinding      struct{}
	formMultipartBinding struct{}
	xmlBinding           struct{}
)

func bindDefault(method, contentType string) Binding {
	var binding Binding
	if method == http.MethodGet {
		binding = &formBinding{}
	} else {
		switch contentType {
		case MIMEJSON:
			binding = &jsonBinding{}
		case MIMEXML, MIMEXML2:
			binding = &xmlBinding{}
		case MIMEPOSTForm:
			binding = &formPostBinding{}
		case MIMEMultipartPOSTForm:
			binding = &formMultipartBinding{}
		default:
			binding = &formBinding{}
		}
	}

	return binding
}

func (*jsonBinding) Bind(req *http.Request, data interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}
	return nil
}

func (*formBinding) Bind(req *http.Request, data interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(1 << 32)

	if err := mapForm(data, req.PostForm); err != nil {
		return err
	}
	return nil
}

func (*formPostBinding) Bind(req *http.Request, data interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := mapForm(data, req.PostForm); err != nil {
		return err
	}
	return nil
}

func (*formMultipartBinding) Bind(req *http.Request, data interface{}) error {
	if err := req.ParseMultipartForm(32 << 10); err != nil {
		return err
	}
	if err := mapForm(data, req.MultipartForm.Value); err != nil {
		return err
	}
	return nil
}

func (*xmlBinding) Bind(req *http.Request, data interface{}) error {
	if err := xml.NewDecoder(req.Body).Decode(data); err != nil {
		return err
	}
	return nil
}