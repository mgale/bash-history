package main

import (
	"bufio"
	"log"
	"os"
	"os/user"

	"github.com/mgale/bash-history.git/internal/events"
	"github.com/spf13/cobra"
)

var historyFile string

var loadBashHistoryFile = &cobra.Command{
	Use:   "load-bash-history",
	Short: "Load existing history file into the database",
	Run: func(cmd *cobra.Command, args []string) {

		// Create an events detailed channel
		streamEventsChannel := make(chan events.ReadEvent, 1000)

		client.LogConnInfo()
		// Setup NATIOS to consume the events detailed channel
		client.EncodedConn.BindSendChan(natsTopic, streamEventsChannel)
		log.Println("streaming data to: ", natsURL)

		currentUser, err := user.Current()
		if err != nil {
			log.Fatalf(err.Error())
		}

		file, err := os.Open(historyFile)
		if err != nil {
			log.Fatalf("failed to open")

		}

		scanner := bufio.NewScanner(file)

		// The bufio.ScanLines is used as an
		// input to the method bufio.Scanner.Split()
		// and then the scanning forwards to each
		// new line using the bufio.Scanner.Scan()
		// method.
		scanner.Split(bufio.ScanLines)
		var text []string

		for scanner.Scan() {
			text = append(text, scanner.Text())
		}

		file.Close()

		// and then a loop iterates through
		// and prints each of the slice values.
		counter := 0
		log.Println("Loading file into database...")
		for _, each_ln := range text {
			e := events.ReadEvent{
				Pid:      0,
				Username: currentUser.Username,
				Line:     each_ln,
			}
			streamEventsChannel <- e
			counter++
		}
		log.Println("Loaded ", counter, " lines into database")
	},
}

func init() {
	rootCmd.AddCommand(loadBashHistoryFile)
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	loadBashHistoryFile.PersistentFlags().StringVar(&historyFile, "file", dirname+"/.bash_history", "file to load")

}
