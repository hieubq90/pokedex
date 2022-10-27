package models

type Pokemon struct {
	ID         int     `json:"id"`
	Num        int     `json:"num"`
	Name       string  `json:"name"`
	Sprite     string  `json:"sprite"`
	Image      string  `json:"image"`
	Types      string  `json:"types"`
	Species    string  `json:"species"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	Evolution  string  `json:"evolution"`
	Weaknesses string  `json:"weaknesses"`
	HP         int     `json:"hp"`
	Attack     int     `json:"attack"`
	Defense    int     `json:"defense"`
	SpAttack   int     `json:"sp_attack"`
	SpDefense  int     `json:"sp_defense"`
	Speed      int     `json:"speed"`
	Total      int     `json:"total"`
	PokedexUrl string
}
