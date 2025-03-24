package repository

import (
	"backend/src/model"
	"errors"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByEmail(username string) (*model.User, error)
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	var existingUser model.User
	if err := r.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user with this username or email already exists")
	}

	if err := r.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("email = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	var existingUser model.User
	if err := r.DB.First(&existingUser, user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := r.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}
