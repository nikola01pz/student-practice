package http

import (
	"bettingAPI/internal/mysql"
	"errors"
	"log"
	"net/mail"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func validateLoginRequest(logReq loginRequest) (bool, error) {
	if IsUsernameValid(logReq.User) {
		if IsPasswordValid(logReq.Password) {
			return true, nil
		}
	} else {
		isValid, err := IsEmailValid(logReq.User)
		if err != nil {
			return false, err
		}
		if isValid {
			if IsPasswordValid(logReq.User) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (h *handler) validateUserAndPassword(logReq loginRequest) (mysql.User, error) { // treba dohvatit cijelog usera pa odmah na cijelom useru provjeriti i password
	user, err := h.db.FindUserByUsername(logReq.User)
	if err != nil {
		user = mysql.User{}
	}
	if user.Username == logReq.User {
		if checkPasswordHash(logReq.Password, user.PasswordHash) {
			return user, nil
		} else {
			return mysql.User{}, nil
		}
	}

	user, err = h.db.FindUserByEmail(logReq.User)
	if err != nil {
		user = mysql.User{}
	}
	if user.Email == logReq.User {
		if checkPasswordHash(logReq.Password, user.PasswordHash) {
			return user, nil
		} else {
			return mysql.User{}, nil
		}
	}
	return user, err
}

func validateRegistrationRequest(regReq registrationRequest) (bool, error) {
	isValid, err := IsEmailValid(regReq.Email)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	if !IsUsernameValid(regReq.Username) {
		return false, nil
	}
	if !IsPasswordValid(regReq.Password) {
		return false, nil
	}
	if !IsNameValid(regReq.FirstName, regReq.LastName) {
		return false, nil
	}
	isValid, err = IsBirthDateValid(regReq.BirthDate)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}
	return true, nil
}

func isUserUnique(db *mysql.DB, regReq registrationRequest) (bool, error) {
	user, err := db.FindUserByUsername(regReq.Username)
	var isUnique = false
	if err != nil {
		isUnique = true
	}
	if user.Username == regReq.Username {
		return false, errors.New("username taken")
	}
	user, err = db.FindUserByEmail(regReq.Email)
	if err != nil {
		isUnique = true
	}
	if user.Email == regReq.Email {
		return false, errors.New("email taken")
	}
	return isUnique, nil
}

func IsEmailValid(email string) (bool, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func IsUsernameValid(username string) bool {
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

func IsPasswordValid(password string) bool {
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

func IsNameValid(firstName string, lastName string) bool {
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

func IsBirthDateValid(birthDate string) (bool, error) {
	birthDate2, err := time.Parse("02-01-2006", birthDate)
	if err != nil {
		log.Println("Error parsing time")
		return false, err
	}
	currentTime := time.Now().AddDate(-18, 0, 0)
	if birthDate2.Before(currentTime) {
		return true, nil
	}
	return false, nil
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
