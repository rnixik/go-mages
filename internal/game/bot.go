package game

import "fmt"

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
	fmt.Printf("BOT: got event to make decision: %+v\n", event)
	demoEvent := &DemoEvent{"Demo message value"}
	b.botClient.sendEventToGame(demoEvent)
}
