package lobby

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
)

var lastClientId uint64
var lastRoomId uint64

type GameEventsDispatcher interface {
	DispatchGameEvent(client ClientPlayer, event interface{})
	OnClientRemoved(client ClientPlayer)
	OnClientJoined(client ClientPlayer)
	AddBotCommand(client ClientPlayer)
	StartMainLoop()
	Status() string
}

type NewGameFunc func(playersClients []ClientPlayer) GameEventsDispatcher
type NewBotFunc func(botId uint64) ClientPlayer

// Lobby is the first place for connected clients. It passes commands to games.
type Lobby struct {
	// Registered clients.
	clients map[uint64]*Client

	// Outbound events to the clients.
	broadcast chan interface{}

	// Register requests from the clients.
	register chan ClientSender

	// Unregister requests from clients.
	unregister chan ClientSender

	// Commands from clients
	clientCommands chan *ClientCommand

	// Started games
	games []GameEventsDispatcher

	// Rooms created by clients
	rooms map[*Client]*Room

	newGameFunc      NewGameFunc
	newBotFunc       NewBotFunc
	maxPlayersInRoom int
}

func NewLobby(newGameFunc NewGameFunc, newBotFunc NewBotFunc, maxPlayersInRoom int) *Lobby {
	return &Lobby{
		broadcast:        make(chan interface{}),
		register:         make(chan ClientSender),
		unregister:       make(chan ClientSender),
		clients:          make(map[uint64]*Client),
		clientCommands:   make(chan *ClientCommand),
		games:            make([]GameEventsDispatcher, 0),
		rooms:            make(map[*Client]*Room),
		newGameFunc:      newGameFunc,
		newBotFunc:       newBotFunc,
		maxPlayersInRoom: maxPlayersInRoom,
	}
}

func (l *Lobby) Run() {
	log.Println("Go lobby")

	go func() {
		for {
			select {
			case event, ok := <-l.broadcast:
				if !ok {
					continue
				}
				for _, client := range l.clients {
					client.transportClient.SendEvent(event)
				}
			}
		}
	}()

	for {
		select {
		case tc := <-l.register:
			atomic.AddUint64(&lastClientId, 1)
			lastClientIdSafe := atomic.LoadUint64(&lastClientId)
			tc.SetId(lastClientIdSafe)

			client := &Client{
				lobby:           l,
				transportClient: tc,
			}
			l.clients[client.Id()] = client
		case tc := <-l.unregister:
			if client, ok := l.clients[tc.Id()]; ok {
				delete(l.clients, client.Id())
				l.onClientLeft(client)
				client.transportClient.Close()
			}
		case clientCommand := <-l.clientCommands:
			l.onClientCommand(clientCommand)
		}
	}
}

func (l *Lobby) RegisterTransportClient(tc ClientSender) {
	l.register <- tc
}

func (l *Lobby) UnregisterTransportClient(tc ClientSender) {
	l.unregister <- tc
}

func (l *Lobby) HandleClientCommand(tc ClientSender, clientCommand *ClientCommand) {
	if client, ok := l.clients[tc.Id()]; ok {
		clientCommand.client = client
		l.clientCommands <- clientCommand
	}
}

func (l *Lobby) broadcastEvent(event interface{}) {
	l.broadcast <- event
}

func (l *Lobby) onJoinCommand(c *Client, nickname string) {
	c.nickname = nickname

	broadcastEvent := &ClientBroadCastJoinedEvent{
		Id:       c.Id(),
		Nickname: c.Nickname(),
	}
	l.broadcastEvent(broadcastEvent)

	clientsInList := make([]*ClientInList, 0)
	for _, client := range l.clients {
		clientInList := &ClientInList{
			Id:       client.Id(),
			Nickname: client.Nickname(),
		}
		clientsInList = append(clientsInList, clientInList)
	}

	roomsInList := make([]*RoomInList, 0)
	for _, room := range l.rooms {
		roomInList := room.toRoomInList()
		roomsInList = append(roomsInList, roomInList)
	}

	event := &ClientJoinedEvent{
		YourId:       c.Id(),
		YourNickname: c.Nickname(),
		Clients:      clientsInList,
		Rooms:        roomsInList,
	}
	c.transportClient.SendEvent(event)
}

