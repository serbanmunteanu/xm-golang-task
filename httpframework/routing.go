package httpframework

import (
	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/auth"
	"github.com/serbanmunteanu/xm-golang-task/company"
	"github.com/serbanmunteanu/xm-golang-task/di"
	"github.com/serbanmunteanu/xm-golang-task/swagger"
)

type RouteRegister interface {
	Register(routerGroup *gin.RouterGroup)
}

type RouteGroup struct {
	groupPrefix    string
	routeRegisters []RouteRegister
	middlewares    []gin.HandlerFunc
}

func Initialize(router *gin.Engine, di *di.DI) {
	authHandler := auth.NewAuthHandler(di.Jwt, di.UserRepository)

	routeGroups := []RouteGroup{
		{
			groupPrefix: "",
			routeRegisters: []RouteRegister{
				swagger.NewSwaggerController(),
				auth.NewAuthController(di.UserRepository, di.Jwt),
			},
			middlewares: []gin.HandlerFunc{},
		},
		{
			groupPrefix: "/api",
			routeRegisters: []RouteRegister{
				company.NewCompanyController(di.CompanyRepository, di.Producer),
			},
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
