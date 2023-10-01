package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.mkabdelrahman.net/internal/data"
	"greenlight.mkabdelrahman.net/internal/jsonparser"
	"greenlight.mkabdelrahman.net/internal/validator"
)

type movieDataFromUser struct {
	Title   *string       `json:"title"`
	Year    *int32        `json:"year"`
	Runtime *data.Runtime `json:"runtime"`
	Genres  []string      `json:"genres"`
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// parse

	input, err := readMovieJSONfromRequest(w, r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	movie := adaptInputToMovieType(input)
	// validate
	if err := app.validateMovie(w, r, movie); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.respondWithCreatedMovie(w, r, movie)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonparser.WriteJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)

	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	input, err := readMovieJSONfromRequest(w, r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres // Note that we don't need to dereference a slice.
	}

	if err := app.validateMovie(w, r, movie); err != nil {
		app.failedValidationResponse(w, r, err)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = jsonparser.WriteJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = jsonparser.WriteJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func readMovieJSONfromRequest(w http.ResponseWriter, r *http.Request) (*movieDataFromUser, error) {
	var input movieDataFromUser
	err := jsonparser.ReadJSON(w, r, &input)
	return &input, err
}

func adaptInputToMovieType(input *movieDataFromUser) *data.Movie {
	return &data.Movie{
		Title:   *input.Title,
		Year:    *input.Year,
		Runtime: *input.Runtime,
		Genres:  input.Genres,
	}
}

func (app *application) validateMovie(w http.ResponseWriter, r *http.Request, movie *data.Movie) map[string]string {
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.IsValid() {
		return v.Errors

	}
	return nil
}

func (app *application) respondWithCreatedMovie(w http.ResponseWriter, r *http.Request, movie *data.Movie) {
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err := jsonparser.WriteJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
