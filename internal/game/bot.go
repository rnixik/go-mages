package game

import (
	"log"
	"math/rand"
)

type Bot struct {
	botClient *BotClient
}

func newBot(botClient *BotClient) *Bot {
	return &Bot{
		botClient: botClient,
	}
}

func (b *Bot) run() {
	for {
		select {
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

	random := rand.Intn(4)
	var command *CastCommand
	switch random {
	case 0:
		command = &CastCommand{"fireball"}
	case 1:
		command = &CastCommand{"lightning"}
	case 2:
		command = &CastCommand{"rocks"}
	case 3:
		command = &CastCommand{"comet"}
	}

	b.botClient.sendCommandToGame("Cast", command)
}
