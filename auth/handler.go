package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwtPkg "github.com/golang-jwt/jwt/v4"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/user/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler interface {
	GetAuthentication() gin.HandlerFunc
	GetAuthorization() gin.HandlerFunc
}

type authHandler struct {
	authTrustedProxy string
	jwt              jwt.Jwt
	userRepository   repository.UserRepository
}

func NewAuthHandler(jwt jwt.Jwt, userRepository repository.UserRepository) AuthHandler {
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
		userId, err := primitive.ObjectIDFromHex(claims["sub"].(string))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authentication": "failed", "err": err.Error()})
			return
		}
		context.Set("userId", userId)
		context.Next()
	}
}

func (a *authHandler) GetAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, ok := context.Get("userId")
		if !ok {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": "cannot read the user"})
			return
		}
		currentUser, err := a.userRepository.Read(userId.(string))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": err.Error()})
			return
		}
		if currentUser.Role != "default" {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": "you don't have enough permissions"})
			return
		}
		context.Next()
	}
}
