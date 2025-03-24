package model

// User struct maps to the existing "users" table
type User struct {
	ID               uint   `gorm:"column:id;primaryKey" json:"id"`
	Username         string `gorm:"column:username" json:"username"`
	Email            string `gorm:"column:email" json:"email"`
	Password         string `gorm:"column:password" json:"password"`
	Verified         bool   `gorm:"column:verified" json:"verified"`
	VerificationCode string `gorm:"column:verification_code" json:"verification_code"`
}
