package main

import (
	"bettingAPI/pkg/router"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	r := router.Router()
	log.Fatal(http.ListenAndServe("localhost:5000", r))

}
