package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takumi616/go-restapi-hands-on/entity"
	"github.com/takumi616/go-restapi-hands-on/testutil"
)

func TestListTask(t *testing.T) {
	t.SkipNow()
	type want struct {
		status  int
		rspFile string
	}

	//Prepare two test case (ok and empty response data pattern)
	tests := map[string]struct {
		tasks map[entity.TaskID]*entity.Task
		want  want
	}{
		"ok": {
			tasks: map[entity.TaskID]*entity.Task{
				1: {
					ID:     1,
					Title:  "test1",
					Status: "todo",
				},
				2: {
					ID:     2,
					Title:  "test2",
					Status: "done",
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json.golden",
			},
		},
		"empty": {
			tasks: map[entity.TaskID]*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		//Execute as parallel tests
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			//Create test http request and response writer
			w := httptest.NewRecorder()
			//r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			//Send http request
			//sut := ListTask{Store: &store.TaskStore{Tasks: tt.tasks}}
			//sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
