package game

import (
	"fmt"
	"github.com/rnixik/go-mages/internal/lobby"
)

type Game struct {
	playersClients []lobby.ClientPlayer
	status         string
}

func (g Game) DispatchGameEvent(client lobby.ClientPlayer, event interface{}) {
	fmt.Printf("got event to game from client '%s': %+v\n", client.Nickname(), event)
}

func (g Game) OnClientRemoved(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) OnClientJoined(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) StartMainLoop() {

}

func (g Game) Status() string {
	return g.status
}

func NewGame(playersClients []lobby.ClientPlayer) *Game {
	return &Game{
		playersClients: playersClients,
	}
}
