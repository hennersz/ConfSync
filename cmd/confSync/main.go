package main

import (
	"os"

	"github.com/hennersz/ConfSync/internal/orchestrator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Starting")
	err := orchestrator.SyncAndUpdate("https://github.com/hennersz/ConfSyncTestRepo.git", "./source")

	if err != nil {
		log.Err(err).Send()
		os.Exit(1)
	}
}
