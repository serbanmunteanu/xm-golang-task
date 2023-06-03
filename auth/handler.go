package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	jwt2 "github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authTrustedProxy string
	jwt              *jwt2.Jwt
	userRepository   *user.UserRepository
}

func NewAuthHandler(jwt *jwt2.Jwt, userRepository *user.UserRepository) *AuthHandler {
	return &AuthHandler{
		authTrustedProxy: os.Getenv("AUTH_TRUSTED_PROXY"),
		jwt:              jwt,
		userRepository:   userRepository,
	}
}

func (a *AuthHandler) GetAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims, err := a.jwt.Validate(context.Request.Header.Get("Authorization"))
		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Authentication": "failed", "err": err.Error()})
			return
		}
		err, ok := claims.Valid().(jwt.ValidationError)
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

func (a *AuthHandler) GetAuthorization() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId, ok := context.Get("userId")
		if !ok {
			context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"Authorization": "failed", "err": "cannot read the user"})
			return
		}
		currentUser, err := a.userRepository.FindOneBy(bson.M{"_id": userId})
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
