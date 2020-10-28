package orchestrator

import (
	"github.com/hennersz/ConfSync/internal/sync"
	"github.com/hennersz/ConfSync/internal/updater"
)

func SyncAndUpdate(sourceRepo, workDir string) error {
	syncer := sync.NewGitSyncer(sourceRepo, workDir)

	shouldUpdate, err := syncer.Sync()
	if err != nil {
		return err
	}

	u, err := updater.NewUpdater(workDir)
	if err != nil {
		return err
	}

	if shouldUpdate {

		err = u.Update()

		if err != nil {
			return err
		}
	}

	return nil
}
