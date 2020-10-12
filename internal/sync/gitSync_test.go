package sync

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type testFile struct {
	name string
	data []byte
}

var testFiles = []testFile{
	{
		"test.txt",
		[]byte("hello"),
	},
}

func addFiles(t *testing.T, repo *git.Repository, repoDir string, files []testFile) error {
	t.Helper()
	for _, file := range files {
		fileName := filepath.Join(repoDir, file.name)
		if err := ioutil.WriteFile(fileName, file.data, 0644); err != nil {
			return err
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err := w.AddGlob("."); err != nil {
		return err
	}

	_, err = w.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func createRepo(t *testing.T, files []testFile) (*git.Repository, string) {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error creating tempdir: %v", err)
	}

	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Could not init repo: %v", err)
	}

	if len(files) > 0 {
		err := addFiles(t, repo, dir, files)

		if err != nil {
			t.Fatalf("An error occured while adding files to repo: %v", err)
		}
	}
	return repo, dir
}

func TestClone(t *testing.T) {
	_, sourceRepo := createRepo(t, testFiles)

	destRepo, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error creating tempdir: %v", err)
	}

	syncer := NewGitSyncer(sourceRepo, destRepo)
	err = syncer.Sync()
	if err != nil {
		t.Fatalf("Error syncing repo: %v", err)
	}

	fileName := filepath.Join(destRepo, testFiles[0].name)
	res, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("expected file not in cloned directory, %v", err)
	}
	if !bytes.Equal(res, testFiles[0].data) {
		t.Errorf("Expected %v, got %v", testFiles[0].data, res)
	}
}

func TestCloneFailIntoNonEmpty(t *testing.T) {
	_, sourceRepo := createRepo(t, testFiles)

	destRepo, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error creating tempdir: %v", err)
	}

	ioutil.WriteFile(filepath.Join(destRepo, "aFile.txt"), []byte("data"), 0644)

	syncer := NewGitSyncer(sourceRepo, destRepo)
	err = syncer.Sync()
	if err != ErrDestNotEmpty {
		t.Errorf("expected %v, got %v", ErrDestNotEmpty, err)
	}
}

func TestCloneOkToNonExistentDir(t *testing.T) {
	_, sourceRepo := createRepo(t, testFiles)

	destRepo, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error creating tempdir: %v", err)
	}

	syncer := NewGitSyncer(sourceRepo, filepath.Join(destRepo, "subdir"))
	err = syncer.Sync()
	if err != nil {
		t.Errorf("Unexpected error occured: %v", err)
	}
}

func TestUpdateOnSecondRun(t *testing.T) {
	sourceRepo, sourceDir := createRepo(t, testFiles)

	destRepo, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error creating tempdir: %v", err)
	}

	syncer := NewGitSyncer(sourceDir, destRepo)
	err = syncer.Sync()
	if err != nil {
		t.Fatalf("An unexpected error occurred: %v", err)
	}

	extraFiles := []testFile{
		{
			"test2.txt",
			[]byte("hello"),
		},
	}

	err = addFiles(t, sourceRepo, sourceDir, extraFiles)
	if err != nil {
		t.Fatalf("An error occured while adding files to repo: %v", err)
	}

	err = syncer.Sync()
	if err != nil {
		t.Fatalf("An unexpected error occurred: %v", err)
	}

	fileName := filepath.Join(destRepo, extraFiles[0].name)
	res, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("expected file not in cloned directory, %v", err)
	}
	if !bytes.Equal(res, extraFiles[0].data) {
		t.Errorf("Expected %v, got %v", testFiles[0].data, res)
	}
}
