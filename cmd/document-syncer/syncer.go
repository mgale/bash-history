package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mgale/bash-history.git/internal/defaults"
	"github.com/mgale/bash-history.git/internal/events"
	"github.com/mgale/bash-history.git/internal/natsioclient"
	"github.com/spf13/cobra"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/xorcare/pointer"
)

var client *natsioclient.NatsIOClient
var tsclient *typesense.Client
var verbose bool

var consumeCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync events to the search engine",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		streamEventsChannel := make(chan events.DocumentEvent, 1000)
		client, err = natsioclient.NewNatsIOClient(natsURL, natsDocumentSynerTopic)
		if err != nil {
			log.Fatal(err)
		}
		client.LogConnInfo()

		tsclient = typesense.NewClient(typesense.WithServer(typesenseURL),
			typesense.WithAPIKey(typesenseAPIKey),
			typesense.WithConnectionTimeout(5*time.Second),
			typesense.WithCircuitBreakerMaxRequests(50),
			typesense.WithCircuitBreakerInterval(2*time.Minute),
			typesense.WithCircuitBreakerTimeout(1*time.Minute))

		err = createCollection(tsclient)
		if err != nil {
			log.Fatal(err)
		}
		client.EncodedConn.BindRecvChan(natsDocumentSynerTopic, streamEventsChannel)
		client.EncodedConn.PublishRequest(natsCommandTopic, natsDocumentSynerTopic, "loadAllHistoryEvents")
		err = handleEvents(cmd.Context(), streamEventsChannel, verbose, tsclient)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)
	consumeCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func handleEvents(ctx context.Context, streamEventsChannel chan events.DocumentEvent, verbose bool, tsclient *typesense.Client) error {

	docBufferMax := 500
	docBuffer := make([]events.Document, 0, docBufferMax)
	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-streamEventsChannel:
			if verbose {
				log.Printf("Received event: %+v", event)
			}
			//Send to search engine
			doc := events.Document{
				ID:        fmt.Sprint(event.ID),
				Timestamp: event.Timestamp,
				Username:  event.Username,
				Command:   event.Command,
			}
			docBuffer = append(docBuffer, doc)
			if len(docBuffer) >= docBufferMax {
				err := bulkDocUpload(&docBuffer, tsclient)
				if err != nil {
					return err
				}
			}
		case <-time.After(3 * time.Second):
			if verbose {
				log.Println("No events received for 3 second")
			}
			if len(docBuffer) != 0 {
				err := bulkDocUpload(&docBuffer, tsclient)
				if err != nil {
					return err
				}
			}
		}
	}
}

func bulkDocUpload(docBuffer *[]events.Document, tsclient *typesense.Client) error {
	log.Println("Uploading", len(*docBuffer), "documents to typesense")
	importParams := &api.ImportDocumentsParams{
		Action:    pointer.String("upsert"),
		BatchSize: pointer.Int(len(*docBuffer)),
	}

	anything := make([]interface{}, 0, len(*docBuffer))
	for _, doc := range *docBuffer {
		anything = append(anything, doc)
	}
	results, err := tsclient.Collection(defaults.CollectionName).Documents().Import(anything, importParams)
	if err != nil {
		return err
	}
	for _, r := range results {
		if !r.Success {
			log.Printf("Imported document success: %v, error: %s", r.Success, r.Error)
		}
	}
	log.Println("Done uploading", len(*docBuffer), "documents to typesense")
	*docBuffer = (*docBuffer)[:0]
	return nil
}

func createCollection(tsclient *typesense.Client) error {

	_, err := tsclient.Collection(defaults.CollectionName).Retrieve()
	//If collection exists, delete it
	if err == nil {
		log.Println("Deleting collection", defaults.CollectionName)
		_, err := tsclient.Collection(defaults.CollectionName).Delete()
		if err != nil {
			return err
		}
	}

	tokenSeparators := []string{"-", "/", "_", "="}

	collectionSchema := &api.CollectionSchema{
		Name: "bash-history",
		Fields: []api.Field{
			{
				Name: "timestamp",
				Type: "int64",
			},
			{
				Name:  "username",
				Type:  "string",
				Facet: pointer.Bool(true),
			},
			{
				Name:  "command",
				Type:  "string",
				Infix: pointer.Bool(true),
			},
		},
		DefaultSortingField: stringToPtr("timestamp"),
		TokenSeparators:     &tokenSeparators,
	}
	log.Println("Creating collection", defaults.CollectionName)
	_, err = tsclient.Collections().Create(collectionSchema)
	return err

}

func stringToPtr(s string) *string {
	return &s
}
