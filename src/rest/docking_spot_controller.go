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

func (dsc *dockingSpotController) List(c *gin.Context) {
	shared.ListEntities(c, dsc.dockingSpotRepository.ListDockingSpots)
}

func AddDockingSpotRoutes(r *gin.Engine, dockingSpotRepository *repository.DockingSpotRepository) {
	dsc := dockingSpotController{dockingSpotRepository: dockingSpotRepository}

	dockingSpots := r.Group("/docking-spots")
	{
		dockingSpots.POST("", dsc.Create)
		dockingSpots.GET("/:id", dsc.Get)
		dockingSpots.GET("/list", dsc.List) // This line adds the List route
		dockingSpots.PUT("/:id", dsc.Update)
		dockingSpots.DELETE("/:id", dsc.Delete)
	}
}
