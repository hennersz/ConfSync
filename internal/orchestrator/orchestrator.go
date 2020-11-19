package orchestrator

import (
	"github.com/hennersz/ConfSync/internal/sync"
	"github.com/hennersz/ConfSync/internal/updater"
	"github.com/pkg/errors"
)

func SyncAndUpdate(sourceRepo, workDir string) error {
	syncer := sync.NewGitSyncer(sourceRepo, workDir)

	shouldUpdate, err := syncer.Sync()
	if err != nil {
		return err
	}

	if shouldUpdate {
		u, err := updater.New().SrcDir(workDir).Build()
		if err != nil {
			return errors.Wrap(err, "An error occured reading the config file")
		}

		err = u.Update()

		if err != nil {
			return errors.Wrap(err, "An error occured while updating config")
		}
	}

	return nil
}
