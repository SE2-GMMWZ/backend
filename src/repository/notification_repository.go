package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db}
}

func (r *NotificationRepository) CreateNotification(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

func (r *NotificationRepository) GetNotificationByID(id uuid.UUID) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.First(&notification, "notification_id = ?", id).Error
	return &notification, err
}

func (r *NotificationRepository) UpdateNotification(notification *model.Notification) error {
	return r.db.Save(notification).Error
}

func (r *NotificationRepository) DeleteNotification(id uuid.UUID) error {
	return r.db.Delete(&model.Notification{}, "notification_id = ?", id).Error
}

func (r *NotificationRepository) ListNotifications(limit, offset int, query string) ([]model.Notification, error) {
	var notifications []model.Notification
	db := r.db.Model(&model.Notification{})

	// Basic query parsing: support "user_id=...", "message=..."
	if query != "" {
		parts := strings.Split(query, "&")
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := kv[0]
			value := kv[1]
			switch key {
			case "user_id":
				db = db.Where("user_id = ?", value)
			case "message":
				db = db.Where("message ILIKE ?", "%"+value+"%")
			}
		}
	}

	err := db.Order("notification_id").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}
