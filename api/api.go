package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

var db = make(map[string]User)

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", createUser(db))
			r.Get("/", getAllUsers(db))
			r.Get("/{id}", getSpecificUser(db))
			r.Delete("/{id}", deleteUser(db))
			r.Put("/{id}", updateUser(db))
		})
	})

	return r
}

func getAllUsers(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, Response{Data: db}, http.StatusOK)
	}
}

func createUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body User

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Error: "Invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		id := uuid.New().String()

		var user User = User{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Biography: body.Biography,
		}

		db[id] = user

		sendJSON(w, Response{Data: id}, http.StatusCreated)

	}
}

func getSpecificUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		data, ok := db[id]

		if !ok {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: data}, http.StatusOK)
	}
}

func deleteUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		data, ok := db[id]

		if !ok {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}

		delete(db, id)

		sendJSON(w, Response{Data: data}, http.StatusOK)
	}
}

func updateUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		_, ok := db[id]

		if !ok {
			sendJSON(w, Response{Error: "user not found"}, http.StatusNotFound)
			return
		}

		var body User

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Error: "Invalid body"}, http.StatusUnprocessableEntity)
			return
		}

		db[id] = body

		sendJSON(w, Response{Data: body}, http.StatusOK)

	}
}
