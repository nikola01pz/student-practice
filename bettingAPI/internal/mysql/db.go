package mysql

import (
	"bettingAPI/internal/source"
	"database/sql"
	"fmt"
	"log"
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
	Name  string  `json:"name"`
	Value float64 `json:"value"`
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

func (d *DB) InsertOffers(offers []source.Offer) {
	query := "INSERT INTO `bettingdb`.`offers`(offer_id, game, time_played, tv_channel, has_statistics) SELECT ?,?,?,?,? WHERE NOT EXISTS(SELECT * FROM `bettingdb`.`offers` WHERE offer_id=?)"

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
	row := d.conn.QueryRow("SELECT game, time_played, tv_channel, has_statistics from `bettingdb`.`offers` where `bettingdb`.`offers`.`offer_id`=?", offerID)
	err := row.Scan(&offer.Name, &offer.Time, &offer.TvChannel, &offer.HasStatistics)
	if err != nil {
		log.Printf("Error getting offer by id: %s", err)
	}
	rowsTips, err := d.conn.Query("select tip, coefficient from `bettingdb`.`offer_tips` where `bettingdb`.`offer_tips`.`offer_id`=?", offerID)
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
