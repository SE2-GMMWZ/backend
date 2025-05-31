package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (r *NotificationRepository) ListNotifications(limit, offset int) ([]model.Notification, error) {
	var notifications []model.Notification
	err := r.db.Order("notification_id").Limit(limit).Offset(offset).Find(&notifications).Error
	return notifications, err
}
