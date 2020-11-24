package orchestrator

import (
	"fmt"

	"github.com/hennersz/ConfSync/internal/sync"
	"github.com/hennersz/ConfSync/internal/updater"
	"github.com/pkg/errors"
)

func SyncAndUpdate(sourceRepo, workDir string) error {
	syncer := sync.NewGitSyncer(sourceRepo, workDir)

	shouldUpdate, err := syncer.Sync()
	if err != nil {
		return fmt.Errorf("error occurred syncing: %w", err)
	}

	if shouldUpdate {
		u, err := updater.New().SrcDir(workDir).Build()
		if err != nil {
			return errors.Wrap(err, "An error occurred reading the config file")
		}

		err = u.Update()

		if err != nil {
			return errors.Wrap(err, "An error occurred while updating config")
		}
	}

	return nil
}
