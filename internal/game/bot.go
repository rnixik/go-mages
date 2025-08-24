package game

import (
	"github.com/rnixik/go-mages/internal/lobby"
	"math/rand"
	"time"
)

type Bot struct {
	botClient          *BotClient
	room               *lobby.Room
	delayedCastCommand *CastCommand
}

func newBot(botClient *BotClient, room *lobby.Room) *Bot {
	return &Bot{
		botClient: botClient,
		room:      room,
	}
}

func (b *Bot) run() {
	attackTicker := time.NewTicker(1300 * time.Millisecond)
	defer attackTicker.Stop()

	defenseTicker := time.NewTicker(500 * time.Millisecond)
	defer defenseTicker.Stop()

	for {
		select {
		case <-attackTicker.C:
			if b.room.Game() != nil && b.room.Game().Status() == StatusEnded {
				break
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
		case <-defenseTicker.C:
			if b.room.Game() != nil && b.room.Game().Status() == StatusEnded {
				break
			}
			if b.delayedCastCommand != nil {
				b.botClient.sendCommandToGame("Cast", b.delayedCastCommand)
				b.delayedCastCommand = nil
			}

		case event := <-b.botClient.incomingEvents:
			if b.room.Game() != nil && b.room.Game().Status() == StatusEnded {
				break
			}

			b.dispatchEvent(event)
		}
	}
}

func (b *Bot) dispatchEvent(event interface{}) {
	castEvent, ok := event.(CastEvent)
	if !ok {
		return
	}
	if castEvent.OriginPlayerId == b.botClient.ID() {
		return
	}

	// cast shield spell after some delay

	random := rand.Intn(4)
	var command *CastCommand
	switch random {
	case 0:
		command = &CastCommand{SpellShieldFireball}
	case 1:
		command = &CastCommand{SpellShieldLightning}
	case 2:
		command = &CastCommand{SpellShieldRocks}
	case 3:
		command = &CastCommand{SpellShieldComet}
	}

	b.delayedCastCommand = command
}
