package session

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"hash"
	"time"
)

var hashFunc = sha256.New

// encode encodes a cookie value.
//
// It serializes, optionally encrypts, signs with a message authentication code,
// and finally encodes the value.
//
// The name argument is the cookie name. It is stored with the encoded value.
// The value argument is the value to be encoded. It can be any value that can
// be encoded using the currently selected serializer
//
// It is the client's responsibility to ensure that value, when encoded using
// the current serialization/encryption settings on s and then base64-encoded,
// is shorter than the maximum permissible length.
func encodeCookie(hashKey []byte, name string, value interface{}) (string, error) {
	var b []byte
	var err error

	// serialize using gob
	if b, err = serialize(value); err != nil {
		return "", err
	}
	b = encode(b)

	// create mac for "name|date|value", Extra pipe to be used later
	b = []byte(fmt.Sprintf("%s|%d|%s|", name, time.Now().UTC().Unix(), b))
	mac := createMac(hmac.New(hashFunc, hashKey), b[:len(b)-1])

	// append mac, remove name
	b = append(b, mac...)[len(name)+1:]

	// encode to base64
	b = encode(b)

	return string(b), nil
}

// decode decodes a cookie value
//
// It decodes, verifies a message authentication code, and finally deserializes the value.
//
// The name argument is the cookie name. It must be the same name used when it was stored.
//
// The value argument is the encoded cookie value. The dst argument is where the cookie will
// be decoded. It must be a pointer.
func decodeCookie(hashKey []byte, name string, value string, dst interface{}) error {
	// decode from base64
	b, err := decode([]byte(value))
	if err != nil {
		return err
	}

	// verify mac, value is "date|value|mac".
	parts := bytes.SplitN(b, []byte("|"), 3)
	if len(parts) != 3 {
		return fmt.Errorf("value is invalid")
	}
	h := hmac.New(hashFunc, hashKey)
	b = append([]byte(name+"|"), b[:len(b)-len(parts[2])-1]...)
	if err = verifyMac(h, b, parts[2]); err != nil {
		return err
	}

	b, err = decode(parts[1])
	if err != nil {
		return err
	}

	// deserialize.
	if err = deserialize(b, dst); err != nil {
		return err
	}
	return nil
}

// encodes a value using gob
func serialize(value interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decodes a value using gob
func deserialize(src []byte, dst interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	if err := dec.Decode(dst); err != nil {
		return err
	}

	return nil
}

// encode encodes a value using base64
func encode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// decode decodes a cookie using base64
func decode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

// createMac creates a message authentication code (MAC).
func createMac(h hash.Hash, value []byte) []byte {
	h.Write(value)
	return h.Sum(nil)
}

// verifyMac verifies that a message authentication code (MAC) is valid.
func verifyMac(h hash.Hash, value []byte, mac []byte) error {
	mac2 := createMac(h, value)
	// check that both macs are of equal length, as subtle.ConstantTimeCompare
	if len(mac) == len(mac2) && subtle.ConstantTimeCompare(mac, mac2) == 1 {
		return nil
	}
	return fmt.Errorf("value is invalid")
}
