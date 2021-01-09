package game

import "github.com/rnixik/go-mages/internal/lobby"

type MatchMaker struct {
	waitingClientId uint64
}

func NewMatchMaker() *MatchMaker {
	return &MatchMaker{}
}

func (mm *MatchMaker) MakeMatch(client lobby.ClientPlayer, foundFunc func(clientsIds []lobby.ClientPlayer), notFoundFunc func(), addBotFunc func() lobby.ClientPlayer) {
	// only with bots
	botClient := addBotFunc()
	foundFunc([]lobby.ClientPlayer{client, botClient})

	// game with players
	//if mm.waitingClientId != 0 {
	//	foundFunc([]uint64{clientId, mm.waitingClientId})
	//	mm.waitingClientId = 0
	//	return
	//}
	//
	//mm.waitingClientId = clientId
}

func (mm *MatchMaker) Cancel(client lobby.ClientPlayer) {
	if mm.waitingClientId == client.Id() {
		mm.waitingClientId = 0
	}
}
