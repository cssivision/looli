package looli

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"reflect"
	"strconv"
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

	if err := mapForm(data, req.Form); err != nil {
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

func mapForm(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	// fmt.Println(typ, val, typ.NumField())

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get("json")
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Struct {
				err := mapForm(structField.Addr().Interface(), form)
				if err != nil {
					return err
				}
				continue
			}
		}

		inputValue, exists := form[inputFieldName]
		if !exists {
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
		} else {
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return nil
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
