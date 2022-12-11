package http

import (
	"bettingAPI/internal/mysql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type handler struct {
	db *mysql.DB
}

func NewHandler(d *mysql.DB) *handler {
	return &handler{
		db: d,
	}
}

type registrationRequest struct {
	Username  string `json:"username" `
	Email     string `json:"email" `
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
}

type loginRequest struct {
	User     string `json:"user" `
	Password string `json:"password"`
}

func (h *handler) GetLeagueOffers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var leagueOffers = h.db.GetLeagueOffers()
	json.NewEncoder(w).Encode(leagueOffers)
}

func (h *handler) GetOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	offerID, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Printf("Error converting id from string to int: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offer := h.db.GetOfferByID(offerID)
	json.NewEncoder(w).Encode(offer)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from register form: %s", err)
	}
	var regReq registrationRequest
	err = json.Unmarshal(body, &regReq)
	if err != nil {
		log.Printf("Error unmarshalling data from register form: %s", err)
	}
	var isValid = true
	var user mysql.User
	err = regReq.validateEmail()
	if err != nil {
		isValid = false
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		err = h.db.IsEmailUsed(regReq.Email)
		if err != nil {
			isValid = false
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			user.Email = regReq.Email
		}
	}

	err = regReq.validateUsername()
	if err != nil {
		isValid = false
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		err = h.db.IsUsernameUsed(regReq.Username)
		if err != nil {
			isValid = false
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			user.Username = regReq.Username
		}
	}

	err = regReq.validatePassword()
	if err != nil {
		isValid = false
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		user.PasswordHash = string(hashPassword(regReq.Password))
	}
	err = regReq.validateName()
	if err != nil {
		isValid = false
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		user.FirstName = regReq.FirstName
		user.LastName = regReq.LastName
	}

	birthDate, err := regReq.validateBirthDate()
	if err != nil {
		isValid = false
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		user.BirthDate = birthDate
	}
	if isValid {
		h.db.InsertUser(user)
		log.Printf("user inserted")
	} else {
		log.Printf("user not inserted")
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from login form: %s", err)
	}
	var logReq loginRequest
	err = json.Unmarshal(body, &logReq)
	if err != nil {
		log.Printf("Error unmarshalling data from login form: %s", err)
	}
	err = h.checkUserAndPassword(logReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		log.Printf("login successful")
	}
}
