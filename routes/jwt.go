package routes

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TOKEN_EXPIRE_HRS = 3

func IssueJWTForUser(user User) (string, error) {

	sec, err := lookupSecret(user.UserName)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user"] = user.UserName
	token.Claims["role"] = user.Role
	token.Claims["iss"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(TOKEN_EXPIRE_HRS * time.Hour).Unix()

	tokenString, err := token.SignedString(sec)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func IssuePwdJWTForUser(user User) (string, error) {
	sec, err := lookupSecret(user.UserName)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["user"] = user.UserName
	token.Claims["role"] = "pwd"
	token.Claims["exp"] = time.Now().Add(TOKEN_EXPIRE_HRS * time.Hour).Unix()

	tokenString, err := token.SignedString(sec)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RenewJWTfromJWT(inToken string) (string, error) {

	token, err := ValidateJWT(inToken)
	if err != nil {
		return "", err
	}

	sec, err := lookupSecret(token.Claims["user"].(string))
	if err != nil {
		return "", err
	}
	token.Claims["exp"] = time.Now().Add(TOKEN_EXPIRE_HRS * time.Hour).Unix()

	tokenString, err := token.SignedString(sec)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateJWT(inToken string) (*jwt.Token, error) {

	token, err := jwt.Parse(inToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return lookupSecret(token.Claims["user"].(string))
	})

	if token.Valid && err == nil {
		return token, nil
	}

	return nil, err
}

func lookupSecret(user string) ([]byte, error) {

	secret := os.Getenv("JWT_SECRET")
	if len(secret) == 0 {
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Generated new secret for %s: %s", user, b)
		sErr := os.Setenv("JWT_SECRET", string(b))
		if sErr != nil {
			return nil, sErr
		}
		return b, nil
	} else {
		return []byte(secret), nil
	}

}
