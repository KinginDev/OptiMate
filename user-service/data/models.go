package data

import "gorm.io/gorm"

type Models struct {
	User User
}

type User struct {
	ID       int             `json:"id" gorm:"type=UUID;primary_key"`
	Username string          `json:"username" gorm:"unique;not null"`
	Email    string          `json:"email" gorm:"unique;not null"`
	Password string          `json:"password" gorm:"not null"`
	Tokens   []PersonalToken `json:"tokens" gorm:"foreignKey:UserID"`
}

type PersonalToken struct {
	ID        int    `json:"id" gorm:"type=UUID;primary_key"`
	UserID    int    `json:"user_id" gorm:"not null"`
	Token     string `json:"token" gorm:"unique;not null"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}
