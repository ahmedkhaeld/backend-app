package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	// figure the id
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	// call the database method that get the movie
	movie, err := app.models.DB.Get(id)
	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(w, err)
	}

}
func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	// call the database method that get list of movies
	movies, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteMovie(w http.ResponseWriter, r *http.Request) {}
func (app *application) insertMovie(w http.ResponseWriter, r *http.Request) {}
func (app *application) updateMovie(w http.ResponseWriter, r *http.Request) {}
func (app *application) searchMovie(w http.ResponseWriter, r *http.Request) {}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	// call the database to fetch list of all available genres
	genres, err := app.models.DB.AllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}
