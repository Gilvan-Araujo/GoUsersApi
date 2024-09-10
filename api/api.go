package api

import (
	"encoding/json"
	"net/http"
	"sync"

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

func createUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		var body User

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Message: "Please provide FirstName LastName and bio for the user"}, http.StatusBadRequest)
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
		mu.Unlock()
	}
}

var mu sync.Mutex

func getAllUsers(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		sendJSON(w, Response{Data: db}, http.StatusOK)
		mu.Unlock()
	}
}

func getSpecificUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		id := chi.URLParam(r, "id")

		data, ok := db[id]

		if !ok {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		sendJSON(w, Response{Data: data}, http.StatusOK)
		mu.Unlock()
	}
}

func deleteUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		id := chi.URLParam(r, "id")

		data, ok := db[id]

		if !ok {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		delete(db, id)

		sendJSON(w, Response{Data: data}, http.StatusOK)
		mu.Unlock()
	}
}

func updateUser(db map[string]User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		id := chi.URLParam(r, "id")

		_, ok := db[id]

		if !ok {
			sendJSON(w, Response{Message: "The user with the specified ID does not exist"}, http.StatusNotFound)
			return
		}

		var body User

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(w, Response{Message: "Please provide name and bio for the user"}, http.StatusUnprocessableEntity)
			return
		}

		db[id] = body

		sendJSON(w, Response{Data: body}, http.StatusOK)
		mu.Unlock()
	}
}
