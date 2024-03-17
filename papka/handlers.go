package main

import (
	"encoding/json"
	"net/http"
)

type Actor struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Birth  string `json:"birth"`
}

type Movie struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Release     string  `json:"release"`
	Rate        int     `json:"rate"`
	Actor_list  []Actor `json:"actor_list"`
}

// Получение списка актеров
func GetActors(w http.ResponseWriter, r *http.Request) {
	actors := make([]Actor, 0)
	rows, err := db.Query("SELECT id, name, gender, birth from actors")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var act Actor
		if err := rows.Scan(&act.ID, &act.Name, &act.Gender, &act.Birth); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		actors = append(actors, act)
	}
	json.NewEncoder(w).Encode(actors)
}

// Добавление информации об актере
func AddActor(w http.ResponseWriter, r *http.Request) {
	var act Actor

	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO actors (name, gender, birth) values ($1 $2 $3)",
		act.Name, act.Gender, act.Birth)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Изменение информации об актере
func UpdActor(w http.ResponseWriter, r *http.Request) {
	var act Actor

	if err := json.NewDecoder(r.Body).Decode(&act); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE actors SET name = $1, gender = $2, birth = $3 WHERE id = $4",
		act.Name, act.Gender, act.Birth, act.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление информации об актере
func DelActor(w http.ResponseWriter, r *http.Request) {
	actorID := r.URL.Query().Get("id")
	if actorID == "" {
		http.Error(w, "No actor ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM actors WHERE id = $1", actorID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Получение списка фильмов
func GetMovies(w http.ResponseWriter, r *http.Request) {
	sortField := r.URL.Query().Get("sort_field")
	sortOrder := r.URL.Query().Get("sort_order")

	if sortField == "" {
		sortField = "rate"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		http.Error(w, "Not correct sort order", http.StatusBadRequest)
		return
	}
	rows, err := db.Query(`SELECT m.id, m.name, m.description, m.release, m.rate, a.id, a.name, a.gender, a.birth
	FROM movies m
	LEFT JOIN movie_actors ma ON m.id = ma.movie_id
	LEFT JOIN actors a ON ma.actor_id = a.id
	ORDER BY $1 $2`, sortField, sortOrder)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	moviesMap := make(map[int]*Movie)

	for rows.Next() {
		var movieID, movieRate int
		var movieName, movieDescription, movieRelease string
		var actorID int
		var actorName, actorGender, actorBirth string

		err := rows.Scan(&movieID, &movieName, &movieDescription, &movieRelease, &movieRate,
			&actorID, &actorName, &actorGender, &actorBirth)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := moviesMap[movieID]; !ok {
			moviesMap[movieID] = &Movie{
				ID:          movieID,
				Name:        movieName,
				Description: movieDescription,
				Release:     movieRelease,
				Rate:        movieRate,
				Actor_list:  make([]Actor, 0),
			}
		}

		moviesMap[movieID].Actor_list = append(moviesMap[movieID].Actor_list, Actor{
			ID:     actorID,
			Name:   actorName,
			Gender: actorGender,
			Birth:  actorBirth,
		})
	}

	movies := make([]Movie, 0)

	for _, movie := range moviesMap {
		movies = append(movies, *movie)
	}

	json.NewEncoder(w).Encode(movies)
}

// Добавление ифнормации о фильме
func AddMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO movies (name, description, release, rate)VALUES ($1, $2, $3, $4)",
		movie.Name, movie.Description, movie.Release, movie.Rate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var movieID int

	err = db.QueryRow("SELECT id FROM movies WHERE name = $1", movie.Name).Scan(&movieID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, actor := range movie.Actor_list {
		_, err := db.Exec("INSERT INTO movie_actos (movie_id, actor_id) VALUES ($1, $2)",
			movieID, actor.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

// Изменение ифнормации о фильме
func UpdMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE movies SET name = $1, description = $2, release = $3, rate = $4 WHERE id = $5",
		movie.Name, movie.Description, movie.Release, movie.Rate, movie.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM movie_actors WHERE movie_id = $1", movie.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, actor := range movie.Actor_list {
		_, err := db.Exec("INSERT INTO movie_actors (movie_id, actor_id) VALUES ($1, $2)",
			movie.ID, actor.ID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// Удаление информации о фильме
func DelMovie(w http.ResponseWriter, r *http.Request) {
	movieID := r.URL.Query().Get("id")
	if movieID == "" {
		http.Error(w, "No movie ID", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM movies WHERE id = $1", movieID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
