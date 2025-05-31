package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
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



func (r *CommentRepository) ListComments(limit, offset int, query string) ([]model.Comment, error) {
	var comments []model.Comment
	db := r.db.Model(&model.Comment{})

	// Basic query parsing: support "guide_id=...", "user_id=...", "content=..."
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
			case "guide_id":
				db = db.Where("guide_id = ?", value)
			case "user_id":
				db = db.Where("user_id = ?", value)
			case "content":
				db = db.Where("content ILIKE ?", "%"+value+"%")
			}
		}
	}

	err := db.Order("comment_id").Limit(limit).Offset(offset).Find(&comments).Error
	return comments, err
}
