package model

import (
	"encoding/json"
	"github.com/google/uuid"
)

type DockingSpot struct {
	DockID          uuid.UUID       `gorm:"column:dock_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"dock_id"`
	Name            string          `gorm:"column:name" json:"name"`
	Location        json.RawMessage `gorm:"column:location;type:jsonb" json:"location" swaggertype:"object"` // latitude, longitude, town
	Description     *string         `gorm:"column:description" json:"description,omitempty"`
	OwnerID         uuid.UUID       `gorm:"column:owner_id;type:uuid" json:"owner_id"`
	Services        *string         `gorm:"column:services" json:"services,omitempty"`
	ServicesPricing *float64        `gorm:"column:services_pricing" json:"services_pricing,omitempty"`
	PricePerNight   float64         `gorm:"column:price_per_night" json:"price_per_night"`
	PricePerPerson  *float64        `gorm:"column:price_per_person" json:"price_per_person,omitempty"`
	Availability    string          `gorm:"column:availability" json:"availability"` // 'available' or 'unavailable'
}
