// Package models
package models

// PersonalToken is a model for the personal token
type PersonalToken struct {
	ID        string `json:"id" gorm:"type=UUID;primary_key"`
	UserID    string `json:"user_id" gorm:"not null"`
	Token     string `json:"token" gorm:"unique;not null"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Revoked   bool   `json:"revoked"`
}
