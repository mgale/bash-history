package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mgale/bash-history.git/internal/defaults"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var natsURL string
var natsDocumentSynerTopic string
var natsCommandTopic string
var typesenseAPIKey string
var typesenseURL string

var rootCmd = &cobra.Command{
	Use:   "document-syncer",
	Short: "document-syncer is a command line tool for loading and syncing history events to the search engine",
}

func Execute() int {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func(cancel context.CancelFunc) {
		<-sig
		cancel()
	}(cancel)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return 1
	}

	return 0
}

func init() {
	viper.SetDefault("NATSIO_URL", "nats://localhost:4222")
	viper.SetDefault("NATSIO_COMMAND_TOPIC", defaults.CommandTopic)
	viper.SetDefault("NATSIO_DOCUMENT_SYNCER_TOPIC", defaults.DocumentSyncTopic)
	viper.SetDefault("TYPESENSE_API_KEY", "typesense123")
	viper.SetDefault("TYPESENSE_URL", "http://localhost:8108")
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVar(&natsURL, "natsio-url", viper.GetString("NATSIO_URL"), "Default natsio URL, will load env NATSIO_URL if set")
	rootCmd.PersistentFlags().StringVar(&natsCommandTopic, "natsio-command-topic", viper.GetString("NATSIO_COMMAND_TOPIC"), "Default natsio command topic, will load env NATSIO_COMMAND_TOPIC if set")
	rootCmd.PersistentFlags().StringVar(&natsDocumentSynerTopic, "natsio-document-syncer-topic", viper.GetString("NATSIO_DOCUMENT_SYNCER_TOPIC"), "Default natsio document syncer topic, will load env NATSIO_DOCUMENT_SYNCER_TOPIC if set")
	rootCmd.PersistentFlags().StringVar(&typesenseAPIKey, "typesense-api-key", viper.GetString("TYPESENSE_API_KEY"), "Default typesense API key, will load env TYPESENSE_API_KEY if set")
	rootCmd.PersistentFlags().StringVar(&typesenseURL, "typesense-url", viper.GetString("TYPESENSE_URL"), "Default typesense URL, will load env TYPESENSE_URL if set")
}
