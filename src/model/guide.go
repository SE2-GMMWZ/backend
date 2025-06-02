package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Guide struct {
	GuideID         uuid.UUID        `gorm:"column:guide_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"guide_id"`
	Title           string           `gorm:"column:title" json:"title"`
	Content         string           `gorm:"column:content" json:"content"`
	AuthorID        uuid.UUID        `gorm:"column:author_id;type:uuid" json:"author_id"`
	PublicationDate time.Time        `gorm:"column:publication_date" json:"publication_date"`
	Images          *json.RawMessage `gorm:"column:images;type:jsonb" json:"images,omitempty" swaggertype:"object"`
	Links           *json.RawMessage `gorm:"column:links;type:jsonb" json:"links,omitempty" swaggertype:"object"`
	Location        json.RawMessage  `gorm:"column:location;type:jsonb" json:"location" swaggertype:"object"` // latitude, longitude
	IsApproved      bool             `gorm:"column:is_approved" json:"is_approved"`
	Comments        []Comment        `gorm:"foreignKey:GuideID" json:"comments"`
}

func (g Guide) MarshalJSON() ([]byte, error) {
	type Alias Guide
	aux := Alias(g)
	if aux.Comments == nil {
		aux.Comments = make([]Comment, 0)
	}
	return json.Marshal(aux)
}
