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

func (r *CommentRepository) ListComments(limit, offset int) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Order("comment_id").Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}
