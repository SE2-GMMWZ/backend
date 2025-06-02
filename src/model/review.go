package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ReviewID      uuid.UUID `gorm:"column:review_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"review_id"`
	ReviewerID    uuid.UUID `gorm:"column:reviewer_id;type:uuid" json:"reviewer_id"`
	DockingSpotID uuid.UUID `gorm:"column:docking_spot_id;type:uuid" json:"docking_spot_id"`
	Rating        float64   `gorm:"column:rating" json:"rating"`
	Comment       *string   `gorm:"column:comment" json:"comment,omitempty"`
	DateOfReview  time.Time `gorm:"column:date_of_review" json:"date_of_review"`
}
