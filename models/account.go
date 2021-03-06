package models

import (
	"fmt"
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
	Active      bool      `gorm:"default:true" json:"active"`
	Contacts    []Contact `gorm:"ForeignKey:AccountID" json:"contacts"`
}

/* LoginDetails struct used to fetch login credentials
from json request
*/
type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* UpdateAccountDetails struct used to fetch account credentials
from json request in order to update an existing account
*/
type UpdateAccountDetails struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// RefreshToken struct to fetch refresh_token from json request
type MapRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// ChangePassword struct to fetch new password from json request
type ChangePassword struct {
	Password      string `json:"password"`
	PasswordAgain string `json:"password_again"`
}

// ResetPassword struct to fetch account's email from json request
type ResetPassword struct {
	Email string `json:"email"`
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
	err := DBConnection.Table("account").Where("email=? AND active=?", email, true).First(account).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utl.Message(104, "account deactivated or it does not exist, "+
				"check credentials, create one or request for reactivation")
		}
		log.Printf("WARNING | An error occurred while fetching account from database to login/authenticate "+
			"it: %v\n", err)
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
		"type":          authDetails.TokenType,
	}
	response := utl.Message(0, "authentication successful")
	response["tokens"] = tokens
	return response
}

// FetchAccount public method that takes in an account id and returns
// account data
func (account *Account) FetchAccount(accountId uint) map[string]interface{} {
	err := DBConnection.Table("account").Where("id=?", accountId).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utl.Message(104, "account not found")
		}
		log.Printf("WARNING | An error occurred while fetching account from database: %v\n", err)
		return utl.Message(105, "failed to fetch account details, try again")
	}

	// account not found
	if account.Email == "" {
		return utl.Message(104, "account details not found")
	}

	// remove the password
	account.Password = ""
	response := utl.Message(0, "your account details have been fetched successfully")
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

// DeactivateAccount public method that set's an account to in active
func (account *Account) DeactivateAccount(req *http.Request) map[string]interface{} {
	// create a channel
	ch := make(chan int, 1)

	// close the channel after use to avoid memory leaks
	defer func() {
		close(ch)
		ch = nil
	}()

	// delete access token details from redis using go routines
	go func() {
		accessDetails, err := middlewares.ExtractTokenFromRequest(req)
		if err != nil {
			ch <- 0
		}

		deleted, delErr := auth.DeleteAuthenticationDetails(accessDetails.AccessUuid)
		if delErr != nil || deleted == 0 {
			ch <- 0
		}
		ch <- 1
	}()

	accountId := req.Context().Value("account").(uint) // fetch account id from context
	err := DBConnection.Table("account").Where("id=?", accountId).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while fetching account from database to deactivate it: %v\n", err)
		return utl.Message(105, "failed to deactivate account, try again")
	}

	// account not found
	if account.Email == "" {
		return utl.Message(104, "account details not found")
	}

	// deactivate an account
	account.Active = false
	DBConnection.Model(&account).Where("active=?", true).Update("active", false)

	// block until a value is received in the channel
	val := <-ch
	if val == 0 {
		return utl.Message(100, "account deactivated successfully but redis details were not cleared")
	}
	return utl.Message(0, "account deactivated successfully")
}

// UpdateAccount public function that is used to make updates to an
// existing account
func UpdateAccount(updateAccount *UpdateAccountDetails, accountId uint) map[string]interface{} {
	account := &Account{}

	// check for empty data
	if updateAccount.FirstName == "" || updateAccount.LastName == "" || updateAccount.Email == "" ||
		updateAccount.PhoneNumber == "" {
		return utl.Message(102, "the following fields are required: first_name, "+
			"last_name, email, phone_number")
	}

	// validate email
	if err := checkmail.ValidateFormat(updateAccount.Email); err != nil {
		return utl.Message(102, "provide a valid email address")
	}

	// email address and phone number must be unique
	tmp := &Account{}
	emailErr := DBConnection.Table("account").Where("email=? AND id NOT IN (?)",
		updateAccount.Email, []uint{accountId}).First(tmp).Error
	if emailErr != nil && emailErr != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while validating email address: %v\n", emailErr.Error())
		return utl.Message(105, "failed to validate email, try again later")
	}

	if tmp.Email != "" {
		return utl.Message(101, "email address already exists")
	}

	phoneErr := DBConnection.Table("account").Where("phone_number=? AND id NOT IN (?)",
		updateAccount.PhoneNumber, []uint{accountId}).First(tmp).Error
	if phoneErr != nil && phoneErr != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while validating phone number: %v\n", phoneErr.Error())
		return utl.Message(105, "failed to validate phone number, try again later")
	}

	if tmp.PhoneNumber != "" {
		return utl.Message(101, "phone number already exists")
	}

	// update the account
	DBConnection.Model(account).Where("id=?", accountId).Updates(map[string]interface{}{"first_name": updateAccount.FirstName,
		"last_name": updateAccount.LastName, "email": updateAccount.Email, "phone_number": updateAccount.PhoneNumber})

	// fetch and return account
	DBConnection.First(account, accountId)
	account.Password = ""
	response := utl.Message(0, "account has been updated")
	response["data"] = account
	return response
}

