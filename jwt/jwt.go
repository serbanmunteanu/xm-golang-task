package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	jwtPkg "github.com/golang-jwt/jwt/v4"
	"github.com/serbanmunteanu/xm-golang-task/config"
)

type Jwt interface {
	CreateJwt(payload interface{}) (string, error)
	Validate(tokenString string) (jwtPkg.MapClaims, error)
}

type jwt struct {
	cfg config.JwtConfig
}

func NewJwt(cfg config.JwtConfig) Jwt {
	return &jwt{
		cfg: cfg,
	}
}

func (j *jwt) CreateJwt(payload interface{}) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(j.cfg.PrivateKey)
	if err != nil {
		return "", err
	}
	key, err := jwtPkg.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	claims := make(jwtPkg.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(time.Duration(j.cfg.Ttl) * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwtPkg.NewWithClaims(jwtPkg.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwt) Validate(tokenString string) (jwtPkg.MapClaims, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(j.cfg.PublicKey)
	if err != nil {
		return nil, err
	}
	key, err := jwtPkg.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwtPkg.Parse(tokenString, func(token *jwtPkg.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtPkg.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method %s", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwtPkg.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
