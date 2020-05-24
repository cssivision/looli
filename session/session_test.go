package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCookie(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		cookie := NewCookie("name", "value", nil)
		assert.Equal(t, cookie.Name, "name")
		assert.Equal(t, cookie.Value, "value")
		assert.Equal(t, cookie.Path, defaultOptions.Path)
		assert.Equal(t, cookie.MaxAge, defaultOptions.MaxAge)
		assert.Equal(t, cookie.HttpOnly, defaultOptions.HttpOnly)
	})

	t.Run("with options", func(t *testing.T) {
		options := &Options{
			Path:     "",
			Domain:   "",
			MaxAge:   10,
			Secure:   true,
			HttpOnly: false,
		}

		cookie := NewCookie("name", "value", options)
		assert.Equal(t, cookie.Name, "name")
		assert.Equal(t, cookie.Value, "value")
		assert.Equal(t, cookie.Path, options.Path)
		assert.Equal(t, cookie.MaxAge, options.MaxAge)
		assert.Equal(t, cookie.HttpOnly, options.HttpOnly)
		assert.Equal(t, cookie.Secure, true)
		assert.Equal(t, cookie.HttpOnly, false)
	})
}
