package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type portController struct {
	portRepository *repository.PortRepository
}

func (pc *portController) Create(c *gin.Context) {
	var port model.Port
	if err := c.ShouldBindJSON(&port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := pc.portRepository.CreatePort(&port); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create port"})
		return
	}

	c.JSON(http.StatusCreated, port)
}

func (pc *portController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	port, err := pc.portRepository.GetPortByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "port not found"})
		return
	}

	c.JSON(http.StatusOK, port)
}

func (pc *portController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var port model.Port
	if err := c.ShouldBindJSON(&port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	port.PortID = id
	if err := pc.portRepository.UpdatePort(&port); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update port"})
		return
	}

	c.JSON(http.StatusOK, port)
}

func (pc *portController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := pc.portRepository.DeletePort(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete port"})
		return
	}

	c.Status(http.StatusNoContent)
}

func AddPortRoutes(r *gin.Engine, portRepository *repository.PortRepository) {
	pc := portController{portRepository: portRepository}

	ports := r.Group("/ports")
	{
		ports.POST("", pc.Create)
		ports.GET("/:id", pc.Get)
		ports.PUT("/:id", pc.Update)
		ports.DELETE("/:id", pc.Delete)
	}
}
