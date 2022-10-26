package pokedex

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hieubq90/pokedex/pkg/pokedex"
	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:     "crawl",
	Aliases: []string{"crawl"},
	Short:   "Crawl pokemon data from https://pokemondb.net/",
	Run: func(cmd *cobra.Command, args []string) {
		var answer FileTypeAnswer
		err := survey.Ask(fileTypeQuestion, &answer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Start crawling data & save to %s file\n", answer.FileType)
		res := pokedex.Crawl()
		pokedex.CrawlDetailInfo(res)
		//fmt.Println(res)
		for _, p := range res {
			fmt.Println(*p)
		}
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)
}
