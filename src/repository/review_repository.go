package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"strconv"
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



func (r *ReviewRepository) ListReviews(limit, offset int, query string) ([]model.Review, error) {
	var reviews []model.Review
	db := r.db.Model(&model.Review{})

	// Basic query parsing: support "reviewer_id=...", "rating=...", "comment=..."
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
			case "reviewer_id":
				db = db.Where("reviewer_id = ?", value)
			case "rating":
				// Support exact or minimum rating
				if rating, err := strconv.ParseFloat(value, 64); err == nil {
					db = db.Where("rating >= ?", rating)
				}
			case "comment":
				db = db.Where("comment ILIKE ?", "%"+value+"%")
			}
		}
	}

	err := db.Order("review_id").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}
