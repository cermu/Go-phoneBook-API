package models

import (
	"github.com/badoux/checkmail"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/jinzhu/gorm"
	"log"
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
	contact.AccountID = accountId
	DBConnection.Table("contact").Create(contact)
	if contact.ID <= 0 {
		return utl.Message(105, "failed to save contact, tyr again")
	}

	response := utl.Message(0, "contact has been created")
	response["data"] = contact
	return response
}

// FetchContactsByAccountId public method that fetches contacts belonging to a specified account
func (contact *Contact) FetchContactsByAccountId(accountId uint) map[string]interface{} {
	// query contact table by account_id
	contacts := make([]*Contact, 0) // results will be stored in a slice of type Contact pointer
	err := DBConnection.Table("contact").Where("account_id=?", accountId).Find(&contacts).Error
	if err != nil {
		log.Printf("WARNING | An error occurred while fetching contacts for account: %d. Error: %v\n",
			accountId, err.Error())
		return utl.Message(105, "failed to fetch contacts, try again later")
	}

	// return the results
	response := utl.Message(0, "contacts fetched successfully")
	response["data"] = contacts
	return response
}

// FetchContactById public method that fetches a contact by its id passed in the URI
func (contact *Contact) FetchContactById(contactId uint) map[string]interface{} {
	// fetch contact from DB
	result := &Contact{}
	err := DBConnection.Table("contact").Where("id=?", contactId).First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utl.Message(104, "contact not found")
		}
		log.Printf("WARNING | An error occurred while fetching contact from the DB: %v\n", err.Error())
		return utl.Message(105, "failed to fetch contact, try again later.")
	}

	// return results
	response := utl.Message(0, "contact fetched successfully")
	response["data"] = result
	return response
}

// UpdateContact public method that is called to make updates to an existing contact record
func (contact *Contact) UpdateContact(contactId uint) map[string]interface{} {
	// validate email
	if contact.Email != "" {
		if err := checkmail.ValidateFormat(contact.Email); err != nil {
			return utl.Message(102, "email address is not valid")
		}
	}

	// validate phone number
	// should not be less than 9 digits and more than 12 chars
	// accepted: 0712345678, 254712345678, 712345678
	// store: 712345678
	if contact.PhoneNumber != "" {
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
	}

	// update contact record
	err := DBConnection.Table("contact").Model(contact).Where("id=?", contactId).Updates(contact).Error
	if err != nil {
		log.Printf("WARNING | An error occurred while updating contact: %v\n", err.Error())
		return utl.Message(105, "failed to update contact, try again later")
	}

	// fetch and return updated contact
	DBConnection.Table("contact").First(contact, contactId)
	response := utl.Message(0, "contact updated successfully")
	response["data"] = contact
	return response
}

// DeleteContact public method to remove a contact record from database
func (contact *Contact) DeleteContact(contactId uint) map[string]interface{} {
	/*
		Soft delete a record if there is a DeletedAt column. the column will only be updated with the deletion time.
		For permanent deletion, add `.Unscoped()` before .Delete()
	*/
	err := DBConnection.Table("contact").Where("id=?", contactId).Delete(contact).Error
	if err != nil {
		log.Printf("WARNING | An error has occurred while deleting contact: %v\n", err.Error())
		return utl.Message(105, "failed to delete contact, try again later")
	}
	return utl.Message(0, "contact deleted successfully")
}
