package routes

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/redis.v3"
)

const TOKEN_EXPIRE_HRS = 24

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
	
	if len(inToken) == 0 {
		return nil, fmt.Errorf("Token length is zero")
	}

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

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	secKey := strings.Join([]string{"user:", user, ":secret"}, "")

	ttl, _ := client.TTL(secKey).Result()
	secret, err := client.Get(secKey).Bytes()
	if err == redis.Nil || ttl < 3*time.Hour {
		fmt.Printf("New secret for %s", user)
		b, err := makeRandomKey()
		if err != nil {
			return nil, err
		}
		if err := client.Set(secKey, b, 8*time.Hour).Err(); err != nil {
			return nil, err
		}
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func makeRandomKey() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
