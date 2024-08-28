// Package models contains the User struct,
package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	User User
}

type User struct {
	ID        string          `json:"id" gorm:"type=UUID;primary_key"`
	Email     string          `json:"email" gorm:"unique;not null" valid:"required~Email is required,email~Email must be a valid email address"`
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

func (u *User) ComparePassword(password string) bool {
	// u.Password should be the hashed password retrieved from the database
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
