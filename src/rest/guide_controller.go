package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type guideController struct {
	guideRepository *repository.GuideRepository
}

func (gc *guideController) Create(c *gin.Context) {
	var guide model.Guide
	if err := c.ShouldBindJSON(&guide); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := gc.guideRepository.CreateGuide(&guide); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create guide"})
		return
	}

	c.JSON(http.StatusCreated, guide)
}

func (gc *guideController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	guide, err := gc.guideRepository.GetGuideByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "guide not found"})
		return
	}

	c.JSON(http.StatusOK, guide)
}

func (gc *guideController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var guide model.Guide
	if err := c.ShouldBindJSON(&guide); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	guide.GuideID = id
	if err := gc.guideRepository.UpdateGuide(&guide); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update guide"})
		return
	}

	c.JSON(http.StatusOK, guide)
}

func (gc *guideController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := gc.guideRepository.DeleteGuide(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete guide"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (gc *guideController) List(c *gin.Context) {
	shared.ListEntities(c, gc.guideRepository.ListGuides)
}

func AddGuideRoutes(r *gin.Engine, guideRepository *repository.GuideRepository) {
	gc := guideController{guideRepository: guideRepository}

	guides := r.Group("/guides")
	{
		guides.POST("", gc.Create)
		guides.GET("/:id", gc.Get)
		guides.GET("/list", gc.List)
		guides.PUT("/:id", gc.Update)
		guides.DELETE("/:id", gc.Delete)
	}
}
