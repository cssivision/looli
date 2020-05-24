package looli

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	data := JSON{
		"name": "cssivision",
	}

	err := &Error{
		Err:  errors.New("error happen"),
		Code: http.StatusOK,
		Meta: data,
	}

	assert.Equal(t, "error happen", err.Error())
	assert.Equal(t, http.StatusOK, err.Code)
	want, e := json.Marshal(data)
	assert.Nil(t, e)
	get, e := json.Marshal(err.Meta)
	assert.Nil(t, e)
	assert.Equal(t, want, get)
}
