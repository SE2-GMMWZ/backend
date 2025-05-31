package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
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



func (r *PortRepository) ListPorts(limit, offset int, query string) ([]model.Port, error) {
	var ports []model.Port
	db := r.db.Model(&model.Port{})

	// Basic query parsing: support "name=...", "owner_id=...", "is_approved=..."
	if query != "" {
		parts := strings.Split(query, "&")
		for _, part := range parts {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := kv[0]
			value := kv[1]
			switch key {
			case "name":
				db = db.Where("name ILIKE ?", "%"+value+"%")
			case "owner_id":
				db = db.Where("owner_id = ?", value)
			case "is_approved":
				db = db.Where("is_approved = ?", value)
			}
		}
	}

	err := db.Order("port_id").Limit(limit).Offset(offset).Find(&ports).Error
	return ports, err
}

func (r *PortRepository) UpdatePort(port *model.Port) error {
	return r.db.Save(port).Error
}

func (r *PortRepository) DeletePort(id uuid.UUID) error {
	return r.db.Delete(&model.Port{}, "port_id = ?", id).Error
}
