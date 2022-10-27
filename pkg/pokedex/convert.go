package pokedex

import (
	"github.com/goccy/go-json"
	"github.com/hieubq90/pokedex/pkg/models"
	"os"
)

func SaveToJson(listPokemon []*models.Pokemon) {
	jsonData, _ := json.MarshalIndent(listPokemon, "", " ")
	_ = os.WriteFile("pokedex.json", jsonData, 0644)
}

func SaveToCsv(listPokemon []*models.Pokemon) {
	jsonData, _ := json.MarshalIndent(listPokemon, "", " ")
	_ = os.WriteFile("pokedex.csv", jsonData, 0644)
}

func ConvertToCsv(listPokemon []*models.Pokemon) string {
	return ""
}
