package main

import (
	"log"
	"os"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var funcName string
var temp *template.Template

var bhFuncNameCmd = &cobra.Command{
	Use:   "bash-function",
	Short: "Generate a bash function to search through the shell history",
	Long:  "The function can be sourced in your .bashrc or .bash_profile file. or added by: source <(history-searcher bash-function)",
	Run: func(cmd *cobra.Command, args []string) {
		temp = template.Must(template.ParseFS(myTemplates, "templates/bh-func-search.tmpl"))
		err := temp.Execute(os.Stdout, funcName)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	viper.SetDefault("BH_FUNC_NAME", "bh")
	rootCmd.AddCommand(bhFuncNameCmd)
	bhFuncNameCmd.PersistentFlags().StringVar(&funcName, "name", viper.GetString("BH_FUNC_NAME"), "Name of the function to generate")
}
