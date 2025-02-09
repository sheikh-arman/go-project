package handler

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Handle() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	log.Info().Msg("Starting Handle server")
	log.Error().Msg("This is not error actually")
	// test()
	testStandard()
}
