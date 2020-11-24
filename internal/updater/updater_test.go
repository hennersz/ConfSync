package updater_test

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hennersz/ConfSync/internal/updater"
)

type task struct {
	fileName string
	args     []string
}

type testOperator struct {
	name           string
	submittedTasks []task
	runCalls       []time.Time
}

func (to *testOperator) SubmitTask(fileName string, args []string) error {
	to.submittedTasks = append(to.submittedTasks, task{fileName, args})

	return nil
}

func (to *testOperator) Name() string {
	return to.name
}

func (to *testOperator) Run() error {
	to.runCalls = append(to.runCalls, time.Now())

	return nil
}

func TestNoConfig(t *testing.T) {
	t.Parallel()

	_, err := updater.New().SrcDir("testdata/empty").Build()

	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected %v, got %v", os.ErrNotExist, err)
	}
}

func TestNonJsonConfig(t *testing.T) {
	t.Parallel()

	_, err := updater.New().SrcDir("testdata/badJson").Build()

	var jsonError *json.SyntaxError
	if !errors.As(err, &jsonError) {
		t.Errorf("Expected a syntax error, got %T", err)
	}
}

func TestSubmitTask(t *testing.T) {
	t.Parallel()

	to := &testOperator{name: "test"}

	u, err := updater.New().SrcDir("testdata/simple").WithOperator(to).Build()
	if err != nil {
		t.Fatalf("An unexpected error occurred: %v", err)
	}

	err = u.Update()
	if err != nil {
		t.Fatalf("An unexpected error occurred: %v", err)
	}

	if len(to.submittedTasks) != 1 {
		t.Errorf("Incorrect number of tasks submitted, expected %v, got %v", 1, len(to.submittedTasks))
	}

	task := to.submittedTasks[0]

	if task.fileName != "test1.txt" {
		t.Errorf("Incorrect filename submitted, expected %v, got %v", "test1.txt", task.fileName)
	}

	if !reflect.DeepEqual(task.args, []string{"arg1"}) {
		t.Errorf("Incorrect args submitted, expected %v, got %v", []string{"arg1"}, task.args)
	}

	if len(to.runCalls) != 1 {
		t.Errorf("Operator run function called incorrect amount of times, expected %v, got %v", 1, len(to.runCalls))
	}
}

func TestOperatorNotFound(t *testing.T) {
	t.Parallel()

	u, err := updater.New().SrcDir("testdata/simple").Build()
	if err != nil {
		t.Fatalf("An unexpected error occurred: %v", err)
	}

	err = u.Update()

	if !errors.Is(err, updater.OperatorNotFoundError{OperatorName: "test"}) {
		t.Errorf("Expected %v, got %v", updater.OperatorNotFoundError{OperatorName: "test"}, err)
	}
}
