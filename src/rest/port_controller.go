package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type portController struct {
	portRepository *repository.PortRepository
}

// Create godoc
// @Summary      Create port
// @Description  Creates a new port
// @Tags         Ports
// @Accept       json
// @Produce      json
// @Param        body  body      model.Port  true  "Port to create"
// @Success      201   {object}  model.Port
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /ports [post]
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

// Get godoc
// @Summary      Get port
// @Description  Returns a single port by its UUID
// @Tags         Ports
// @Produce      json
// @Param        id   path      string      true  "Port UUID"
// @Success      200  {object}  model.Port
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /ports/{id} [get]
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

// Update godoc
// @Summary      Update port
// @Description  Updates a port by its UUID
// @Tags         Ports
// @Accept       json
// @Produce      json
// @Param        id    path      string      true  "Port UUID"
// @Param        body  body      model.Port  true  "Port data to update"
// @Success      200   {object}  model.Port
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /ports/{id} [put]
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

// Delete godoc
// @Summary      Delete port
// @Description  Deletes a port by its UUID
// @Tags         Ports
// @Produce      json
// @Param        id   path      string  true  "Port UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /ports/{id} [delete]
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

// List godoc
// @Summary      List ports
// @Description  Returns paginated list of ports; supports query parameters forwarded to repository
// @Tags         Ports
// @Produce      json
// @Param        name         query     string false "Filter by port name (partial match)"
// @Param        owner_id     query     string false "Filter by owner UUID"
// @Param        is_approved  query     bool   false "Filter by approval status"
// @Param        limit        query     int    false "Items per page"
// @Success      200    {object}  []model.Port
// @Router       /ports/list [get]
func (pc *portController) List(c *gin.Context) {
	shared.ListWithQuery(c, pc.portRepository.ListPorts)
}

func AddPortRoutes(r *gin.Engine, portRepository *repository.PortRepository) {
	pc := portController{portRepository: portRepository}

	ports := r.Group("/ports")
	{
		ports.POST("", pc.Create)
		ports.GET("/:id", pc.Get)
		ports.GET("/list", pc.List)
		ports.PUT("/:id", pc.Update)
		ports.DELETE("/:id", pc.Delete)
	}
}
