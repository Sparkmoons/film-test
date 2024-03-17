package main

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user = user dbname = db_name sslmode = disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("actors", GetActors)
	http.HandleFunc("/actors/add", AddActor)
	http.HandleFunc("/actors/update", UpdActor)
	http.HandleFunc("actors/delete", DelActor)

	http.HandleFunc("/movies", GetMovies)
	http.HandleFunc("/movies/add", AddMovie)
	http.HandleFunc("/movies/update", UpdMovie)
	http.HandleFunc("/movies/delete", DelMovie)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

