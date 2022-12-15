package mysql

import (
	"bettingAPI/internal/source"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DB struct {
	conn *sql.DB
}

func NewDB() *DB {
	return &DB{
		conn: ConnectDB(),
	}
}

func ConnectDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Lozinka123#@tcp(localhost:3306)/bettingdb")
	if err != nil {
		log.Fatal()
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connectiong to DB: %s", err)
	} else {
		fmt.Println("Connected to DB")
	}
	return db
}

type League struct {
	Name         string
	Elaborations []Elaboration
}

type Elaboration struct {
	Tips []Tip
	ID   []int
}

type Tip struct {
	Name  string
	Value float64
}

type Offer struct {
	Number        string
	ID            int
	Name          string
	Time          string
	Tips          []Tip
	TvChannel     string
	HasStatistics bool
}

type LeagueOffers struct {
	ID     int    `db:"league_id" json:"id"`
	Title  string `db:"title" json:"title"`
	Offers []int  `json:"offers"`
}

type OfferByID struct {
	Name          string `json:"game"`
	Time          string `json:"time"`
	TvChannel     string `json:"tv_channel"`
	HasStatistics bool   `json:"statistics"`
	Tips          []Tip  `json:"tips"`
}

type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username" `
	Email        string    `db:"email" `
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	BirthDate    time.Time `db:"birth_date"`
	Balance      float32   `db:"balance"`
}

type OfferTip struct {
	OfferID     int     `db:"id"`
	Tip         string  `db:"tip"`
	Coefficient float64 `db:"coefficient"`
}

type BetSlipRequest struct {
	Username string  `json:"username" `
	Stake    float32 `json:"stake"`
	Bets     []Bet   `json:"bets"`
}

type Bet struct {
	OfferID int    `json:"offer"`
	Tip     string `json:"tip"`
}

type UserBetSlip struct {
	ID          int     `db:"id"`
	UserID      int     `db:"user_id"`
	Stake       float32 `db:"stake"`
	Coefficient float32 `db:"Coefficient"`
	Payout      float32 `db:"payout"`
}

func (d *DB) InsertOffers(offers []source.Offer) {
	query := "INSERT INTO `bettingdb`.`offers`(id, game, time_played, tv_channel, has_statistics) SELECT ?,?,?,?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`offers` WHERE id=?)"

	for i := range offers {
		_, err := d.conn.Exec(query, offers[i].ID, offers[i].Name, offers[i].Time, offers[i].TvChannel, offers[i].HasStatistics, offers[i].ID)
		if err != nil {
			log.Printf("impossible to insert offers: %s", err)
		}
	}
}

