package game

import (
	"github.com/rnixik/go-mages/internal/lobby"
)

type ClientPlayer interface {
	SendEvent(event interface{})
	Id() uint64
	Nickname() string
}

type Game struct {
	playersClients []ClientPlayer
}

func (g Game) DispatchGameEvent(client lobby.ClientPlayer, event interface{}) {
	panic("implement me")
}

func (g Game) OnClientRemoved(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) OnClientJoined(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) AddBotCommand(client lobby.ClientPlayer) {
	panic("implement me")
}

func (g Game) StartMainLoop() {
	panic("implement me")
}

func (g Game) Status() string {
	panic("implement me")
}

func NewGame(playersClients []ClientPlayer) *Game {
	return &Game{
		playersClients: playersClients,
	}
}
