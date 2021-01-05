package lobby

// ClientSender represents interface which sends events to connected players.
type ClientSender interface {
	SendEvent(event interface{})
	Id() uint64
	SetId(id uint64)
	Close()
}

type ClientPlayer interface {
	SendEvent(event interface{})
	Id() uint64
	Nickname() string
}

type Client struct {
	lobby *Lobby

	transportClient ClientSender

	nickname string
	room     *Room
}

func (c *Client) SendEvent(event interface{}) {
	c.transportClient.SendEvent(event)
}

func (c *Client) Id() uint64 {
	return c.transportClient.Id()
}

func (c *Client) Nickname() string {
	return c.nickname
}
