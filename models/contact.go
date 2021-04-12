package models

import "github.com/jinzhu/gorm"

// Contact struct to store contact information
type Contact struct {
	gorm.Model         // fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`will be added
	FirstName   string `gorm:"size:15" json:"first_name"`
	LastName    string `gorm:"size:15" json:"last_name"`
	PhoneNumber string `gorm:"type:varchar(15);not null;index:idx_phone" json:"phone_number"`
	Email       string `gorm:"size:255;not null;index:idx_email" json:"email"`
	AccountID   uint   `json:"account_id"` // this is a foreign_key from the account table
}
