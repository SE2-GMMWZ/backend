package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type dockingSpotController struct {
	dockingSpotRepository *repository.DockingSpotRepository
}

// Create godoc
// @Summary      Create docking-spot
// @Description  Creates a new docking spot
// @Tags         DockingSpots
// @Accept       json
// @Produce      json
// @Param        body  body      model.DockingSpot  true  "Docking-spot to create"
// @Success      201   {object}  model.DockingSpot
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /docking-spots [post]
func (dsc *dockingSpotController) Create(c *gin.Context) {
	var spot model.DockingSpot
	if err := c.ShouldBindJSON(&spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := dsc.dockingSpotRepository.CreateDockingSpot(&spot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create docking spot"})
		return
	}

	c.JSON(http.StatusCreated, spot)
}

// Get godoc
// @Summary      Get docking-spot
// @Description  Returns a single docking spot by its UUID
// @Tags         DockingSpots
// @Produce      json
// @Param        id   path      string            true  "Docking-spot UUID"
// @Success      200  {object}  model.DockingSpot
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /docking-spots/{id} [get]
func (dsc *dockingSpotController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	spot, err := dsc.dockingSpotRepository.GetDockingSpotByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "docking spot not found"})
		return
	}

	c.JSON(http.StatusOK, spot)
}

// Update godoc
// @Summary      Update docking-spot
// @Description  Updates a docking spot by its UUID
// @Tags         DockingSpots
// @Accept       json
// @Produce      json
// @Param        id    path      string            true  "Docking-spot UUID"
// @Param        body  body      model.DockingSpot true  "Docking-spot data to update"
// @Success      200   {object}  model.DockingSpot
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /docking-spots/{id} [put]
func (dsc *dockingSpotController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var spot model.DockingSpot
	if err := c.ShouldBindJSON(&spot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	spot.DockID = id
	if err := dsc.dockingSpotRepository.UpdateDockingSpot(&spot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update docking spot"})
		return
	}

	c.JSON(http.StatusOK, spot)
}

// Delete godoc
// @Summary      Delete docking-spot
// @Description  Deletes a docking spot by its UUID
// @Tags         DockingSpots
// @Produce      json
// @Param        id   path      string  true  "Docking-spot UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /docking-spots/{id} [delete]
func (dsc *dockingSpotController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := dsc.dockingSpotRepository.DeleteDockingSpot(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete docking spot"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary      List docking-spots
// @Description  Returns paginated list of docking spots; supports query parameters forwarded to repository
// @Tags         DockingSpots
// @Produce      json
// @Param        name         query     string false "Filter by docking spot name (partial match)"
// @Param        owner_id     query     string false "Filter by owner UUID"
// @Param        availability query     string false "Filter by availability ('available' or 'unavailable')"
// @Param        limit        query     int    false "Items per page"
// @Success      200    {object}  []model.DockingSpot
// @Router       /docking-spots/list [get]
func (dsc *dockingSpotController) List(c *gin.Context) {
	shared.ListWithQuery(c, dsc.dockingSpotRepository.ListDockingSpots)
}

func AddDockingSpotRoutes(r *gin.Engine, dockingSpotRepository *repository.DockingSpotRepository) {
	dsc := dockingSpotController{dockingSpotRepository: dockingSpotRepository}

	dockingSpots := r.Group("/docking-spots")
	{
		dockingSpots.POST("", dsc.Create)
		dockingSpots.GET("/:id", dsc.Get)
		dockingSpots.GET("/list", dsc.List)
		dockingSpots.PUT("/:id", dsc.Update)
		dockingSpots.DELETE("/:id", dsc.Delete)
	}
}
