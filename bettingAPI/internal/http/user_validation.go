package http

import (
	"bettingAPI/internal/mysql"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func isLoginRequestValid(logReq loginRequest) bool {
	if (isUsernameValid(logReq.User) || isEmailValid(logReq.User)) && isPasswordValid(logReq.Password) {
		return true
	}
	return false
}

func (h *handler) Login(user string, password string) (*mysql.User, error) {
	var u *mysql.User
	var err error
	if IsEmail(user) {
		u, err = h.db.FindUserByEmail(user)
	} else {
		u, err = h.db.FindUserByUsername(user)
	}
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	if checkPasswordHash(password, u.PasswordHash) {
		return u, nil
	}
	return nil, nil
}

func IsEmail(user string) bool {
	return strings.Contains(user, "@")
}

func isRegistrationRequestValid(regReq registrationRequest) bool {
	if !isEmailValid(regReq.Email) {
		return false
	}
	if !isUsernameValid(regReq.Username) {
		return false
	}
	if !isPasswordValid(regReq.Password) {
		return false
	}
	if !isNameValid(regReq.FirstName, regReq.LastName) {
		return false
	}
	if !isBirthDateValid(regReq.BirthDate) {
		return false
	}
	return true
}

func (h *handler) isEmailUnique(email string) (bool, error) {
	user, err := h.db.FindUserByEmail(email)
	if err != nil {
		return false, err
	}
	return user == nil, nil
}

func (h *handler) isUsernameUnique(username string) (bool, error) {
	user, err := h.db.FindUserByUsername(username)
	if err != nil {
		return false, err
	}
	return user == nil, nil
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isUsernameValid(username string) bool {
	if 3 >= len(username) || len(username) >= 20 {
		return false
	}
	for _, char := range username {
		if !unicode.IsLetter(char) {
			if unicode.IsNumber(char) {
				return false
			}
		}
	}
	return true
}

func isPasswordValid(password string) bool {
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
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
	if 7 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdNoSpaces || !pswdLength {
		return false
	}
	return true
}

func isNameValid(firstName string, lastName string) bool {
	for _, char := range firstName {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	for _, char := range lastName {
		if !unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func isBirthDateValid(birthDate string) bool {
	bd, err := time.Parse("02-01-2006", birthDate)
	if err != nil {
		log.Println("Error parsing time")
		return false
	}
	currentTime := time.Now().AddDate(-18, 0, 0)
	return bd.Before(currentTime)
}

func hashPassword(password string) (string, error) { // isto treba vracat error, status 500
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
