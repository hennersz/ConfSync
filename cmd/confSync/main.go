package main

import (
	"fmt"

	"github.com/hennersz/ConfSync/internal/orchestrator"
)

func main() {
	err := orchestrator.SyncAndUpdate("https://github.com/hennersz/ConfSyncTestRepo.git", "./source")

	if err != nil {
		fmt.Printf("An error occured: %v", err)
	}
}
