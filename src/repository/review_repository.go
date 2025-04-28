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

func (r *ReviewRepository) ListReviews(limit, offset int) ([]model.Review, error) {
	var reviews []model.Review
	err := r.db.Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}