func (d *DB) InsertLeagues(leagues *source.LeaguesData) {
	query := "INSERT INTO `bettingdb`.`leagues`(title) SELECT ? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`leagues` WHERE `bettingdb`.`leagues`.`title`=?)"

	for i := range leagues.Leagues {
		res, err := d.conn.Exec(query, leagues.Leagues[i].Name, leagues.Leagues[i].Name)
		if err != nil {
			log.Printf("Impossible to insert leagues: %s", err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			log.Printf("Impossible to insert leagues: %s", err)
		}
		if ra == 0 {
			continue
		}
		q := "INSERT INTO `bettingdb`.`league_offers`(league_id, offer_id) VALUES(?,?)"
		league_id, err := res.LastInsertId()
		if err != nil {
			log.Printf("Impossible to get last insert id: %s", err)
		}
		for j := range leagues.Leagues[i].Elaborations {
			for k := range leagues.Leagues[i].Elaborations[j].ID {
				_, err := d.conn.Exec(q, league_id, leagues.Leagues[i].Elaborations[j].ID[k])
				if err != nil {
					log.Printf("Impossible to insert league offer: %s", err)
				}
			}
		}
	}
}

func (d *DB) InsertTips(offers []source.Offer) {
	query := "INSERT INTO `bettingdb`.`offer_tips`(offer_id, tip, coefficient) SELECT ?,?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`offer_tips` WHERE `bettingdb`.`offer_tips`.`offer_id`=? AND `bettingdb`.`offer_tips`.`tip`=?)"

	for i := range offers {
		for j := range offers[i].Tips {
			_, err := d.conn.Exec(query, offers[i].ID, offers[i].Tips[j].Name, offers[i].Tips[j].Value, offers[i].ID, offers[i].Tips[j].Name)
			if err != nil {
				log.Printf("Impossible to insert tips: %s", err)
			}
		}
	}
}

func (d *DB) GetLeagueOffers() []LeagueOffers {
	rowsLeagues, err := d.conn.Query("SELECT * FROM `bettingdb`.`leagues`")
	if err != nil {
		log.Printf("Impossible to scan from leagues table: %s", err)
	}
	defer rowsLeagues.Close()
	var leagueOffers []LeagueOffers
	for rowsLeagues.Next() {
		var league LeagueOffers
		if err := rowsLeagues.Scan(&league.ID, &league.Title); err != nil {
			log.Printf("Impossible to scan from leagues table: %s", err)
		}
		leagueOffers = append(leagueOffers, league)
	}
	for i := range leagueOffers {
		rowsOffers, err := d.conn.Query("SELECT offer_id FROM `bettingdb`.`league_offers` WHERE `bettingdb`.`league_offers`.`league_id`=?", leagueOffers[i].ID)
		if err != nil {
			log.Printf("Impossible to select from league_offers table: %s", err)
		}
		defer rowsOffers.Close()
		for rowsOffers.Next() {
			var Offer int
			if err := rowsOffers.Scan(&Offer); err != nil {
				log.Printf("Impossible to scan from league_offers table: %s", err)
			}
			leagueOffers[i].Offers = append(leagueOffers[i].Offers, Offer)
		}
	}
	return leagueOffers
}

func (d *DB) GetOfferByID(offerID int) interface{} {
	var offer OfferByID
	row := d.conn.QueryRow("SELECT game, time_played, tv_channel, has_statistics from `bettingdb`.`offers` where `bettingdb`.`offers`.`id`=?", offerID)
	err := row.Scan(&offer.Name, &offer.Time, &offer.TvChannel, &offer.HasStatistics)
	if err != nil {
		log.Printf("Error getting offer by id: %s", err)
	}
	rowsTips, err := d.conn.Query("SELECT tip, coefficient FROM `bettingdb`.`offer_tips` WHERE `bettingdb`.`offer_tips`.`offer_id`=?", offerID)
	if err != nil {
		log.Printf("Impossible to select from offer_tips table: %s", err)
	}
	defer rowsTips.Close()
	for rowsTips.Next() {
		var tip Tip
		if err := rowsTips.Scan(&tip.Name, &tip.Value); err != nil {
			log.Printf("Impossible to scan from offer_tips table: %s", err)
		}
		offer.Tips = append(offer.Tips, tip)
	}
	return offer
}

func (d *DB) FindUserByEmail(email string) (*User, error) {
	row := d.conn.QueryRow("SELECT username, email, first_name, last_name, password_hash, balance FROM `bettingdb`.`users` WHERE `bettingdb`.`users`.`email`=?", email)
	var user User
	err := row.Scan(&user.Username, &user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) FindUserByUsername(username string) (*User, error) { // moze vratit 3 slucaja
	row := d.conn.QueryRow("SELECT id,username, email, first_name, last_name, password_hash, balance FROM `bettingdb`.`users` WHERE `bettingdb`.`users`.`username`=?", username)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.PasswordHash, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) InsertUser(regReq User) error {
	var user = regReq
	query := "INSERT INTO `bettingdb`.`users`(username, email, password_hash, first_name, last_name, birth_date, balance) VALUES(?,?,?,?,?,?,?)"
	_, err := d.conn.Exec(query, user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.BirthDate, 100)
	if err != nil {
		log.Printf("impossible to insert user: %s", err)
		return err
	}
	return nil
}

func (d *DB) GetOfferTipCoefficients(bets []Bet) ([]OfferTip, error) {
	var offerTips []OfferTip
	for i := range bets {
		query := "SELECT offer_id, tip, coefficient FROM `bettingdb`.`offer_tips` WHERE `bettingdb`.`offer_tips`.`offer_id`=? AND `bettingdb`.`offer_tips`.`tip`=?"
		rows, err := d.conn.Query(query, bets[i].OfferID, bets[i].Tip)
		if err != nil {
			log.Printf("Impossible to select from offer_tips table for coefficient: %s", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var ot OfferTip
			if err := rows.Scan(&ot.OfferID, &ot.Tip, &ot.Coefficient); err != nil {
				log.Printf("Impossible to scan from offer_tips table for coefficient: %s", err)
			}
			offerTips = append(offerTips, ot)
		}
	}
	return offerTips, nil
}

func (d *DB) UpdateUserBalance(user User, updatedUserBalance float32) error {
	query := "UPDATE `bettingdb`.`users` SET `bettingdb`.`users`.`balance`=? WHERE `bettingdb`.`users`.`id`=?"
	_, err := d.conn.Exec(query, updatedUserBalance, user.ID)
	if err != nil {
		log.Printf("Impossible to update user balance: %s", err)
		return err
	}
	return nil
}

func (d *DB) InsertUserBetSlip(userBetSlip UserBetSlip, betSlip BetSlipRequest) error { // prvo insert listica pa onda insertaj pripadajuce betove na temelju lastinsertid
	query1 := "INSERT INTO `bettingdb`.`user_bet_slips`(user_id, stake, coefficient, payout) VALUES(?,?,?,?)"
	res, err := d.conn.Exec(query1, userBetSlip.UserID, userBetSlip.Stake, userBetSlip.Coefficient, userBetSlip.Payout)
	if err != nil {
		log.Printf("impossible to insert user: %s", err)
		return err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		log.Printf("Impossible to insert leagues: %s", err)
		return err
	}
	if ra == 0 {
		return nil
	}

	offerTips, err := d.GetOfferTipCoefficients(betSlip.Bets)
	if err != nil {
		return err
	}
	user_bet_slip_id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Impossible to get last insert id: %s", err)
		return err
	}
	for i := range offerTips {
		q := "INSERT INTO `bettingdb`.`bet`(user_bet_slip_id, offer_id, tip, coefficient) VALUES(?,?,?,?)"

		_, err = d.conn.Exec(q, user_bet_slip_id, offerTips[i].OfferID, offerTips[i].Tip, offerTips[i].Coefficient)
		if err != nil {
			log.Printf("Impossible to insert league offer: %s", err)
			return err
		}
	}
	return nil
}
