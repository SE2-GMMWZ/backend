package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PortRepository struct {
	db *gorm.DB
}

func NewPortRepository(db *gorm.DB) *PortRepository {
	return &PortRepository{db}
}

func (r *PortRepository) CreatePort(port *model.Port) error {
	return r.db.Create(port).Error
}

func (r *PortRepository) GetPortByID(id uuid.UUID) (*model.Port, error) {
	var port model.Port
	err := r.db.First(&port, "port_id = ?", id).Error
	return &port, err
}

func (r *PortRepository) UpdatePort(port *model.Port) error {
	return r.db.Save(port).Error
}

func (r *PortRepository) DeletePort(id uuid.UUID) error {
	return r.db.Delete(&model.Port{}, "port_id = ?", id).Error
}
