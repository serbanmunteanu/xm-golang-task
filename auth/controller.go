package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/user"
	"github.com/serbanmunteanu/xm-golang-task/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	userRepository *user.UserRepository
	jwt            *jwt.Jwt
}

func NewAuthController(userRepository *user.UserRepository, jwt *jwt.Jwt) *AuthController {
	return &AuthController{
		userRepository: userRepository,
		jwt:            jwt,
	}
}

func (ac *AuthController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/auth")

	router.POST("/register", ac.singUp)
	router.POST("/login", ac.signIn)
}

func (ac *AuthController) singUp(context *gin.Context) {
	var signUp *SignUpInput

	if err := context.ShouldBindJSON(&signUp); err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	hashedPassword, err := utils.HashPassword(signUp.Password)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	newUser := &user.UserDbModel{
		Name:      signUp.Name,
		Email:     strings.ToLower(signUp.Email),
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Role:      "default",
		Verified:  false,
	}
	err = ac.userRepository.Insert(newUser)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	accessToken, err := ac.jwt.CreateJwt(newUser.ID)

	context.JSON(
		http.StatusCreated,
		gin.H{"status": "success", "user": user.MapToUserResponse(newUser), "accessToken": accessToken},
	)
}

func (ac *AuthController) signIn(context *gin.Context) {
	var credentials *SignInInput

	if err := context.ShouldBindJSON(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userRepository.FindOneBy(bson.M{"email": strings.ToLower(credentials.Email)})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err = utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		return
	}

	accessToken, err := ac.jwt.CreateJwt(user.ID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}
