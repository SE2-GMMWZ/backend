package model

import (
	"github.com/google/uuid"
	"time"
)

type Booking struct {
	BookingID     uuid.UUID `gorm:"column:booking_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"booking_id"`
	SailorID      uuid.UUID `gorm:"column:sailor_id;type:uuid" json:"sailor_id"`
	DockID        uuid.UUID `gorm:"column:dock_id;type:uuid" json:"dock_id"`
	StartDate     time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate       time.Time `gorm:"column:end_date" json:"end_date"`
	PaymentMethod string    `gorm:"column:payment_method" json:"payment_method"`
	PaymentStatus string    `gorm:"column:payment_status" json:"payment_status"`
	People        int       `gorm:"column:people" json:"people"`
}
