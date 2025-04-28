package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Guide struct {
	GuideID         uuid.UUID        `gorm:"column:guide_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"guide_id"`
	Title           string           `gorm:"column:title" json:"title"`
	Content         string           `gorm:"column:content" json:"content"`
	AuthorID        uuid.UUID        `gorm:"column:author_id;type:uuid" json:"author_id"`
	PublicationDate time.Time        `gorm:"column:publication_date" json:"publication_date"`
	Images          *json.RawMessage `gorm:"column:images;type:jsonb" json:"images,omitempty"`
	Links           *json.RawMessage `gorm:"column:links;type:jsonb" json:"links,omitempty"`
	Location        json.RawMessage  `gorm:"column:location;type:jsonb" json:"location"` // latitude, longitude
	IsApproved      bool             `gorm:"column:is_approved" json:"is_approved"`
}
