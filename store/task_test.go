package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/takumi616/go-restapi-hands-on/clock"
	"github.com/takumi616/go-restapi-hands-on/entity"
	"github.com/takumi616/go-restapi-hands-on/testutil"
)

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	//ListTasks() fetches all records from DB, so the result of this function test will be chenged
	//if some records are added by other test cases.
	//To prevent this, creating transaction will be necessary.
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)

	// Execute rollback when finishing this test case
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	//Add some records to test DB
	wants := prepareTasks(ctx, t, tx)

	//Check if ListTasks() works correctly
	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	//Compare result data to expected one
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()

	//Clean up records
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("Failed to initialize task table: %v", err)
	}

	//Add some records before testing ListTasks()
	//Use FixedClocker to compare wants to got
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}

	//LastInsertedId is the id which is inserted first
	//when inserting records at the same time
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	//Create test data inserted into DB
	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	//Creates sqlmock database connection and a mock to manage expectations.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	//Set expectation to mock
	mock.ExpectExec(
		"INSERT INTO task",
	).WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	//Execute AddTask method
	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
