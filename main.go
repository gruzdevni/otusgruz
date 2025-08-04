package main

import (
	"context"
	"os"

	"otusgruz/cmd"
	"otusgruz/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	exitCode := 0

	log.Ctx(ctx).Info().Msg("the application is launching")

	err = cmd.Run(ctx, conf)
	if err != nil {
		exitCode = 1
	}

	os.Exit(exitCode)
}
