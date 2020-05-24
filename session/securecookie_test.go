package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeCookie(t *testing.T) {
	secret := []byte("secret")
	str, err := EncodeCookie(secret, nil, "sess", "cdcdf")
	assert.Nil(t, err)
	assert.NotNil(t, str)
}

func TestDecodeCookie(t *testing.T) {
	dst := make(map[interface{}]interface{})
	secret := []byte("secret")
	err := DecodeCookie(secret, nil, "sess", "cdcdf", &dst)
	assert.NotNil(t, err)
}
