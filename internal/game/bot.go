package game

import "github.com/rnixik/go-mages/internal/lobby"

type Bot struct {
	lobby.Client
	nickname string
	id       uint64
}

func NewBot(botId uint64) lobby.ClientPlayer {
	return &Bot{
		nickname: "Bot",
		id:       botId,
	}
}

func (b *Bot) SendEvent(event interface{}) {

}

func (b *Bot) Id() uint64 {
	return b.id
}

func (b *Bot) Nickname() string {
	return b.nickname
}
