package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwtPkg "github.com/golang-jwt/jwt/v4"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/user"
)

type IHandler interface {
	GetAuthentication() gin.HandlerFunc
	GetAuthorization() gin.HandlerFunc
}

type authHandler struct {
	authTrustedProxy string
	jwt              jwt.Jwt
	userRepository   user.IRepository
}

func NewAuthHandler(jwt jwt.Jwt, userRepository user.IRepository) IHandler {
	return &authHandler{
		authTrustedProxy: os.Getenv("AUTH_TRUSTED_PROXY"),
		jwt:              jwt,
		userRepository:   userRepository,
	}
}

func (a *authHandler) GetAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, err := a.jwt.Validate(context.Request.Header.Get("Authorization"))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authentication": "failed", "err": err.Error()})
			return
		}
		err, ok := claims.Valid().(jwtPkg.ValidationError)
		if ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authentication": "failed", "err": err.Error()})
			return
		}
		context.Set("user", claims["sub"])
		context.Next()
	}
}

func (a *authHandler) GetAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		loggedUser, ok := context.Get("user")
		if !ok {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": "cannot read the user"})
			return
		}
		bytes, err := json.Marshal(loggedUser)
		if err != nil {
			panic(err)
		}
		var model user.Model
		err = json.Unmarshal(bytes, &model)
		if err != nil {
			panic(err)
		}
		if model.Role != "default" {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": "you don't have enough permissions"})
			return
		}
		context.Next()
	}
}
