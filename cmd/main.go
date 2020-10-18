package main

import (
	"fmt"
	"os"

	"github.com/hennersz/ConfSync/internal/sync"
	"github.com/hennersz/ConfSync/internal/updater"
)

func main() {
	syncer := sync.NewGitSyncer("https://github.com/hennersz/ConfSync.git", "./test")
	err := syncer.Sync()

	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}

	u, err := updater.NewUpdater("./test")

	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}

	err = u.Update()

	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
