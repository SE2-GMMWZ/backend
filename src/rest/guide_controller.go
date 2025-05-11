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

// Create godoc
// @Summary      Create guide
// @Description  Creates a new guide
// @Tags         Guides
// @Accept       json
// @Produce      json
// @Param        body  body      model.Guide  true  "Guide to create"
// @Success      201   {object}  model.Guide
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /guides [post]
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

// Get godoc
// @Summary      Get guide
// @Description  Returns a single guide by its UUID
// @Tags         Guides
// @Produce      json
// @Param        id   path      string       true  "Guide UUID"
// @Success      200  {object}  model.Guide
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /guides/{id} [get]
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

// Update godoc
// @Summary      Update guide
// @Description  Updates a guide by its UUID
// @Tags         Guides
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Guide UUID"
// @Param        body  body      model.Guide  true  "Guide data to update"
// @Success      200   {object}  model.Guide
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /guides/{id} [put]
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

// Delete godoc
// @Summary      Delete guide
// @Description  Deletes a guide by its UUID
// @Tags         Guides
// @Produce      json
// @Param        id   path      string  true  "Guide UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /guides/{id} [delete]
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

// List godoc
// @Summary      List guides
// @Description  Returns paginated list of guides; supports query parameters forwarded to repository
// @Tags         Guides
// @Produce      json
// @Param        page   query     int false "Page number"
// @Param        limit  query     int false "Items per page"
// @Success      200    {object}  map[string]interface{}
// @Router       /guides/list [get]
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
