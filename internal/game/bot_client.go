package game

import (
	"github.com/rnixik/go-mages/internal/lobby"
)

type BotClient struct {
	lobby.Client
	id              uint64
	incomingEvents  chan interface{}
	outgoingActions chan interface{}
}

func NewBotClient(botId uint64) lobby.ClientPlayer {
	botClient := &BotClient{
		id:              botId,
		incomingEvents:  make(chan interface{}),
		outgoingActions: make(chan interface{}),
	}
	botClient.SetNickname("BotClient")

	bot := newBot(botClient)
	go botClient.sendingActionsToGame()
	go bot.run()

	return botClient
}

func (bc *BotClient) SendEvent(event interface{}) {
	bc.incomingEvents <- event
}

func (bc *BotClient) Id() uint64 {
	return bc.id
}

func (bc *BotClient) sendingActionsToGame() {
	for {
		select {
		case playerAction := <-bc.outgoingActions:
			bc.Room().Game().DispatchGameEvent(bc, playerAction)
		}
	}
}
