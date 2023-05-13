package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/typesense/typesense-go/typesense"
)

var tsclient *typesense.Client
var dryrun bool
var outPutFile string
var debugLogging bool
var cmdStr string

var searchCmd = &cobra.Command{
	Use:   "search [search term] wrapped with quotes",
	Short: "search through the shell history",
	Run: func(cmd *cobra.Command, args []string) {

		tsclient = typesense.NewClient(typesense.WithServer(typesenseURL),
			typesense.WithAPIKey(typesenseAPIKey),
			typesense.WithConnectionTimeout(3*time.Second),
			typesense.WithCircuitBreakerMaxRequests(50),
			typesense.WithCircuitBreakerInterval(2*time.Minute),

			typesense.WithCircuitBreakerTimeout(1*time.Minute))

		if debugLogging {
			f, err := tea.LogToFile("/tmp/bash_history_debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}
			defer f.Close()
			log.SetOutput(f)

		} else {
			log.SetOutput(ioutil.Discard)
		}
		log.Println("#############################################")
		log.Println("Starting debug log")

		defaultQuery := convertArgsToQuery(args)
		p := tea.NewProgram(initialModel(defaultQuery))
		_, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
		if len(cmdStr) > 0 {
			fmt.Println("# Executing command: ", cmdStr)
			if dryrun {
				fmt.Println("Dryrun: not executing command")
			} else {
				writeCmdToFile(cmdStr, outPutFile)
			}
		}
	},
}

func init() {
	viper.SetDefault("BH_DEBUG", false)
	viper.SetDefault("BH_DRYRUN", false)
	viper.SetDefault("BH_OUTPUT_FILE", false)
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Print the command to be executed instead of executing it.")
	searchCmd.PersistentFlags().StringVar(&outPutFile, "output-file", viper.GetString("BH_OUTPUT_FILE"), "Output file to write the command to be executed to.")
	searchCmd.PersistentFlags().BoolVar(&debugLogging, "debug", viper.GetBool("BH_DEBUG"), "Enable debug logging, can use env: BH_DEBUG")
}

// convertArgsToQuery converts the args to a query string
func convertArgsToQuery(args []string) []rune {
	var defaultQuery []rune
	if len(args) > 0 {
		defaultQuery = []rune(strings.Join(args, ""))
	} else {
		defaultQuery = []rune{'*'}
	}

	return defaultQuery
}
