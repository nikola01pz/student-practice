package http

import (
	"errors"
	"log"
	"net/mail"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func (regReq *registrationRequest) validateUsername() error {
	if 3 >= len(regReq.Username) || len(regReq.Username) >= 20 {
		return errors.New("username length must be greater than 3 and less than 20 characters")
	}
	for _, char := range regReq.Username {
		if !unicode.IsLetter(char) {
			if unicode.IsNumber(char) {
				return errors.New("only letters are allowed for username")
			}
		}
	}
	return nil
}

func (regReg *registrationRequest) validatePassword() error {
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range regReg.Password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 7 < len(regReg.Password) && len(regReg.Password) < 60 {
		pswdLength = true
	}
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdNoSpaces || !pswdLength {
		return errors.New("password does not meet the requirements")
	}
	return nil
}

func (regReq *registrationRequest) validateEmail() error {
	_, err := mail.ParseAddress(regReq.Email)
	if err != nil {
		return errors.New("email is not valid")
	}
	return nil
}

func (regReq *registrationRequest) validateName() error {
	for _, char := range regReq.FirstName {
		if !unicode.IsLetter(char) {
			return errors.New("firstname does not meet the requirements")
		}
	}
	for _, char := range regReq.LastName {
		if !unicode.IsLetter(char) {
			return errors.New("lastname does not meet the requirements")
		}
	}
	return nil
}

func (regReq *registrationRequest) validateBirthDate() (time.Time, error) {
	currentTime := time.Now()
	birthDate, err := time.Parse("02-01-2006", regReq.BirthDate)
	if err != nil {
		log.Println("Error parsing time")
	}
	timeSpan := currentTime.Sub(birthDate)
	if timeSpan.Hours() < 157788 {
		return birthDate, errors.New("you are too young to register")
	}
	return birthDate, nil
}

func (h *handler) checkUserAndPassword(logReq loginRequest) error {
	if (h.db.IsUsernameUsed(logReq.User)) == nil && (h.db.IsEmailUsed(logReq.User) == nil) {
		return errors.New("wrong username, email or password")
	}
	var passwordHashFromDB = h.db.StoredUserPassword(logReq.User)
	if !checkPasswordHash(logReq.Password, passwordHashFromDB) {
		return errors.New("wrong username, email or password")
	}
	return nil
}

func hashPassword(password string) []byte {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password")
	}
	return hash
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
