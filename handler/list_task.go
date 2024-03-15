package handler

import (
	"net/http"

	"github.com/takumi616/go-restapi-hands-on/entity"
	"github.com/takumi616/go-restapi-hands-on/store"
)

type ListTask struct {
	Store *store.TaskStore
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

// Http handler to get all tasks
func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	//Get all tasks from db
	tasks := lt.Store.All()

	//Create response data
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}

	RespondJSON(ctx, w, rsp, http.StatusOK)
}
