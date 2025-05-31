package repository

import (
	"backend/src/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "user_id = ?", id).Error
	return &user, err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "user_id = ?", id).Error
}

func (r *UserRepository) ListUsers(limit, offset int, query string) ([]model.User, error) {
	var users []model.User
	db := r.db.Model(&model.User{})

	if query != "" {
		likePattern := "%" + query + "%"
		db = db.Where(
			"email ILIKE ? OR name ILIKE ? OR surname ILIKE ? OR phone_number ILIKE ?",
			likePattern, likePattern, likePattern, likePattern,
		)
	}

	err := db.Order("user_id ASC").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
