package repository

import (
	"backend/src/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db}
}

func (r *CommentRepository) CreateComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetCommentByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, "comment_id = ?", id).Error
	return &comment, err
}

func (r *CommentRepository) UpdateComment(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *CommentRepository) DeleteComment(id uuid.UUID) error {
	return r.db.Delete(&model.Comment{}, "comment_id = ?", id).Error
}

func (r *CommentRepository) ListComments(limit, page int, guideID, userID, content string) ([]model.Comment, int, int, error) {
	var comments []model.Comment
	db := r.db.Model(&model.Comment{})

	// Apply filters based on provided parameters
	if guideID != "" {
		db = db.Where("guide_id = ?", guideID)
	}
	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	if content != "" {
		db = db.Where("content ILIKE ?", "%"+content+"%")
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("comment_id").Limit(limit).Offset(offset).Find(&comments).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return comments, page, totalPages, err
}
