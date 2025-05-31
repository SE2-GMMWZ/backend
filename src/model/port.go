package model

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Port struct {
	PortID       uuid.UUID        `gorm:"column:port_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"port_id"`
	Name         string           `gorm:"column:name" json:"name"`
	Location     json.RawMessage  `gorm:"column:location;type:jsonb" json:"location" swaggertype:"object"` // latitude, longitude, town
	Description  *string          `gorm:"column:description" json:"description,omitempty"`
	OwnerID      uuid.UUID        `gorm:"column:owner_id;type:uuid" json:"owner_id"`
	DockingSpots *json.RawMessage `gorm:"column:docking_spots;type:jsonb" json:"docking_spots,omitempty" swaggertype:"object"`
	Services     *json.RawMessage `gorm:"column:services;type:jsonb" json:"services,omitempty" swaggertype:"object"`
	IsApproved   bool             `gorm:"column:is_approved" json:"is_approved"`
}
