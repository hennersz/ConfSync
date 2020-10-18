package updater_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hennersz/ConfSync/internal/updater"
)

type testFile struct {
	name    string
	data    []byte
	destDir string
}

var testFiles = []testFile{
	{
		"test.txt",
		[]byte("hello"),
		"test/dir",
	},
}

func initTestDir(t *testing.T, files []testFile) (string, error) {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	conf := make(map[string]string)

	sourceDir := path.Join(dir, "source")
	err = os.MkdirAll(sourceDir, 0755)
	if err != nil {
		return "", err
	}

	destDir := path.Join(dir, "destination")
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filePath := path.Join(destDir, file.destDir)
		conf[file.name] = filePath
		err = ioutil.WriteFile(path.Join(sourceDir, file.name), file.data, 0644)
		if err != nil {
			return "", err
		}
	}
	confJson, err := json.Marshal(conf)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path.Join(dir, "source", "config.json"), confJson, 0644)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func TestWriteFiles(t *testing.T) {
	rootDir, err := initTestDir(t, testFiles)
	if err != nil {
		t.Fatalf("Unexpected error occured: %v", err)
	}

	u, err := updater.NewUpdater(path.Join(rootDir, "source"))
	if err != nil {
		t.Fatalf("Unexpected error occured: %v", err)
	}

	err = u.Update()

	if err != nil {
		t.Fatalf("Unexpected error occured: %v", err)
	}

	data, err := ioutil.ReadFile(path.Join(rootDir, "destination", "test/dir/test.txt"))
	if err != nil {
		t.Fatalf("An error occured while reading test file: %v", err)
	}

	if !bytes.Equal(testFiles[0].data, data) {
		t.Errorf("Expected %v, got %v", testFiles[0].data, data)
	}
}

func TestNoConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error occured creating temporary directory: %v", err)
	}

	_, err = updater.NewUpdater(dir)

	if err == nil {
		t.Error("Expected an error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected %v, got %v", os.ErrNotExist, err)
	}
}

func TestNonJsonConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Error occured creating temporary directory: %v", err)
	}

	err = ioutil.WriteFile(path.Join(dir, "config.json"), []byte("notJson"), 0644)
	if err != nil {
		t.Fatalf("Error occured while creating config file: %v", err)
	}

	_, err = updater.NewUpdater(dir)

	if err == nil {
		t.Error("Expected an error")
	}
}
