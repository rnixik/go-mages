package game

import (
	"encoding/json"
	"github.com/rnixik/go-mages/internal/lobby"
	"log"
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

// maxShieldCastDiffMs is the maximum time difference in milliseconds between the shield spell cast and the attack spell cast
const maxShieldCastDiffMs = int64(900)

// attackCastDelayMs is the minimum time in milliseconds between two attack spell casts
const attackCastDelayMs = int64(2000)

// shieldCastDelayMs is the minimum time in milliseconds between two shield spell casts
const shieldCastDelayMs = int64(900)

const maxHP = 1000

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
		hp:             maxHP,
	}
}

type Game struct {
	players            []*Player
	status             string
	broadcastEventFunc func(event interface{})
	mutex              sync.Mutex
	statusMx           sync.Mutex
	room               *lobby.Room
}

func NewGame(playersClients []lobby.ClientPlayer, room *lobby.Room, broadcastEventFunc func(event interface{})) *Game {
	players := make([]*Player, len(playersClients))
	for i, client := range playersClients {
		players[i] = newPlayer(client)
	}

	log.Printf("new game created: %s vs %s\n", players[0].client.Nickname(), players[1].client.Nickname())

	return &Game{
		status:             StatusStarted,
		players:            players,
		broadcastEventFunc: broadcastEventFunc,
		mutex:              sync.Mutex{},
		room:               room,
	}
}

func (g *Game) DispatchGameCommand(client lobby.ClientPlayer, commandName string, commandData interface{}) {
	if g.isGameEnded() {
		return
	}

	eventDataJson, ok := commandData.(json.RawMessage)
	if !ok {
		log.Printf("cannot decode event data for event name = %s\n", commandName)
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
		g.updatePlayerSpell(client.ID(), castCommand.SpellId)
		break
	}
}

func (g *Game) OnClientRemoved(client lobby.ClientPlayer) {
	if g.isGameEnded() {
		return
	}
	winnerID := uint64(0)
	for _, p := range g.players {
		if p.client.ID() != client.ID() {
			winnerID = p.client.ID()
		}
	}
	g.endGame(winnerID)
}

func (g *Game) OnClientJoined(client lobby.ClientPlayer) {
	log.Printf("client '%s' joined game\n", client.Nickname())
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
	g.statusMx.Lock()
	defer g.statusMx.Unlock()

	return g.status
}

func (g *Game) updatePlayerSpell(clientID uint64, spellId string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, p := range g.players {
		now := time.Now()
		if p.client.ID() == clientID {
			_, isShield := shieldsMap[spellId]
			if isShield {
				if p.lastSpellIdShield != "" && now.Sub(p.lastCastTimeShield).Milliseconds() < shieldCastDelayMs {
					return
				}
				p.lastSpellIdShield = spellId
				p.lastCastTimeShield = now
				p.spellWasSentShield = false
			} else {
				if p.lastSpellId != "" && now.Sub(p.lastCastTime).Milliseconds() < attackCastDelayMs {
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

func (g *Game) endGame(winnerPlayerId uint64) {
	g.statusMx.Lock()
	g.status = StatusEnded
	g.statusMx.Unlock()

	g.broadcastEventFunc(EndGameEvent{WinnerPlayerId: winnerPlayerId})
	g.room.OnGameEnded()
}

func (g *Game) checkAttackFromP1ToP2(p1 *Player, p2 *Player) {
	if p2.hp <= 0 {
		g.endGame(p1.client.ID())
	}

	if p1.lastSpellIdShield != "" && !p1.spellWasSentShield {
		p1.spellWasSentShield = true
		g.broadcastEventFunc(CastEvent{SpellId: p1.lastSpellIdShield, OriginPlayerId: p1.client.ID()})
	}

	if !p1.hasActiveSpell {
		return
	}

	if !p1.spellWasSent {
		p1.spellWasSent = true
		g.broadcastEventFunc(CastEvent{SpellId: p1.lastSpellId, OriginPlayerId: p1.client.ID()})
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

	damage := 100
	var isShieldMatch bool

	if p2.lastSpellIdShield != "" {
		shieldCastDiff := p2.lastCastTimeShield.Sub(p1.lastCastTime).Milliseconds()
		if shieldCastDiff < 0 || shieldCastDiff > maxShieldCastDiffMs {
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
				damage = int(shieldCastDiff / maxShieldCastDiffMs * 100)
			}
		}
	}

	p2.hp -= damage

	g.broadcastEventFunc(DamageEvent{
		SpellId:        p1.lastSpellId,
		Damage:         damage,
		TargetPlayerId: p2.client.ID(),
		TargetPlayerHp: p2.hp,
		ShieldWorked:   isShieldMatch,
	})
}

func (g *Game) isGameEnded() bool {
	g.statusMx.Lock()
	defer g.statusMx.Unlock()

	return g.status == StatusEnded
}
