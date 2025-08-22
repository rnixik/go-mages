package game

import (
	"github.com/rnixik/go-mages/internal/lobby"
	"log"
	"math/rand"
	"time"
)

type Bot struct {
	botClient *BotClient
	room      *lobby.Room
}

func newBot(botClient *BotClient, room *lobby.Room) *Bot {
	return &Bot{
		botClient: botClient,
		room:      room,
	}
}

func (b *Bot) run() {
	ticker := time.NewTicker(1300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if b.room.Game() != nil && b.room.Game().Status() == GameStatusEnded {
				return
			}

			random := rand.Intn(4)
			var command *CastCommand
			switch random {
			case 0:
				command = &CastCommand{SpellFireball}
			case 1:
				command = &CastCommand{SpellLightning}
			case 2:
				command = &CastCommand{SpellRocks}
			case 3:
				command = &CastCommand{SpellComet}
			}

			b.botClient.sendCommandToGame("Cast", command)
		case event := <-b.botClient.incomingEvents:
			b.dispatchEvent(event)
		}
	}
}

func (b *Bot) dispatchEvent(event interface{}) {
	log.Printf("BOT: got event to make decision: %+v\n", event)
	castEvent, ok := event.(CastEvent)
	if !ok {
		return
	}
	if castEvent.OriginPlayerId == b.botClient.Id() {
		return
	}
	log.Printf("BOT: got spell %s", castEvent.SpellId)
}
