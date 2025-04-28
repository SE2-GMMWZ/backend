package model

import "github.com/google/uuid"

// User struct maps to the existing "users" table
type User struct {
	UserID      uuid.UUID `gorm:"column:user_id;primaryKey;type:uuid;default:uuid_generate_v4()" json:"user_id"`
	Name        string    `gorm:"column:name" json:"name"`
	Surname     string    `gorm:"column:surname" json:"surname"`
	Email       string    `gorm:"column:email;unique" json:"email"`
	PhoneNumber *string   `gorm:"column:phone_number" json:"phone_number,omitempty"`
	Password    string    `gorm:"column:password" json:"-"`
	Role        string    `gorm:"column:role" json:"role"`
}
