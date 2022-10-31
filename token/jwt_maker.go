package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

const minSecretKey = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secret_key string) (Maker, error) {
	if len(secret_key) < minSecretKey {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKey)
	}

	return &JWTMaker{secretKey: secret_key}, nil
}

func (maker *JWTMaker) CreateToken(user_id int64, duration time.Duration) (string, error) {
	payload, err := NewPayload(user_id, duration)
	if err != nil {
		return "", nil
	}

	jwt_token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwt_token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) { //verify jwt algorithm token header
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwt_token, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwt_token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil

}
