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

func (r *UserRepository) ListUsers(limit, page int, email, name, surname, phoneNumber, role, query string) ([]model.User, int, int, error) {
	var users []model.User
	db := r.db.Model(&model.User{})

	// Apply specific field filters
	if email != "" {
		db = db.Where("email ILIKE ?", "%"+email+"%")
	}
	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname != "" {
		db = db.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if phoneNumber != "" {
		db = db.Where("phone_number ILIKE ?", "%"+phoneNumber+"%")
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}

	// Apply general search query across multiple fields
	if query != "" {
		likePattern := "%" + query + "%"
		db = db.Where(
			"email ILIKE ? OR name ILIKE ? OR surname ILIKE ? OR phone_number ILIKE ?",
			likePattern, likePattern, likePattern, likePattern,
		)
	}

	offset := (page - 1) * limit
	var totalRecords int64
	db.Count(&totalRecords)

	err := db.Order("user_id ASC").Limit(limit).Offset(offset).Find(&users).Error
	totalPages := int((totalRecords + int64(limit) - 1) / int64(limit))
	return users, page, totalPages, err
}
