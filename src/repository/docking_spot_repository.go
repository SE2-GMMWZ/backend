package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type DockingSpotRepository struct {
	db *gorm.DB
}

func NewDockingSpotRepository(db *gorm.DB) *DockingSpotRepository {
	return &DockingSpotRepository{db}
}

func (r *DockingSpotRepository) CreateDockingSpot(dock *model.DockingSpot) error {
	return r.db.Create(dock).Error
}

func (r *DockingSpotRepository) GetDockingSpotByID(id uuid.UUID) (*model.DockingSpot, error) {
	var dock model.DockingSpot
	err := r.db.First(&dock, "dock_id = ?", id).Error
	return &dock, err
}

func (r *DockingSpotRepository) UpdateDockingSpot(dock *model.DockingSpot) error {
	return r.db.Save(dock).Error
}

func (r *DockingSpotRepository) DeleteDockingSpot(id uuid.UUID) error {
	return r.db.Delete(&model.DockingSpot{}, "dock_id = ?", id).Error
}

func (r *DockingSpotRepository) ListDockingSpots(limit, offset int, query string) ([]model.DockingSpot, error) {
	var spots []model.DockingSpot
	db := r.db.Model(&model.DockingSpot{})

	// Basic query parsing: support "name=...", "owner_id=...", "availability=..."
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
			case "availability":
				db = db.Where("availability = ?", value)
			}
		}
	}

	err := db.Order("dock_id").Limit(limit).Offset(offset).Find(&spots).Error
	return spots, err
}
