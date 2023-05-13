package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mgale/bash-history.git/internal/defaults"
	"github.com/mgale/bash-history.git/internal/natsioclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var natsURL string
var natsTopic string
var client *natsioclient.NatsIOClient
var err error

var rootCmd = &cobra.Command{
	Use:   "history-exporter",
	Short: "history-exporter is a command line tool for tracing bash readline calls",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		client, err = natsioclient.NewNatsIOClient(natsURL, natsTopic)
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		client.Close()
	},
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
	viper.SetDefault("NATSIO_TOPIC", defaults.BashHistoryTopic)
	viper.AutomaticEnv()
	rootCmd.PersistentFlags().StringVar(&natsURL, "natsio-url", viper.GetString("NATSIO_URL"), "Natsio URL, will load env NATSIO_URL if set")
	rootCmd.PersistentFlags().StringVar(&natsTopic, "natsio-topic", viper.GetString("NATSIO_TOPIC"), "Natsio topic, will load env NATSIO_TOPIC if set")
}
