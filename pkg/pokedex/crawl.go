package pokedex

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/hieubq90/pokedex/pkg/models"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
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
		fmt.Println("Crawling", r.URL)
	})

	err := c.Visit(NationalUrl)
	if err != nil {
		fmt.Println(err.Error())
	}

	return _ret
}

func CrawlDetailInfo(listPokemon []*models.Pokemon) {
	total := len(listPokemon)
	doneCh := make(chan struct{})
	desc := fmt.Sprintf("[cyan][%d/%d][reset] Crawling Pokemon's detail infomation", 0, total)
	bar := progressbar.NewOptions(total,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			doneCh <- struct{}{}
		}),
	)

	go func() {
		for i, p := range listPokemon {
			c := colly.NewCollector(
				colly.AllowedDomains("pokemondb.net"),
			)

			tabId := fmt.Sprintf("#tab-basic-%d", p.Num)

			// parse Pokedex data & Base Stats
			c.OnHTML(tabId, func(e *colly.HTMLElement) {
				e.ForEach(".vitals-table", func(i int, element *colly.HTMLElement) {
					if i == 0 {
						// Pokedex data
						parsePokedexData(p, element)
					} else if i == 3 {
						// Base stats
						parseBaseStats(p, element)
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
				p.Weaknesses = strings.Join(weaknesses[:], ", ")
			})

			// parse Evolution
			c.OnHTML(".infocard-list-evo", func(e *colly.HTMLElement) {
				evolution := make([]string, 0)
				e.ForEach("div.infocard", func(_ int, element *colly.HTMLElement) {
					num := parsePokemonNum(element)
					evolution = append(evolution, strconv.Itoa(num))
				})
				p.Evolution = strings.Join(evolution[:], ", ")
			})

			c.OnRequest(func(r *colly.Request) {
				desc = fmt.Sprintf("[cyan][%d/%d][reset] Processing: %s", i+1, total, r.URL)
				_ = bar.Add(1)
				bar.Describe(desc)
			})

			err := c.Visit(p.PokedexUrl)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()

	<-doneCh
}

func parsePokedexData(p *models.Pokemon, e *colly.HTMLElement) {
	e.ForEach("td", func(i int, element *colly.HTMLElement) {
		if i == 2 {
			// Species
			p.Species = element.Text
		} else if i == 3 {
			// Height
			heightStr := strings.Split(element.Text, " ")[0]
			height, _ := strconv.ParseFloat(heightStr, 64)
			p.Height = height
		} else if i == 4 {
			// Weight
			weightStr := strings.Split(element.Text, " ")[0]
			weight, _ := strconv.ParseFloat(weightStr, 64)
			p.Weight = weight
		}
	})
}

func parseBaseStatValue(e *colly.HTMLElement) int {
	statStr := ""
	e.ForEachWithBreak("td.cell-num", func(i int, element *colly.HTMLElement) bool {
		if i == 0 {
			statStr = element.Text
			return true
		}
		return false
	})
	val, _ := strconv.Atoi(statStr)
	return val
}

func parseBaseStats(p *models.Pokemon, e *colly.HTMLElement) {
	e.ForEach("tr", func(i int, element *colly.HTMLElement) {
		if i == 0 {
			// HP
			p.HP = parseBaseStatValue(element)
		} else if i == 1 {
			// Attack
			p.Attack = parseBaseStatValue(element)
		} else if i == 2 {
			// Defense
			p.Defense = parseBaseStatValue(element)
		} else if i == 3 {
			// Sp. Atk
			p.SpAttack = parseBaseStatValue(element)
		} else if i == 4 {
			// Sp. Def
			p.SpDefense = parseBaseStatValue(element)
		} else if i == 5 {
			// Speed
			p.Speed = parseBaseStatValue(element)
		}
	})

	totalStr := e.ChildText("td.cell-total")
	total, _ := strconv.Atoi(totalStr)
	p.Total = total
}
