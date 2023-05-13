package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Embed the entire templates directory
//go:embed templates/*
var myTemplates embed.FS

var typesenseAPIKey string
var typesenseURL string

var rootCmd = &cobra.Command{
	Use:   "history-searcher",
	Short: "history-searcher is a CLI to search through your shell history",
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
	viper.SetDefault("TYPESENSE_API_KEY", "typesense123")
	viper.SetDefault("TYPESENSE_URL", "http://localhost:8108")
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVar(&typesenseAPIKey, "typesense-api-key", viper.GetString("TYPESENSE_API_KEY"), "Default typesense API key, will load env TYPESENSE_API_KEY if set")
	rootCmd.PersistentFlags().StringVar(&typesenseURL, "typesense-url", viper.GetString("TYPESENSE_URL"), "Default typesense URL, will load env TYPESENSE_URL if set")
}
