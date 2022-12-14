package http

import (
	"bettingAPI/internal/mysql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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
	Username  string  `json:"username" `
	Email     string  `json:"email" `
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	BirthDate string  `json:"birth_date"`
	Balance   float32 `json:"balance"`
}

type loginInfo struct {
	Username  string  `json:"username" `
	Email     string  `json:"email" `
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float32 `json:"balance"`
}

type loginRequest struct {
	User     string `json:"user" `
	Password string `json:"password"`
}

type betSlipRequest struct {
	Username string  `json:"username" `
	Stake    float32 `json:"stake"`
	Bets     []bet   `json:"elaboration"`
}

type bet struct {
	OfferID int    `json:"offer"`
	Tip     string `json:"tip"`
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

	isRegReqValid, err := validateRegistrationRequest(regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	isUserUnique, err := isUserUnique(h.db, regReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	if isRegReqValid && isUserUnique {
		passwordHash, err := hashPassword(regReq.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		birthDate, err := time.Parse("02-01-2006", regReq.BirthDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		user := mysql.User{
			Email:        regReq.Email,
			Username:     regReq.Username,
			PasswordHash: passwordHash,
			FirstName:    regReq.FirstName,
			LastName:     regReq.LastName,
			BirthDate:    birthDate,
			Balance:      regReq.Balance,
		}
		err = h.db.InsertUser(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else {
		w.WriteHeader(http.StatusInternalServerError)
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

	isValid, err := validateLoginRequest(logReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if isValid {
		user, err := h.validateUserAndPassword(logReq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		} else if (user != mysql.User{}) {
			userInfo := loginInfo{
				Email:     user.Email,
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Balance:   user.Balance,
			}
			json.NewEncoder(w).Encode(userInfo)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h *handler) BetSlip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from login form: %s", err)
	}
	var betSlip betSlipRequest
	err = json.Unmarshal(body, &betSlip)
	if err != nil {
		log.Printf("Error unmarshalling data from bet slip request: %s", err)
	}
	// validate

	//insert

}
