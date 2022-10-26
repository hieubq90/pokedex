package pokedex

import (
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"convert"},
	Short:   "Convert Pokemon data between CSV & JSON",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
