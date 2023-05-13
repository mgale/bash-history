package main

import (
	"log"

	"github.com/mgale/bash-history.git/internal/events"
	"github.com/mgale/bash-history.git/internal/natsioclient"
	"github.com/spf13/cobra"
)

var client *natsioclient.NatsIOClient
var verbose bool

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "consume events",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		streamEventsChannel := make(chan events.ReadEvent, 1000)
		docEventsChannel := make(chan events.DocumentEvent, 1000)
		client, err = natsioclient.NewNatsIOClient(natsURL, natsTopic)
		if err != nil {
			log.Fatal(err)
		}
		client.LogConnInfo()

		err = createDBFileIfNotExists(dbFilename)
		if err != nil {
			log.Fatal(err)
		}

		db, err := createDBConnection(dbFilename)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		err = createTables(db)
		if err != nil {
			log.Fatal(err)
		}

		err = createTableIndexes(db)
		if err != nil {
			log.Fatal(err)
		}

		client.EncodedConn.BindRecvChan(natsTopic, streamEventsChannel)
		client.EncodedConn.BindSendChan(natsDocumentSynerTopic, docEventsChannel)
		client.EncodedConn.Subscribe(natsCommandTopic, handleHistoricalDataLoad)
		err = handleEvents(cmd.Context(), streamEventsChannel, docEventsChannel, db, verbose)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)
	consumeCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func handleHistoricalDataLoad(subject, reply, command string) {
	log.Printf("Received command: %s", command)
	db, err := createDBConnection(dbFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM bashhistory")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var docEvent events.DocumentEvent
		err := rows.Scan(&docEvent.ID, &docEvent.Timestamp, &docEvent.Username, &docEvent.Command)
		if err != nil {
			log.Fatal(err)
		}
		client.EncodedConn.Publish(reply, docEvent)
	}
}
