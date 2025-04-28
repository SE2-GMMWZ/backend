package model

import (
	"github.com/google/uuid"
	"time"
)

type Notification struct {
	NotificationID uuid.UUID `gorm:"column:notification_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"notification_id"`
	UserID         uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	Message        string    `gorm:"column:message" json:"message"`
	Timestamp      time.Time `gorm:"column:timestamp" json:"timestamp"`
}
