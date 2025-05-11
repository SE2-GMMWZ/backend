package rest

import (
	"backend/src/model"
	"backend/src/repository"
	"backend/src/rest/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type notificationController struct {
	notificationRepository *repository.NotificationRepository
}

// Create godoc
// @Summary      Create notification
// @Description  Creates a new notification
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        body  body      model.Notification  true  "Notification to create"
// @Success      201   {object}  model.Notification
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notifications [post]
func (nc *notificationController) Create(c *gin.Context) {
	var notification model.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := nc.notificationRepository.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

// Get godoc
// @Summary      Get notification
// @Description  Returns a single notification by its UUID
// @Tags         Notifications
// @Produce      json
// @Param        id   path      string               true  "Notification UUID"
// @Success      200  {object}  model.Notification
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /notifications/{id} [get]
func (nc *notificationController) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	notification, err := nc.notificationRepository.GetNotificationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// Update godoc
// @Summary      Update notification
// @Description  Updates a notification by its UUID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Param        id    path      string               true  "Notification UUID"
// @Param        body  body      model.Notification   true  "Notification data to update"
// @Success      200   {object}  model.Notification
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /notifications/{id} [put]
func (nc *notificationController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var notification model.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	notification.NotificationID = id
	if err := nc.notificationRepository.UpdateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// Delete godoc
// @Summary      Delete notification
// @Description  Deletes a notification by its UUID
// @Tags         Notifications
// @Produce      json
// @Param        id   path      string  true  "Notification UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /notifications/{id} [delete]
func (nc *notificationController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := nc.notificationRepository.DeleteNotification(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary      List notifications
// @Description  Returns paginated list of notifications; supports query parameters forwarded to repository
// @Tags         Notifications
// @Produce      json
// @Param        page   query     int false "Page number"
// @Param        limit  query     int false "Items per page"
// @Success      200    {object}  map[string]interface{}
// @Router       /notifications/list [get]
func (nc *notificationController) List(c *gin.Context) {
	shared.ListEntities(c, nc.notificationRepository.ListNotifications)
}

func AddNotificationRoutes(r *gin.Engine, notificationRepository *repository.NotificationRepository) {
	nc := notificationController{notificationRepository: notificationRepository}

	notifications := r.Group("/notifications")
	{
		notifications.POST("", nc.Create)
		notifications.GET("/:id", nc.Get)
		notifications.GET("/list", nc.List)
		notifications.PUT("/:id", nc.Update)
		notifications.DELETE("/:id", nc.Delete)
	}
}
