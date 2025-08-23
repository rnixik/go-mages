package game

import (
	"encoding/json"
	"fmt"
	"github.com/rnixik/go-mages/internal/lobby"
	"sync"
	"time"
)

const SpellShieldFireball = "protect_fireball"
const SpellShieldLightning = "protect_lightning"
const SpellShieldComet = "protect_comet"
const SpellShieldRocks = "protect_rocks"
const SpellFireball = "fireball"
const SpellLightning = "lightning"
const SpellComet = "comet"
const SpellRocks = "rocks"

var shieldsMap = map[string]bool{
	SpellShieldFireball:  true,
	SpellShieldLightning: true,
	SpellShieldComet:     true,
	SpellShieldRocks:     true,
}

const StatusStarted = "started"
const StatusEnded = "ended"

type Player struct {
	client             lobby.ClientPlayer
	lastSpellId        string
	lastSpellIdShield  string
	lastCastTime       time.Time
	lastCastTimeShield time.Time
	spellWasSent       bool
	spellWasSentShield bool
	hasActiveSpell     bool
	hp                 int
}

func newPlayer(client lobby.ClientPlayer) *Player {
	return &Player{
		client:         client,
		lastSpellId:    "",
		lastCastTime:   time.Time{},
		hasActiveSpell: false,
		hp:             100,
	}
}

type Game struct {
	players            []*Player
	status             string
	broadcastEventFunc func(event interface{})
	mutex              sync.Mutex
}

func NewGame(playersClients []lobby.ClientPlayer, broadcastEventFunc func(event interface{})) *Game {
	players := make([]*Player, len(playersClients))
	for i, client := range playersClients {
		players[i] = newPlayer(client)
	}
	return &Game{
		status:             StatusStarted,
		players:            players,
		broadcastEventFunc: broadcastEventFunc,
		mutex:              sync.Mutex{},
	}
}

func (g *Game) DispatchGameCommand(client lobby.ClientPlayer, commandName string, commandData interface{}) {
	if g.isGameEnded() {
		return
	}

	eventDataJson, ok := commandData.(json.RawMessage)
	if !ok {
		fmt.Printf("cannot decode event data for event name = %s\n", commandName)
		return
	}

	switch commandName {
	case "Cast":
		var castCommand CastCommand
		if err := json.Unmarshal(eventDataJson, &castCommand); err != nil {
			return
		}
		if castCommand.SpellId == "" {
			return
		}
		g.updatePlayerSpell(client.Id(), castCommand.SpellId)
		break
	}
}

func (g *Game) OnClientRemoved(client lobby.ClientPlayer) {
	fmt.Printf("client '%s' removed from game\n", client.Nickname())
}

func (g *Game) OnClientJoined(client lobby.ClientPlayer) {
	fmt.Printf("client '%s' joined game\n", client.Nickname())
}

func (g *Game) StartMainLoop() {
	// perform action every 50 milliseconds
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if g.isGameEnded() {
				return
			}

			g.mutex.Lock()
			p1 := g.players[0]
			p2 := g.players[1]
			g.checkAttackFromP1ToP2(p1, p2)
			g.checkAttackFromP1ToP2(p2, p1)
			g.mutex.Unlock()
		}
	}
}

func (g *Game) Status() string {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.status
}

func (g *Game) updatePlayerSpell(clientID uint64, spellId string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, p := range g.players {
		now := time.Now()
		if p.client.Id() == clientID {
			_, isShield := shieldsMap[spellId]
			if isShield {
				if p.lastSpellIdShield != "" && now.Sub(p.lastCastTimeShield).Milliseconds() < 900 {
					return
				}
				p.lastSpellIdShield = spellId
				p.lastCastTimeShield = now
				p.spellWasSentShield = false
			} else {
				if p.lastSpellId != "" && now.Sub(p.lastCastTime).Milliseconds() < 1000 {
					return
				}
				p.lastSpellId = spellId
				p.lastCastTime = now
				p.hasActiveSpell = true
				p.spellWasSent = false
			}

			return
		}
	}
}

func (g *Game) checkAttackFromP1ToP2(p1 *Player, p2 *Player) {
	if p2.hp <= 0 {
		g.status = StatusEnded
		g.broadcastEventFunc(EndGameEvent{WinnerPlayerId: p1.client.Id()})
	}

	if p1.lastSpellIdShield != "" && !p1.spellWasSentShield {
		p1.spellWasSentShield = true
		g.broadcastEventFunc(CastEvent{SpellId: p1.lastSpellIdShield, OriginPlayerId: p1.client.Id()})
	}

	if !p1.hasActiveSpell {
		return
	}

	if !p1.spellWasSent {
		p1.spellWasSent = true
		g.broadcastEventFunc(CastEvent{SpellId: p1.lastSpellId, OriginPlayerId: p1.client.Id()})
	}

	var prepareDurationMs int64 = 500

	now := time.Now()
	castedAgo := now.Sub(p1.lastCastTime)
	if castedAgo.Milliseconds() < prepareDurationMs {
		// too early to check

		return
	}

	if castedAgo.Milliseconds() >= prepareDurationMs {
		// to not check twice
		p1.hasActiveSpell = false
	}

	damage := 10
	var isShieldMatch bool

	if p2.lastSpellIdShield != "" {
		shieldCastDiff := p2.lastCastTimeShield.Sub(p1.lastCastTime).Milliseconds()
		if shieldCastDiff < 0 || shieldCastDiff > 400.0 {
			p2.lastSpellIdShield = ""
		}

		if p2.lastSpellIdShield != "" {
			switch p1.lastSpellId {
			case SpellFireball:
				isShieldMatch = p2.lastSpellIdShield == SpellShieldFireball
			case SpellLightning:
				isShieldMatch = p2.lastSpellIdShield == SpellShieldLightning
			case SpellComet:
				isShieldMatch = p2.lastSpellIdShield == SpellShieldComet
			case SpellRocks:
				isShieldMatch = p2.lastSpellIdShield == SpellShieldRocks
			default:
				isShieldMatch = false
			}

			if isShieldMatch {
				damage = 2
			}
		}
	}

	p2.hp -= damage

	g.broadcastEventFunc(DamageEvent{
		SpellId:        p1.lastSpellId,
		Damage:         damage,
		TargetPlayerId: p2.client.Id(),
		TargetPlayerHp: p2.hp,
		ShieldWorked:   isShieldMatch,
	})
}

func (g *Game) isGameEnded() bool {
	return g.status == StatusEnded
}
