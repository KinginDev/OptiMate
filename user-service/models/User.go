package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Models struct {
	User User
}

type User struct {
	ID       string          `json:"id" gorm:"type=UUID;primary_key"`
	Email    string          `json:"email" gorm:"unique;not null" valid:"required~Email is required,email~Email must be a valid email address"`
	Password string          `json:"password" gorm:"not null" valid:"required~Password is required,minstringlength(8)~Password must be at least 8 characters long"`
	Tokens   []PersonalToken `json:"tokens" gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) CreateUser(db *gorm.DB, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: string(hashedPassword),
	}

	return db.Create(user).Error
}
