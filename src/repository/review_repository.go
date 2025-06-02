package repository

import (
	"backend/src/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db}
}

func (r *ReviewRepository) CreateReview(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) GetReviewByID(id uuid.UUID) (*model.Review, error) {
	var review model.Review
	err := r.db.First(&review, "review_id = ?", id).Error
	return &review, err
}

func (r *ReviewRepository) UpdateReview(review *model.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) DeleteReview(id uuid.UUID) error {
	return r.db.Delete(&model.Review{}, "review_id = ?", id).Error
}

func (r *ReviewRepository) ListReviews(limit, page int, reviewerID string, minRating *float64, comment string, dockingSpotID string) ([]model.Review, int, int, error) {
	var reviews []model.Review
	db := r.db.Model(&model.Review{})

	// Apply filters based on provided parameters
	if reviewerID != "" {
		db = db.Where("reviewer_id = ?", reviewerID)
	}
	if minRating != nil {
		db = db.Where("rating >= ?", *minRating)
	}
	if comment != "" {
		db = db.Where("comment ILIKE ?", "%"+comment+"%")
	}
	if dockingSpotID != "" {
		db = db.Where("docking_spot_id = ?", dockingSpotID)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("review_id").Limit(limit).Offset(offset).Find(&reviews).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return reviews, page, totalPages, err
}
