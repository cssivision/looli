package looli

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"fmt"
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

type BindingStruct interface {
	Validate() error
}

type (
	Binding interface {
		Bind(*http.Request, interface{}) error
	}

	jsonBinding struct{}
	formBinding struct{}
	xmlBinding  struct{}
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
		default: // MIMEPOSTForm, MIMEMultipartPOSTForm
			binding = &formBinding{}
		}
	}

	return binding
}

func (*jsonBinding) Bind(req *http.Request, data interface{}) error {
	return json.NewDecoder(req.Body).Decode(data)
}

func (*formBinding) Bind(req *http.Request, data interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	req.ParseMultipartForm(1 << 32)

	return parseValues(data, req.Form)
}

func (*xmlBinding) Bind(req *http.Request, data interface{}) error {
	return xml.NewDecoder(req.Body).Decode(data)
}

func parseValues(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()

	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)

		if !structField.CanSet() {
			continue
		}

		typeFieldKind := typeField.Type.Kind()
		inputFieldName := typeField.Tag.Get("json")
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			if typeFieldKind == reflect.Struct {
				if err := parseValues(structField.Addr().Interface(), form); err != nil {
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
		if typeFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := typeField.Type.Elem().Kind()
			slice := reflect.MakeSlice(typeField.Type, numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i), false); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if typeFieldKind == reflect.Ptr {
			if err := setWithProperType(typeField.Type.Elem().Kind(), inputValue[0], structField, true); err != nil {
				return err
			}
		} else {
			if err := setWithProperType(typeFieldKind, inputValue[0], structField, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value, isPtrType bool) error {
	if isPtrType {
		structField.Set(reflect.New(structField.Type().Elem()))
		structField = structField.Elem()
	}

	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setIntField(val, structField.Type().Bits(), structField)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUintField(val, structField.Type().Bits(), structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32, reflect.Float64:
		return setFloatField(val, structField.Type().Bits(), structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return tryUnmarshalValue(structField, val)
	}
	return nil
}

func tryUnmarshalValue(v reflect.Value, str string) error {
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		v = v.Addr()
	}

	if v.Type().NumMethod() > 0 {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		i := v.Interface()
		if u, ok := i.(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(str))
		}
		if u, ok := i.(json.Unmarshaler); ok {
			return u.UnmarshalJSON([]byte(str))
		}
	}
	return fmt.Errorf("unknown field type: %v", v.Type())
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
