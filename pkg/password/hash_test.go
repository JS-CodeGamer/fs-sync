package password

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func generate_random_password(bytes int8) string {
	pass_bytes := make([]byte, bytes)
	rand.Read(pass_bytes)
	return base64.StdEncoding.EncodeToString(pass_bytes)
}

func TestHashPassword(t *testing.T) {
	pass := generate_random_password(16)
	hash, err := HashPassword(pass)
	if hash == pass || err != nil {
		t.Fatalf(`hash not matched or error generating hash (%v)`, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pass := generate_random_password(16)
	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf(`error generating hash: %v`, err)
	}
	if !CheckPasswordHash(pass, hash) {
		t.Fatalf(`password and hash didn't match`)
	}
}
