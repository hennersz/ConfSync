package sync

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

var ErrDestNotEmpty = errors.New("destination directory is not empty and not a git repository")

// notEmpy checks if a directory is not empty
// A directory counts as empty if it is empty or it doesn't exist.
func notEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("can't open directory: %s\nerror: %w", name, err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return false, nil
	}

	if err != nil {
		err = fmt.Errorf("can't read directory: %s is empty\nerror: %w", name, err)
	}

	return true, err
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
		return s.pull()
	}

	return s.clone()
}

func (s *GitSyncer) pull() (bool, error) {
	r, err := git.PlainOpen(s.destPath)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return false, ErrDestNotEmpty
	}

	w, err := r.Worktree()
	if err != nil {
		return false, fmt.Errorf("can't get worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{RemoteName: "origin"})

	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("can't pull changes: %w", err)
	}

	return true, nil
}

func (s *GitSyncer) clone() (bool, error) {
	_, err := git.PlainClone(s.destPath, false, &git.CloneOptions{
		URL: s.sourceRepo,
	})
	if err != nil {
		return false, fmt.Errorf("can't clone repo: %s\nerror: %w", s.sourceRepo, err)
	}

	return true, nil
}
