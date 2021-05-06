package models

import (
	"github.com/badoux/checkmail"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/jinzhu/gorm"
	"strings"
)

// Contact struct to store contact information
type Contact struct {
	gorm.Model         // fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`will be added
	FirstName   string `gorm:"size:15" json:"first_name"`
	LastName    string `gorm:"size:15" json:"last_name"`
	PhoneNumber string `gorm:"type:varchar(15);not null" json:"phone_number"`
	Email       string `gorm:"size:255;not null" json:"email"`
	AccountID   uint   `gorm:"not null" json:"account_id"` // this is a foreign_key from the account table
}

// CreateContact public method that allows a user/account to create/save a contact
func (contact *Contact) CreateContact(accountId uint) map[string]interface{} {
	// check for empty data
	if contact.PhoneNumber == "" || contact.Email == "" {
		return utl.Message(102, "the following fields are required: phone_number and email")
	}

	// validate email
	if err := checkmail.ValidateFormat(contact.Email); err != nil {
		return utl.Message(102, "email address is not valid")
	}

	// validate phone number
	// should not be less than 9 digits and more than 12 chars
	// accepted: 0712345678, 254712345678, 712345678
	// store: 712345678
	if len(contact.PhoneNumber) > 12 || len(contact.PhoneNumber) < 9 {
		return utl.Message(102, "enter a valid phone number, between 9 to 12 digits.")
	}

	if strings.HasPrefix(contact.PhoneNumber, "254") {
		phoneNumber_ := contact.PhoneNumber[3:len(contact.PhoneNumber)]
		contact.PhoneNumber = phoneNumber_
	}

	if strings.HasPrefix(contact.PhoneNumber, "0") {
		phoneNumber_ := contact.PhoneNumber[1:len(contact.PhoneNumber)]
		contact.PhoneNumber = phoneNumber_
	}

	// save the contact in DB
	contact.ID = accountId
	DBConnection.Table("contact").Create(contact)
	if contact.ID <= 0 {
		return utl.Message(105, "failed to save contact, tyr again")
	}

	response := utl.Message(0, "contact has been created")
	response["data"] = contact
	return response
}
