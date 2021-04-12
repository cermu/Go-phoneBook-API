package models

import "github.com/jinzhu/gorm"

// Account struct to store user information
// Account has many Contacts, AccountID is the foreign key
type Account struct {
	gorm.Model            // fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`will be added
	FirstName   string    `gorm:"size:15" json:"first_name"`
	LastName    string    `gorm:"size:15" json:"last_name"`
	Email       string    `gorm:"size:255;not null;unique;index:idx_email,unique" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(15);not null;unique;index:idx_phone,unique" json:"phone_number"`
	Password    string    `gorm:"type:varchar(255); not null" json:"password"`
	Contacts    []Contact `gorm:"ForeignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
