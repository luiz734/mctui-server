package app

import (
	"encoding/json"
	"fmt"
	"log"
	"mctui-server/env"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var EXPIRATION_TIME_SEC = 3600

func createToken(secretKey []byte, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Second * time.Duration(EXPIRATION_TIME_SEC)).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(secretKey []byte, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		log.Printf("Invalid jwt token")
		return fmt.Errorf("invalid token")
	}

	// Extract claims if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			log.Printf("Authenticated user %s using jwt", username)
			return nil
		}
		return fmt.Errorf("Missing claims")
	}

	return nil
}

type User struct {
	Username string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	json.NewDecoder(r.Body).Decode(&u)
	log.Printf("User %s attempt to login", u.Username)

	if u.Username == "admin" && u.Password == "1234" {
		secretKey := env.GetJwtSecret()
		tokenString, err := createToken(secretKey, u.Username)
		log.Printf("Creating token for user %s", u.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Errorf("No username found")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "Invalid credentials")
	log.Printf("Invalid credentials for user %s", u.Username)
}
