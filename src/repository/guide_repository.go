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



func (r *GuideRepository) ListGuides(limit, page int, title, authorID string, isApproved *bool) ([]model.Guide, int, int, error) {
	var guides []model.Guide
	db := r.db.Model(&model.Guide{})

	// Apply filters based on provided parameters
	if title != "" {
		db = db.Where("title ILIKE ?", "%"+title+"%")
	}
	if authorID != "" {
		db = db.Where("author_id = ?", authorID)
	}
	if isApproved != nil {
		db = db.Where("is_approved = ?", *isApproved)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("guide_id asc").Limit(limit).Offset(offset).Find(&guides).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return guides, page, totalPages, err
}
