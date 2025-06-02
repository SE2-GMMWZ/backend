package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	err := r.db.Preload("Reviews").First(&dock, "dock_id = ?", id).Error
	return &dock, err
}

func (r *DockingSpotRepository) UpdateDockingSpot(dock *model.DockingSpot) error {
	return r.db.Save(dock).Error
}

func (r *DockingSpotRepository) DeleteDockingSpot(id uuid.UUID) error {
	return r.db.Delete(&model.DockingSpot{}, "dock_id = ?", id).Error
}

func (r *DockingSpotRepository) ListDockingSpots(limit, page int, name, ownerID, availability string) ([]model.DockingSpot, int, int, error) {
	var spots []model.DockingSpot
	db := r.db.Model(&model.DockingSpot{})

	// Apply filters based on provided parameters
	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if ownerID != "" {
		db = db.Where("owner_id = ?", ownerID)
	}
	if availability != "" {
		db = db.Where("availability = ?", availability)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("dock_id").Limit(limit).Offset(offset).Find(&spots).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return spots, page, totalPages, err
}
