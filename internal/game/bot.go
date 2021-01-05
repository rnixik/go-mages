package game

type Bot struct {
	transportClient ClientPlayer
	nickname        string
}

func NewBot(botId uint64) *Bot {
	return &Bot{
		transportClient: nil,
		nickname:        "",
	}
}

func (b *Bot) SendEvent(event interface{}) {
	b.transportClient.SendEvent(event)
}

func (b *Bot) Id() uint64 {
	return b.transportClient.Id()
}

func (b *Bot) Nickname() string {
	return b.nickname
}
