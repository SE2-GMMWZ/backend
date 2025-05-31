package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuideRepository struct {
	db *gorm.DB
}

func NewGuideRepository(db *gorm.DB) *GuideRepository {
	return &GuideRepository{db}
}

func (r *GuideRepository) CreateGuide(guide *model.Guide) error {
	return r.db.Create(guide).Error
}

func (r *GuideRepository) GetGuideByID(id uuid.UUID) (*model.Guide, error) {
	var guide model.Guide
	err := r.db.First(&guide, "guide_id = ?", id).Error
	return &guide, err
}

func (r *GuideRepository) UpdateGuide(guide *model.Guide) error {
	return r.db.Save(guide).Error
}

func (r *GuideRepository) DeleteGuide(id uuid.UUID) error {
	return r.db.Delete(&model.Guide{}, "guide_id = ?", id).Error
}

func (r *GuideRepository) ListGuides(limit, offset int) ([]model.Guide, error) {
	var guides []model.Guide
	err := r.db.Order("guide_id asc").Limit(limit).Offset(offset).Find(&guides).Error
	return guides, err
}
