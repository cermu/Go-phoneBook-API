package models

import (
	"github.com/badoux/checkmail"
	"github.com/cermu/Go-phoneBook-API/auth"
	"github.com/cermu/Go-phoneBook-API/middlewares"
	utl "github.com/cermu/Go-phoneBook-API/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

/* Account struct to store user information
Account has many Contacts, AccountID is the foreign key
*/
type Account struct {
	gorm.Model            // fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`will be added
	FirstName   string    `gorm:"size:15" json:"first_name"`
	LastName    string    `gorm:"size:15" json:"last_name"`
	Email       string    `gorm:"size:255;not null;unique;index:idx_email" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(15);not null;unique;index:idx_phone" json:"phone_number"`
	Password    string    `gorm:"type:varchar(255); not null" json:"password"`
	Contacts    []Contact `gorm:"ForeignKey:AccountID"`
}

/* LoginDetails struct used to fetch login credentials
from json request
*/
type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshToken struct to fetch refresh_token from json request
type MapRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// validateAccountData private method to be used to validate
// incoming requests to create an account
func (account *Account) validateAccountData() (map[string]interface{}, bool) {
	// check for empty data
	if account.FirstName == "" || account.LastName == "" || account.Email == "" ||
		account.PhoneNumber == "" || account.Password == "" {
		return utl.Message(102, "the following fields are required: first_name, "+
			"last_name, email, phone_number, password"), false
	}

	// validate email
	if err := checkmail.ValidateFormat(account.Email); err != nil {
		return utl.Message(102, "provide a valid email address"), false
	}

	// check password
	if len(account.Password) < 6 {
		return utl.Message(102, "password should not be less than six characters"), false
	}

	// email address and phone number must be unique
	tmp := &Account{}
	emailErr := DBConnection.Table("account").Where("email=?", account.Email).First(tmp).Error
	if emailErr != nil && emailErr != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while validatin email address: %v\n", emailErr.Error())
		return utl.Message(105, "failed to validate email, try again later"), false
	}

	if tmp.Email != "" {
		return utl.Message(101, "email address already exists"), false
	}

	phoneErr := DBConnection.Table("account").Where("phone_number=?", account.PhoneNumber).First(tmp).Error
	if phoneErr != nil && phoneErr != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while validatin phone number: %v\n", phoneErr.Error())
		return utl.Message(105, "failed to validate phone number, try again later"), false
	}

	if tmp.PhoneNumber != "" {
		return utl.Message(101, "phone number already exists"), false
	}

	return utl.Message(0, "account data validated successfully"), true
}

// CreateAccount public method that is used to create a new account
func (account *Account) CreateAccount() map[string]interface{} {
	if resp, ok := account.validateAccountData(); !ok {
		return resp
	}

	// hash the password before storing it
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	DBConnection.Create(account) // save account in the DB

	if account.ID <= 0 {
		return utl.Message(105, "failed to save account, try again")
	}

	// remove the password
	account.Password = ""

	response := utl.Message(0, "account has been created")
	response["data"] = account
	return response
}

// Login public function to authenticate users
func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := DBConnection.Table("account").Where("email=?", email).First(account).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utl.Message(104, "account does not exist, check credentials or create one")
		}
		return utl.Message(105, "failed to fetch account, try again")
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return utl.Message(102, "invalid login credentials, try again")
	}

	// create tokens
	authDetails, tokenErr := auth.CreateToken(account.ID)
	if tokenErr != nil {
		return utl.Message(105, "failed to create authentication tokens, try again")
	}

	// save JWT metadata in redis
	redisSaveErr := auth.SaveJWTMetadata(account.ID, authDetails)
	if redisSaveErr != nil {
		log.Printf("WARNING | An error occurred while saving to redis: %v\n", redisSaveErr)
		return utl.Message(105, "failed to create authentication tokens, try again")
	}

	tokens := map[string]string{
		"access_token":  authDetails.AccessToken,
		"refresh_token": authDetails.RefreshToken,
	}
	response := utl.Message(0, "authentication successful")
	response["tokens"] = tokens
	return response
}

// FetchAccount public method that takes in an account id and returns
// account data
func (account *Account) FetchAccount(accountId uint) map[string]interface{} {
	err := DBConnection.Table("account").Where("id=?", accountId).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while fetching account from database: %v\n", err)
		return utl.Message(105, "failed to fetch account details, try again")
	}

	// account not found
	if account.Email == "" {
		return utl.Message(104, "account details not found")
	}

	// remove the password
	account.Password = ""
	response := utl.Message(0, "account details fetched successfully")
	response["data"] = account
	return response
}

// Logout public function to delete auth details and log out a user
func Logout(req *http.Request) map[string]interface{} {
	accessDetails, err := middlewares.ExtractTokenFromRequest(req)
	if err != nil {
		log.Printf("WARNING | The following error occurred while logging out: %v\n", err.Error())
		return utl.Message(105, "failed to log out, try again")
	}

	authDeleted, authDelErr := auth.DeleteAuthenticationDetails(accessDetails.AccessUuid)
	if authDelErr != nil || authDeleted == 0 {
		log.Printf("WARNING | The following error occurred while logging out: %v\n", authDelErr)
		return utl.Message(105, "failed to log out, try again")
	}
	response := utl.Message(0, "logged out successfully")
	return response
}
