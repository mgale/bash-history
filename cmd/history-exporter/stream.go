package main

import (
	"log"

	"github.com/mgale/bash-history.git/internal/events"
	"github.com/spf13/cobra"
)

var verbose bool

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream readline events to a remote server",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		// Create an events detailed channel
		streamEventsChannel := make(chan events.ReadEvent, 1000)

		client.LogConnInfo()
		// Setup NATIOS to consume the events detailed channel
		client.EncodedConn.BindSendChan(natsTopic, streamEventsChannel)
		log.Println("streaming data to: ", natsURL)
		err = readEvents(cmd.Context(), streamEventsChannel, verbose)
		if err != nil {
			log.Printf("stopped reading events: %s", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

}
