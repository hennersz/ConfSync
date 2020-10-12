package sync

import (
	"errors"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

var ErrDestNotEmpty = errors.New("Destination directory is not empty and not a git repository")

func notEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return false, nil
	}
	return true, err // Either not empty or error, suits both cases
}

type GitSyncer struct {
	sourceRepo string
	destPath   string
}

func NewGitSyncer(sourceRepo, destPath string) *GitSyncer {
	return &GitSyncer{
		sourceRepo,
		destPath,
	}
}

func (s *GitSyncer) Sync() error {
	ne, err := notEmpty(s.destPath)
	if err != nil {
		return err
	}

	if ne {
		r, err := git.PlainOpen(s.destPath)
		if err == git.ErrRepositoryNotExists {
			return ErrDestNotEmpty
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			return err
		}

	} else {
		_, err = git.PlainClone(s.destPath, false, &git.CloneOptions{
			URL: s.sourceRepo,
		})
	}
	return err
}
