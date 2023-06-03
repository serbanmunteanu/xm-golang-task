package http_framework

import (
	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/auth"
	"github.com/serbanmunteanu/xm-golang-task/jwt"
	"github.com/serbanmunteanu/xm-golang-task/swagger"
	"github.com/serbanmunteanu/xm-golang-task/user"
)

type RouteRegister interface {
	Register(routerGroup *gin.RouterGroup)
}

type RouteGroup struct {
	groupPrefix    string
	routeRegisters []RouteRegister
	middlewares    []gin.HandlerFunc
}

func Initialize(router *gin.Engine, httpServer *HttpServer) {
	userRepository := user.NewUserRepository(httpServer.mongo, httpServer.config.MongoConfig.Collections.UserCollection)
	jwt := jwt.NewJwt(httpServer.config.JwtConfig)
	authHandler := auth.NewAuthHandler(jwt, userRepository)
	routeGroups := []RouteGroup{
		{
			groupPrefix: "",
			routeRegisters: []RouteRegister{
				swagger.NewSwaggerController(),
				auth.NewAuthController(userRepository, jwt),
			},
			middlewares: []gin.HandlerFunc{},
		},
		{
			groupPrefix:    "/api",
			routeRegisters: []RouteRegister{},
			middlewares: []gin.HandlerFunc{
				authHandler.GetAuthentication(),
				authHandler.GetAuthorization(),
			},
		},
	}

	for _, routeGroup := range routeGroups {
		group := router.Group(routeGroup.groupPrefix)
		for _, middle := range routeGroup.middlewares {
			group.Use(middle)
		}
		for _, route := range routeGroup.routeRegisters {
			route.Register(group)
		}
	}
}
