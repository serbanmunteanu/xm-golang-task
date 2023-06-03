package company

import "github.com/gin-gonic/gin"

type CompanyController interface {
	Create(context *gin.Context)
	Read(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	Register(routerGroup *gin.RouterGroup)
}

type companyController struct {
}

func (c companyController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/company")

	router.POST("/", c.Create)
	router.GET("/:id", c.Read)
	router.PATCH("/:id", c.Update)
	router.DELETE("/:id", c.Delete)
}

func (c companyController) Create(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c companyController) Read(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c companyController) Update(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c companyController) Delete(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewCompanyController() CompanyController {
	return companyController{}
}
