// Package models contains the User struct,
package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Models struct {
	User User
}

type User struct {
	ID    string `json:"id" gorm:"type=UUID;primary_key"`
	Email string `json:"email" gorm:"unique;not null" valid:"required~Email is required,email~Email must be a valid email address"`
	//Add pointers to nullable fields
	Firstname *string         `json:"firstname" gorm:"type:varchar(255)" valid:""`
	Lastname  *string         `json:"lastname" gorm:"type:varchar(255)" valid:""`
	Password  string          `json:"password" gorm:"type:varchar(255);not null" valid:"required~Password is required,minstringlength(8)~Password must be at least 8 characters long"`
	Tokens    []PersonalToken `json:"tokens" gorm:"foreignKey:UserID"`
}

type RegisterInput struct {
	Firstname string `json:"firstname" valid:"required~Firstname is required" `
	LastName  string `json:"lastname" valid:"required~Lastname is required"`
	Email     string `json:"email" valid:"email~Email is not a valid enail,required~Email is required"`
	Password  string `json:"password" valid:"required~Password is required,minstringlength(8)~Password must be at least 8 characters long"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := &User{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *User) GetUserByID(db *gorm.DB, id string) (*User, error) {
	user := &User{}
	err := db.Where("id = ?", id).First(user).Error
	return user, err
}

func GetTokensByUserID(db *gorm.DB, id string) ([]PersonalToken, error) {
	var tokens []PersonalToken
	if err := db.Where("user_id = ?", id).Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (u *User) ComparePassword(password string) bool {
	// u.Password should be the hashed password retrieved from the database
	fmt.Println(u.Password)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
