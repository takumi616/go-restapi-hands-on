package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/takumi616/go-restapi-hands-on/clock"
	"github.com/takumi616/go-restapi-hands-on/config"
	"github.com/takumi616/go-restapi-hands-on/handler"
	"github.com/takumi616/go-restapi-hands-on/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	//An endpoint to check if http server is running correctly
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	//Create validator
	v := validator.New()

	//Get DB handle
	//cleanup is used to close *sql.DB
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	r := store.Repository{Clocker: clock.RealClocker{}}
	//Register http handler
	at := &handler.AddTask{DB: db, Repo: &r, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{DB: db, Repo: &r}
	mux.Get("/tasks", lt.ServeHTTP)

	return mux, cleanup, nil
}
