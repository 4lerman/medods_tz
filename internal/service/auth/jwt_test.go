package auth

import (
	"testing"
)

func TestCreateJWT(t *testing.T) {
	accessTokenSecret := []byte("accessTokenSecret")
	refreshTokenSecret := []byte("refreshTokenSecret")

	accessToken, refreshToken, err := CreateJWT(accessTokenSecret, refreshTokenSecret, 1)
	if err != nil {
		t.Errorf("error creating jwt %v", err)
	}

	if accessToken == "" || refreshToken == ""{
		t.Error("expected token to be not empty")
	}
}
