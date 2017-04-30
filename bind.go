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

	return mapForm(data, req.Form)
}

func (*xmlBinding) Bind(req *http.Request, data interface{}) error {
	return xml.NewDecoder(req.Body).Decode(data)
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

		typeFieldKind := typeField.Type.Kind()
		inputFieldName := typeField.Tag.Get("json")
		if inputFieldName == "" {
			inputFieldName = typeField.Name

			// if "form" tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if typeFieldKind == reflect.Struct {
				if err := mapForm(structField.Addr().Interface(), form); err != nil {
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
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField, isPtrType)
	case reflect.Int8:
		return setIntField(val, 8, structField, isPtrType)
	case reflect.Int16:
		return setIntField(val, 16, structField, isPtrType)
	case reflect.Int32:
		return setIntField(val, 32, structField, isPtrType)
	case reflect.Int64:
		return setIntField(val, 64, structField, isPtrType)
	case reflect.Uint:
		return setUintField(val, 0, structField, isPtrType)
	case reflect.Uint8:
		return setUintField(val, 8, structField, isPtrType)
	case reflect.Uint16:
		return setUintField(val, 16, structField, isPtrType)
	case reflect.Uint32:
		return setUintField(val, 32, structField, isPtrType)
	case reflect.Uint64:
		return setUintField(val, 64, structField, isPtrType)
	case reflect.Bool:
		return setBoolField(val, structField, isPtrType)
	case reflect.Float32:
		return setFloatField(val, 32, structField, isPtrType)
	case reflect.Float64:
		return setFloatField(val, 64, structField, isPtrType)
	case reflect.String:
		if isPtrType {
			structField.Set(reflect.ValueOf(&val))
		} else {
			structField.SetString(val)
		}
	default:
		return errors.New("Unknown type")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value, isPtrType bool) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		if isPtrType {
			switch bitSize {
			case 8:
				pintVal := int8(intVal)
				field.Set(reflect.ValueOf(&pintVal))
			case 16:
				pintVal := int16(intVal)
				field.Set(reflect.ValueOf(&pintVal))
			case 32:
				pintVal := int32(intVal)
				field.Set(reflect.ValueOf(&pintVal))
			case 64:
				pintVal := int64(intVal)
				field.Set(reflect.ValueOf(&pintVal))
			default:
				pintVal := int(intVal)
				field.Set(reflect.ValueOf(&pintVal))
			}
		} else {
			field.SetInt(intVal)
		}
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value, isPtrType bool) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		if isPtrType {
			switch bitSize {
			case 8:
				puintVal := uint8(uintVal)
				field.Set(reflect.ValueOf(&puintVal))
			case 16:
				puintVal := uint16(uintVal)
				field.Set(reflect.ValueOf(&puintVal))
			case 32:
				puintVal := uint32(uintVal)
				field.Set(reflect.ValueOf(&puintVal))
			case 64:
				puintVal := uint64(uintVal)
				field.Set(reflect.ValueOf(&puintVal))
			default:
				puintVal := uint(uintVal)
				field.Set(reflect.ValueOf(&puintVal))
			}
		} else {
			field.SetUint(uintVal)
		}
	}
	return err
}

func setBoolField(val string, field reflect.Value, isPtrType bool) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		if isPtrType {
			field.Set(reflect.ValueOf(&boolVal))
		} else {
			field.SetBool(boolVal)
		}
	}
	return nil
}

func setFloatField(val string, bitSize int, field reflect.Value, isPtrType bool) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		if isPtrType {
			switch bitSize {
			case 32:
				pfloatVal := float32(floatVal)
				field.Set(reflect.ValueOf(&pfloatVal))
			case 64:
				pfloatVal := float64(floatVal)
				field.Set(reflect.ValueOf(&pfloatVal))
			}
		} else {
			field.SetFloat(floatVal)
		}
	}
	return err
}
