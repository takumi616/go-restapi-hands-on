package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/takumi616/go-restapi-hands-on/entity"
	"github.com/takumi616/go-restapi-hands-on/store"
)

type AddTask struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

// Http handler to add a new task
func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}

	//Convert json http request data into go struct type
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//Validate request body
	err := validator.New().Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	//Set request body to Task entity type
	t := &entity.Task{
		Title:  b.Title,
		Status: entity.TaskStatusTodo,
	}

	//Store request data to db
	err = at.Repo.AddTask(ctx, at.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	//Create json response body writing newly created record's id
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
