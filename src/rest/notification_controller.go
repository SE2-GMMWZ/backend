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

func (nc *notificationController) List(c *gin.Context) {
	shared.ListEntities(c, nc.notificationRepository.ListNotifications)
}

func AddNotificationRoutes(r *gin.Engine, notificationRepository *repository.NotificationRepository) {
	nc := notificationController{notificationRepository: notificationRepository}

	notifications := r.Group("/notifications")
	{
		notifications.POST("", nc.Create)
		notifications.GET("/:id", nc.Get)
		notifications.GET("/list", nc.List) // This line adds the List route
		notifications.PUT("/:id", nc.Update)
		notifications.DELETE("/:id", nc.Delete)
	}
}
