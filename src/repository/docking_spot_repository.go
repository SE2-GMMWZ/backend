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
	err := r.db.First(&dock, "dock_id = ?", id).Error
	return &dock, err
}

func (r *DockingSpotRepository) UpdateDockingSpot(dock *model.DockingSpot) error {
	return r.db.Save(dock).Error
}

func (r *DockingSpotRepository) DeleteDockingSpot(id uuid.UUID) error {
	return r.db.Delete(&model.DockingSpot{}, "dock_id = ?", id).Error
}
func (r *DockingSpotRepository) ListDockingSpots(limit, offset int) ([]model.DockingSpot, error) {
	var spots []model.DockingSpot
	err := r.db.Limit(limit).Offset(offset).Find(&spots).Error
	return spots, err
}
