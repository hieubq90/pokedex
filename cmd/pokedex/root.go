package pokedex

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pokedex",
	Short: "pokedex is a crawler & converter data of pokemon",
	Long:  `A CLI application help crawling pokemon data & convert between CSV & JSON`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
