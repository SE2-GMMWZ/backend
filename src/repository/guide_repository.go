package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
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



func (r *GuideRepository) ListGuides(limit, offset int, query string) ([]model.Guide, error) {
	var guides []model.Guide
	db := r.db.Model(&model.Guide{})

	// Basic query parsing: support "title=...", "author_id=...", "is_approved=..."
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
			case "title":
				db = db.Where("title ILIKE ?", "%"+value+"%")
			case "author_id":
				db = db.Where("author_id = ?", value)
			case "is_approved":
				db = db.Where("is_approved = ?", value)
			}
		}
	}

	err := db.Order("guide_id asc").Limit(limit).Offset(offset).Find(&guides).Error
	return guides, err
}
