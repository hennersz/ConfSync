package updater_test

import (
	"testing"

	"github.com/hennersz/ConfSync/internal/updater"
)

func TestLoadConfig(t *testing.T) {
	u, err := updater.NewUpdater("./testdata/")

	if err != nil {
		t.Errorf("Unexpected error occured: %v", err)
	}
	if u.Config["test.txt"] != "/some/place/to/go" {
		t.Error("Config not loaded correctly")
	}
}
