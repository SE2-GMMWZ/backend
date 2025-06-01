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



func (r *PortRepository) ListPorts(limit, page int, name, ownerID string, isApproved *bool) ([]model.Port, int, int, error) {
	var ports []model.Port
	db := r.db.Model(&model.Port{})

	// Apply filters based on provided parameters
	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if ownerID != "" {
		db = db.Where("owner_id = ?", ownerID)
	}
	if isApproved != nil {
		db = db.Where("is_approved = ?", *isApproved)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("port_id").Limit(limit).Offset(offset).Find(&ports).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return ports, page, totalPages, err
}

func (r *PortRepository) UpdatePort(port *model.Port) error {
	return r.db.Save(port).Error
}

func (r *PortRepository) DeletePort(id uuid.UUID) error {
	return r.db.Delete(&model.Port{}, "port_id = ?", id).Error
}
