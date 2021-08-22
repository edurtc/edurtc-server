package server

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"log"
	"time"
)

func (n *NatsSignal) ParseToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(tk *jwt.Token) (interface{}, error) {
		b := n.config.Jwt.Secret
		return b, nil
	})
	return token, err
}

func (n *NatsSignal) ClaimToken(id uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["exp"] = time.Now().Add(time.Minute * 1) // this === 1 minute

	s, err := token.SignedString(n.config.Jwt.Secret)
	if err != nil {
		log.Fatal("error: ", err)
		return "", err
	}

	return s, nil
}

