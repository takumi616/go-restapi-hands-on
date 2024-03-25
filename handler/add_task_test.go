package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/takumi616/go-restapi-hands-on/entity"
	"github.com/takumi616/go-restapi-hands-on/testutil"
)

func TestAddTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}

	//Prepare two test case (ok and bad request pattern)
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		//Execute as parallel tests
		//Run runs function as a subtest of t called name n(first parameter of Run).
		//It runs function in a separate goroutine and blocks
		//until this function returns or calls t.Parallel to become a parallel test.
		t.Run(n, func(t *testing.T) {
			//Parallel signals that this test is to be run in parallel
			//with (and only with) other parallel tests.
			t.Parallel()

			//Create test http request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(
				ctx context.Context, title string,
			) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			//Send http request
			sut := AddTask{
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			//Check http response body
			resp := w.Result()
			testutil.AssertResponse(t,
				resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
