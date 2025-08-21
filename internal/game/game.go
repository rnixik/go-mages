package game

import (
	"encoding/json"
	"fmt"
	"github.com/rnixik/go-mages/internal/lobby"
)

type Game struct {
	playersClients     []lobby.ClientPlayer
	status             string
	broadcastEventFunc func(event interface{})
}

func (g Game) DispatchGameCommand(client lobby.ClientPlayer, commandName string, commandData interface{}) {
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
		g.broadcastEventFunc(CastEvent{SpellId: castCommand.SpellId, OriginPlayerId: client.Id()})
		break
	}
}

func (g Game) OnClientRemoved(client lobby.ClientPlayer) {
	fmt.Printf("client '%s' removed from game\n", client.Nickname())
}

func (g Game) OnClientJoined(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) StartMainLoop() {

}

func (g Game) Status() string {
	return g.status
}

func NewGame(playersClients []lobby.ClientPlayer, broadcastEventFunc func(event interface{})) *Game {
	return &Game{
		playersClients:     playersClients,
		broadcastEventFunc: broadcastEventFunc,
	}
}
