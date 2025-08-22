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
	ShieldWorked   bool   `json:"shieldWorked"`
}

type EndGameEvent struct {
	WinnerPlayerId uint64 `json:"winnerPlayerId"`
}
