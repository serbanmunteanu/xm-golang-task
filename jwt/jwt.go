package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/serbanmunteanu/xm-golang-task/config"
)

type Jwt struct {
	cfg config.JwtConfig
}

func NewJwt(cfg config.JwtConfig) *Jwt {
	return &Jwt{
		cfg: cfg,
	}
}

func (j *Jwt) CreateJwt(payload interface{}) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(j.cfg.PrivateKey)
	if err != nil {
		return "", err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(time.Duration(j.cfg.Ttl) * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *Jwt) Validate(tokenString string) (jwt.MapClaims, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(j.cfg.PublicKey)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method %s", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
