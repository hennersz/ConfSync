package sync

import (
	"errors"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

var ErrDestNotEmpty = errors.New("Destination directory is not empty and not a git repository")

// notEmpy checks if a directory is not empty
// A directory counts as empty if it is empty or it doesn't exist
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

func (s *GitSyncer) Sync() (bool, error) {
	ne, err := notEmpty(s.destPath)
	if err != nil {
		return false, err
	}

	if ne {
		r, err := git.PlainOpen(s.destPath)
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return false, ErrDestNotEmpty
		}

		w, err := r.Worktree()
		if err != nil {
			return false, err
		}

		err = w.Pull(&git.PullOptions{RemoteName: "origin"})

		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false, nil
		} else if err != nil {
			return false, err
		}
	} else {
		_, err = git.PlainClone(s.destPath, false, &git.CloneOptions{
			URL: s.sourceRepo,
		})
		if err != nil {
			return false, err
		}
	}
	return true, err
}
