package main

import (
	"fmt"
	"os"

	"github.com/hennersz/ConfSync/internal/sync"
)

func main() {
	syncer := sync.NewGitSyncer("https://github.com/hennersz/ConfSync.git", "./test")
	err := syncer.Sync()

	if err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
