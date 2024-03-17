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

	if err := createTables(); err != nil {
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

func createTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS actors (
			id SERIAL PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			gender VARCHAR(10),
			birth DATE
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			name VARCHAR(150) NOT NULL,
			description TEXT,
			release DATE,
			rate INT CHECK (rate >= 0 AND rate <= 10)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movie_actors (
			movie_id INT REFERENCES movies(id),
			actor_id INT REFERENCES actor(id),
			PRIMARY KEY (movie_id, actor_id)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
