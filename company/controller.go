package company

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyController interface {
	Create(context *gin.Context)
	Read(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	Register(routerGroup *gin.RouterGroup)
}

type companyController struct {
	repo CompanyRepository
}

func (c companyController) Register(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/company")

	router.POST("/", c.Create)
	router.GET("/:id", c.Read)
	router.PATCH("/:id", c.Update)
	router.DELETE("/:id", c.Delete)
}

func (c companyController) Create(context *gin.Context) {
	var model *Model
	if err := context.ShouldBindJSON(&model); err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	count, err := c.repo.Count(model.Name)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if count > 0 {
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Field name must be unique"})
		return
	}
	err = c.repo.Create(model)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	context.JSON(
		http.StatusCreated,
		model,
	)
}

func (c companyController) Read(context *gin.Context) {
	id := context.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	model, err := c.repo.Read(objectId)
	if err != nil {
		log.Info(err)
		status := http.StatusBadRequest
		if err == mongo.ErrNoDocuments {
			status = http.StatusNotFound
		}
		context.JSON(status, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	context.JSON(
		http.StatusOK,
		model,
	)
}

func (c companyController) Update(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (c companyController) Delete(context *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func NewCompanyController(repo CompanyRepository) CompanyController {
	return companyController{
		repo: repo,
	}
}
