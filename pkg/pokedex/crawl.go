package pokedex

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/hieubq90/pokedex/pkg/models"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
)

const NationalUrl string = "https://pokemondb.net/pokedex/national"
const SpriteUrl = "https://img.pokemondb.net/sprites/home/normal/2x/avif/%s.avif"
const ImageUrl = "https://img.pokemondb.net/artwork/%s.jpg"

func parsePokemonNum(e *colly.HTMLElement) int {
	numStr := ""
	e.ForEachWithBreak("small", func(index int, e *colly.HTMLElement) bool {
		if index > 0 {
			return false
		}
		numStr = strings.Replace(e.Text, "#", "", 1)
		return true
	})
	num, _ := strconv.Atoi(numStr)
	return num
}

func parsePokemonTypes(e *colly.HTMLElement) string {
	types := ""
	e.ForEach(".itype", func(index int, e *colly.HTMLElement) {
		if index > 0 {
			types += fmt.Sprintf(", %s", e.Text)
		} else {
			types += e.Text
		}
	})
	return types
}

func Crawl() []*models.Pokemon {
	_ret := make([]*models.Pokemon, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("pokemondb.net"),
	)

	c.OnHTML(".infocard", func(e *colly.HTMLElement) {
		name := e.ChildText(".ent-name")
		pokedexUrl := e.Request.AbsoluteURL(e.ChildAttr(".ent-name", "href"))
		num := parsePokemonNum(e)
		sprite := fmt.Sprintf(SpriteUrl, strings.ToLower(name))
		image := fmt.Sprintf(ImageUrl, strings.ToLower(name))
		types := parsePokemonTypes(e)
		_ret = append(_ret, &models.Pokemon{
			ID:         num,
			Num:        num,
			Name:       name,
			Sprite:     sprite,
			Image:      image,
			Types:      types,
			PokedexUrl: pokedexUrl,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(NationalUrl)
	if err != nil {
		fmt.Println(err.Error())
	}

	return _ret
}

func CrawlDetailInfo(listPokemon []*models.Pokemon) {
	for i, p := range listPokemon {
		if i > 10 {
			break
		}
		c := colly.NewCollector(
			colly.AllowedDomains("pokemondb.net"),
		)

		tabId := fmt.Sprintf("#tab-basic-%d", p.Num)

		c.OnHTML(tabId, func(e *colly.HTMLElement) {
			e.ForEach(".vitals-table", func(i int, element *colly.HTMLElement) {
				if i == 0 {
					// Pokedex data
					fmt.Println("Pokedex data\n", element.Text)
				} else if i == 3 {
					// Base stats
					fmt.Println("Base stats\n", element.Text)
				}
			})
			weaknesses := make([]string, 0)
			e.ForEach("td.type-fx-cell", func(_ int, element *colly.HTMLElement) {
				if element.Text == "2" {
					weak := strings.Split(element.Attr("title"), " ")[0]
					if !slices.Contains(weaknesses, weak) {
						weaknesses = append(weaknesses, weak)
					}
				}
			})
			fmt.Println("weaknesses:", strings.Join(weaknesses[:], ", "))
			p.Weaknesses = strings.Join(weaknesses[:], ", ")
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		err := c.Visit(p.PokedexUrl)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}
