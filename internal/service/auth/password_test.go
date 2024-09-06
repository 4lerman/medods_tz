package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("random")

	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to be not empty")
	}

	if hash == "random" {
		t.Error("expected password to be hashed")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("random")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !ComparePasswords(hash, []byte("random")) {
		t.Error("expected password to match hash")
	}

	if ComparePasswords(hash, []byte("notrandom")){
		t.Error("expected password to match hash")
	}
}
