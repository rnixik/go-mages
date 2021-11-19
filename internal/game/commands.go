package game

type DemoCommand struct {
	DemoMessage string `json:"demoMessage"`
}

type CastCommand struct {
	SpellId string `json:"spellId"`
}
