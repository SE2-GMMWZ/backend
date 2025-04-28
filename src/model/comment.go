package model

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	CommentID uuid.UUID `gorm:"column:comment_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"comment_id"`
	GuideID   uuid.UUID `gorm:"column:guide_id;type:uuid" json:"guide_id"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid" json:"user_id"`
	Content   string    `gorm:"column:content" json:"content"`
	Timestamp time.Time `gorm:"column:timestamp" json:"timestamp"`
}
