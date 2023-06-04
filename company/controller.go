package company

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serbanmunteanu/xm-golang-task/kafka"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IController interface {
	Create(context *gin.Context)
	Read(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
	Register(routerGroup *gin.RouterGroup)
}

type companyController struct {
	repo     IRepository
	producer kafka.Producer
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
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	err = c.repo.Create(model)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	event, err := json.Marshal(model)
	if err != nil {
		log.Info(err)
	}
	err = c.producer.Produce([]byte(fmt.Sprintf("%s %s", "Company created: ", model.ID.String())), event)
	if err != nil {
		log.Info(err)
	}
	context.JSON(
		http.StatusCreated,
		model,
	)
}

func (c companyController) Read(context *gin.Context) {
	id := context.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	model, err := c.repo.Read(objectID)
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
	id := context.Param("id")
	var updates map[string]interface{}
	if err := context.ShouldBindJSON(&updates); err != nil {
		context.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	updates["updatedAt"] = time.Now()
	err = c.repo.Update(objectID, updates)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	model, err := c.repo.Read(objectID)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	event, err := json.Marshal(updates)
	if err != nil {
		log.Info(err)
	}
	err = c.producer.Produce([]byte(fmt.Sprintf("%s %s", "Company updated: ", model.ID.String())), event)
	if err != nil {
		log.Info(err)
	}
	context.JSON(
		http.StatusOK,
		model,
	)
}

func (c companyController) Delete(context *gin.Context) {
	id := context.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	err = c.repo.Delete(objectID)
	if err != nil {
		log.Info(err)
		context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	err = c.producer.Produce([]byte(fmt.Sprintf("%s %s", "Company delete: ", objectID.String())), []byte(id))
	if err != nil {
		log.Info(err)
	}
	context.JSON(http.StatusOK, gin.H{})
}

func NewCompanyController(repo IRepository, producer kafka.Producer) IController {
	return companyController{
		repo:     repo,
		producer: producer,
	}
}
