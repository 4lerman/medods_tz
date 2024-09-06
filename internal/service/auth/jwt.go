package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	configs "github.com/4lerman/medods_tz/internal/config"
	"github.com/4lerman/medods_tz/utils"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(accessTokenSecret, refreshTokenSecret []byte, userID int) (string, string, error) {
	access_token_exp := configs.Envs.AccessTokenExp
	expiration := time.Second * time.Duration(access_token_exp)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := accessToken.SignedString(accessTokenSecret)
	if err != nil {
		return "", "", err
	}

	refresh_token_exp := configs.Envs.RefreshTokenExp
	expiration = time.Second * time.Duration(refresh_token_exp)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(refreshTokenSecret)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func GetUserFromToken(r *http.Request) int {

		tokenString := getTokenFromRequest(r)

		token, err := ValidateToken(tokenString)
		if err != nil {
			log.Panicf("failed to validate token: %v", err)
		}

		if !token.Valid {
			log.Panic("invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)

		return userID
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(configs.Envs.AccessTokenSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
