package game

type CastEvent struct {
	SpellId        string `json:"spellId"`
	OriginPlayerId uint64 `json:originPlayerId`
}
