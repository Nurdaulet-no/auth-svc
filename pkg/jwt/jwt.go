package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	ttl time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

func (m *Manager) Issue(userID string) (string, error){
	claims := jwtlib.RegisteredClaims{
		Subject: userID,
		ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.ttl)),
		IssuedAt: jwtlib.NewNumericDate(time.Now()),
	}

	t := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return t.SignedString(m.secret)
}


func (m *Manager) Parse(tokenStr string) (string, error){
	token, err := jwtlib.ParseWithClaims(tokenStr, &jwtlib.RegisteredClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		return m.secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwtlib.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil
}