// ChangePassword public method used to change an account's password
func (account *Account) ChangePassword(changePassword *ChangePassword, accountId uint) map[string]interface{} {
	// validate the passwords in request
	if changePassword.Password == "" || changePassword.PasswordAgain == "" {
		return utl.Message(102, "the following fields are required, password, password_again")
	}
	if changePassword.PasswordAgain != changePassword.Password {
		return utl.Message(102, "password change failed, passwords entered did not match")
	}

	// fetch the account from DB
	err := DBConnection.Table("account").Where("id=? AND active=?", accountId, true).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while fetching account from database to change its password it: %v\n", err)
		return utl.Message(105, "password change failed, try again")
	}

	if account.Email == "" {
		return utl.Message(104, "account is deactivated or it does not exist")
	}

	// bcrypt the new password and update the existing one
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(changePassword.PasswordAgain), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	DBConnection.Model(&account).Update("password", string(hashedPassword))

	// respond to the request
	return utl.Message(0, "password changed successfully")
}

// SendResetPasswordLink public method used to reset an account's password
func SendResetPasswordLink(resetPassword *ResetPassword) map[string]interface{} {
	// check for empty data
	if resetPassword.Email == "" {
		return utl.Message(102, "the following field is required: email")
	}

	// validate email
	if err := checkmail.ValidateFormat(resetPassword.Email); err != nil {
		return utl.Message(102, "provide a valid email address")
	}

	// check if account exists
	account := &Account{}
	err := DBConnection.Table("account").Where("email=? AND active=?",
		resetPassword.Email, true).First(&account).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while fetching account from database to send a password reset link: %v\n", err)
		return utl.Message(105, "Sending password reset link had failed, try again later")
	}

	if account.Email == "" {
		return utl.Message(105, "Sending password reset link has failed, provided email does not exist")
	}

	// send email
	done := make(chan bool, 1)
	go func() {
		// generate reset link
		resetLinkMeta, rlErr := auth.GenerateResetPasswordLink()
		if rlErr != nil {
			log.Printf("WARNING | An error occurred while generating reset password link in SendResetPasswordLink method: %v\n", rlErr)
			done <- false
		}

		// save metadata to redis
		saveErr := auth.SaveResetLinkMetadata(account.ID, resetLinkMeta)
		if saveErr != nil {
			log.Printf("WARNING | An error occurred while saving to redis in SendResetPasswordLink method: %v\n", saveErr)
			done <- false
		}

		resetLinkToken := resetLinkMeta.RandomString
		resetLink := fmt.Sprintf("http://localhost:8081/phonebookapi/v1/reset/password/%s", resetLinkToken)

		// send the actual email here
		log.Printf("INFO | An email has been sent to: %v with link: %v", account.Email, resetLink)
		done <- true
	}()
	// <-done
	emailSent := <-done
	if !emailSent {
		return utl.Message(105, "Sending reset password link has failed, try again later")
	}
	return utl.Message(0, "An email has been sent with instructions to reset your password")
}

// ResetPassword public method used to reset an account's password
func ResetAccountPassword(resetLinkToken string, changePassword *ChangePassword) map[string]interface{} {
	// validate the passwords in request
	if changePassword.Password == "" || changePassword.PasswordAgain == "" {
		return utl.Message(102, "the following fields are required, password, password_again")
	}
	if changePassword.PasswordAgain != changePassword.Password {
		return utl.Message(102, "password reset failed, passwords entered did not match")
	}

	// fetch metadata from redis
	accountId, err := utl.RedisClient().Get(resetLinkToken).Result()
	if err != nil {
		log.Printf("WARNING | An error has occurred while fetching data from redis in ResetAccountPassword method: %v\n", err.Error())
		return utl.Message(106, "password reset link has expired")
	}

	// fetch account
	account := &Account{}
	fetchErr := DBConnection.Table("account").Where("id=? AND active=?", accountId, true).First(&account).Error
	if fetchErr != nil && err != gorm.ErrRecordNotFound {
		log.Printf("WARNING | An error occurred while fetching account from database to reset its password it: %v\n", err)
		return utl.Message(105, "password reset failed, try again")
	}

	if account.Email == "" {
		return utl.Message(104, "account is deactivated or it does not exist")
	}

	// reset account's password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(changePassword.PasswordAgain), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)
	DBConnection.Model(&account).Update("password", string(hashedPassword))

	// delete password reset meta from redis
	_, delErr := auth.DeleteAuthenticationDetails(resetLinkToken)
	if delErr != nil {
		log.Printf("WARNING | An error occurred while deleting reset link metadata from redis: %v\n", delErr)
	}

	return utl.Message(0, "password has been reset successfully")
}
