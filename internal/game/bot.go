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

	random := rand.Intn(2)
	var command *CastCommand
	if random == 0 {
		command = &CastCommand{"fireball"}
	} else {
		command = &CastCommand{"lightning"}
	}

	b.botClient.sendCommandToGame("Cast", command)
}
