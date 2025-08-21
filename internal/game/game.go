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

type Player struct {
	client         lobby.ClientPlayer
	lastSpellId    string
	lastCastTime   time.Time
	hasActiveSpell bool
	hp             int
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
		players:            players,
		broadcastEventFunc: broadcastEventFunc,
		mutex:              sync.Mutex{},
	}
}

func (g *Game) DispatchGameCommand(client lobby.ClientPlayer, commandName string, commandData interface{}) {
	fmt.Printf("got game command from client '%s': %+v: %+v\n", client.Nickname(), commandName, commandData)
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
		fmt.Printf("spellId: %s\n", castCommand.SpellId)
		if castCommand.SpellId == "" {
			return
		}
		g.updatePlayerSpell(client.Id(), castCommand.SpellId)
		g.broadcastEventFunc(CastEvent{SpellId: castCommand.SpellId, OriginPlayerId: client.Id()})
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
	return g.status
}

func (g *Game) updatePlayerSpell(clientID uint64, spellId string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for _, p := range g.players {
		if p.client.Id() == clientID {
			p.lastSpellId = spellId
			p.lastCastTime = time.Now()
			p.hasActiveSpell = true

			return
		}
	}
}

func (g *Game) checkAttackFromP1ToP2(p1 *Player, p2 *Player) {
	if !p1.hasActiveSpell {
		return
	}

	now := time.Now()
	castedAgo := now.Sub(p1.lastCastTime)
	if castedAgo.Milliseconds() < 300 {
		// too early to check

		return
	}

	if castedAgo.Milliseconds() > 300 {
		// to not check twice
		p1.hasActiveSpell = false
	}

	_, isShield := shieldsMap[p1.lastSpellId]
	if isShield {
		return
	}

	damage := 10

	_, hasP2Shield := shieldsMap[p2.lastSpellId]
	if hasP2Shield {
		if p2.lastCastTime.Sub(p1.lastCastTime).Milliseconds() > 300.0 {
			hasP2Shield = false
		}

		if hasP2Shield {
			var isShieldMatch bool
			switch p1.lastSpellId {
			case SpellFireball:
				isShieldMatch = p2.lastSpellId == SpellShieldFireball
			case SpellLightning:
				isShieldMatch = p2.lastSpellId == SpellShieldLightning
			case SpellComet:
				isShieldMatch = p2.lastSpellId == SpellShieldComet
			case SpellRocks:
				isShieldMatch = p2.lastSpellId == SpellShieldRocks
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
	})
}