func (l *Lobby) onClientLeft(client *Client) {
	room := client.room
	if room != nil {
		l.onLeftRoom(client, room)
	}
	leftEvent := &ClientLeftEvent{
		Id: client.Id(),
	}
	l.broadcastEvent(leftEvent)
}

func (l *Lobby) onCreateNewRoomCommand(c *Client) {
	_, roomExists := l.rooms[c]
	if roomExists {
		errEvent := &ClientCommandError{errorYouCanCreateOneRoomOnly}
		c.transportClient.SendEvent(errEvent)
		return
	}

	oldRoomJoined := c.room
	if oldRoomJoined != nil {
		l.onLeftRoom(c, oldRoomJoined)
	}

	atomic.AddUint64(&lastRoomId, 1)
	lastRoomIdSafe := atomic.LoadUint64(&lastRoomId)

	room := newRoom(lastRoomIdSafe, c, l)
	l.rooms[c] = room

	event := &ClientCreatedRoomEvent{room.toRoomInList()}
	l.broadcastEvent(event)

	roomJoinedEvent := RoomJoinedEvent{room.toRoomInfo()}
	c.transportClient.SendEvent(roomJoinedEvent)
}

func (l *Lobby) getRoomById(roomId uint64) (room *Room, err error) {
	for _, r := range l.rooms {
		if r.Id() == roomId {
			return r, nil
		}
	}
	return nil, fmt.Errorf("room not found by id = %d", roomId)
}

func (l *Lobby) onLeftRoom(c *Client, room *Room) {
	changedOwner, roomBecameEmpty := room.removeClient(c)
	c.room = nil
	if roomBecameEmpty {
		roomInListRemovedEvent := &RoomInListRemovedEvent{room.Id()}
		l.broadcastEvent(roomInListRemovedEvent)
		l.rooms[c] = nil
		delete(l.rooms, c)
		return
	}
	if changedOwner {
		l.rooms[room.owner.client] = room
		delete(l.rooms, c)
	}
	roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
	l.broadcastEvent(roomInListUpdatedEvent)
}

func (l *Lobby) onJoinRoomCommand(c *Client, roomId uint64) {
	oldRoomJoined := c.room
	if oldRoomJoined != nil && oldRoomJoined.Id() == roomId {
		return
	}
	if oldRoomJoined != nil {
		l.onLeftRoom(c, oldRoomJoined)
	}
	room, err := l.getRoomById(roomId)
	if err == nil {
		room.addClient(c)
		roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
		l.broadcastEvent(roomInListUpdatedEvent)
	} else {
		errEvent := &ClientCommandError{errorRoomDoesNotExist}
		c.transportClient.SendEvent(errEvent)
	}
}

func (l *Lobby) onClientCommand(cc *ClientCommand) {
	if cc.Type == ClientCommandTypeLobby {
		if cc.SubType == ClientCommandLobbySubTypeJoin {
			var nickname string
			if err := json.Unmarshal(cc.Data, &nickname); err != nil {
				return
			}
			l.onJoinCommand(cc.client, nickname)
		} else if cc.SubType == ClientCommandLobbySubTypeCreateRoom {
			l.onCreateNewRoomCommand(cc.client)
		} else if cc.SubType == ClientCommandLobbySubTypeJoinRoom {
			var roomId uint64
			if err := json.Unmarshal(cc.Data, &roomId); err != nil {
				return
			}
			l.onJoinRoomCommand(cc.client, roomId)
		}
	} else if cc.Type == ClientCommandTypeGame {

		l.dispatchGameEvent(cc)
	} else if cc.Type == ClientCommandTypeRoom {
		if cc.client.room == nil {
			return
		}
		cc.client.room.onClientCommand(cc)
	}
}

func (l *Lobby) dispatchGameEvent(cc *ClientCommand) {
	if cc.client.room == nil {
		return
	}
	if cc.client.room.game == nil {
		return
	}
	cc.client.room.game.DispatchGameEvent(cc.client, cc.Data)
}

func (l *Lobby) sendRoomUpdate(room *Room) {
	roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
	l.broadcastEvent(roomInListUpdatedEvent)
}
