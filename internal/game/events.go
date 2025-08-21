package game

type CastEvent struct {
	SpellId        string `json:"spellId"`
	OriginPlayerId uint64 `json:originPlayerId`
}

type DamageEvent struct {
	SpellId        string `json:"spellId"`
	Damage         int    `json:"damage"`
	TargetPlayerId uint64 `json:"targetPlayerId"`
	TargetPlayerHp int    `json:"targetPlayerHp"`
}
