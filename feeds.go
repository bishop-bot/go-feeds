package main

import (
	"context"

	"github.com/bishop-bot/feeds/config"
	"github.com/bishop-bot/feeds/parser"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// feed := &parser.Feed{URL: "https://feeds.bloomberg.com/markets/news.rss"}
	ymlConfig := config.MustLoadConfig("config.yaml")
	log.Info().Msgf("Loaded config: %+v", ymlConfig)

	ctx := context.Background()
	_, err := parser.FetchAll(ctx, &ymlConfig.Feeds)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch feeds")
		return
	}
}
