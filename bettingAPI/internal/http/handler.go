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
	Username  string `json:"username" `
	Email     string `json:"email" `
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
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

type BetSlipRequest struct {
	Username string  `json:"username" `
	Stake    float32 `json:"stake"`
	Bets     []bet   `json:"bets"`
}

type bet struct {
	OfferID int    `json:"offer"`
	Tip     string `json:"tip"`
}

type myError struct {
	Message string `json:"message"`
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

func (h *handler) HandleRegisterRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from register form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var regReq registrationRequest
	err = json.Unmarshal(body, &regReq)
	if err != nil {
		log.Printf("Error unmarshalling data from register form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	isRegReqValid := validateRegistrationRequest(regReq)
	if !isRegReqValid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isEmailUnique, err := h.isEmailUnique(regReq.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isEmailUnique {
		w.WriteHeader(http.StatusBadRequest)
		e := myError{
			Message: "email is already used",
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	isUsernameUnique, err := h.isUsernameUnique(regReq.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isUsernameUnique {
		w.WriteHeader(http.StatusBadRequest)
		e := myError{
			Message: "username is already used",
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	passwordHash, err := hashPassword(regReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	parsedBirthDate, err := time.Parse("02-01-2006", regReq.BirthDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := mysql.User{
		Email:        regReq.Email,
		Username:     regReq.Username,
		PasswordHash: passwordHash,
		FirstName:    regReq.FirstName,
		LastName:     regReq.LastName,
		BirthDate:    parsedBirthDate,
	}
	err = h.db.InsertUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) HandleLoginRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // ponavlja se pa izvuc u posebnu funkciju
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from login form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var logReq loginRequest
	err = json.Unmarshal(body, &logReq)
	if err != nil {
		log.Printf("Error unmarshalling data from login form: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isLoginRequestValid(logReq) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.Login(logReq.User, logReq.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userInfo := loginInfo{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Balance:   user.Balance,
	}
	json.NewEncoder(w).Encode(userInfo)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) HandleBetSlip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error getting data from login form: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var betSlip mysql.BetSlipRequest
	err = json.Unmarshal(body, &betSlip)
	if err != nil {
		log.Printf("Error unmarshalling data from bet slip request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isStakeSufficient(int(betSlip.Stake)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.db.FindBetUserByUsername(betSlip.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !hasSufficientFunds(int(user.Balance)) || !hasOnlyOneTipPerOffer(betSlip.Bets) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payout, isReached := h.isPayoutLimitReached(betSlip)
	if isReached {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userBetSlip := mysql.UserBetSlip{
		UserID: user.ID,
		Stake:  betSlip.Stake,
		Payout: payout,
	}

	err = h.db.InsertUserBetSlip(userBetSlip, betSlip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.db.UpdateUserBalance(userBetSlip.UserID, betSlip.Stake)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
