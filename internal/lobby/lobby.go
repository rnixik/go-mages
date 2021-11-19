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
	DispatchGameCommand(client ClientPlayer, eventName string, eventData interface{})
	OnClientRemoved(client ClientPlayer)
	OnClientJoined(client ClientPlayer)
	StartMainLoop()
	Status() string
}

type NewGameFunc func(playersClients []ClientPlayer, broadcastEventFunc func(event interface{})) GameEventsDispatcher
type NewBotFunc func(botId uint64, sendGameEvent func(client ClientPlayer, eventName string, eventData json.RawMessage)) ClientPlayer

type MatchMaker interface {
	MakeMatch(client ClientPlayer, foundFunc func(clients []ClientPlayer), notFoundFunc func(), addBotFunc func() ClientPlayer)
	Cancel(client ClientPlayer)
}

// Lobby is the first place for connected clients. It passes commands to games.
type Lobby struct {
	// Registered clients.
	clients map[uint64]ClientPlayer

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
	roomsCreatedByClients map[ClientPlayer]*Room

	// Room where client is
	clientsJoinedRooms map[ClientPlayer]*Room

	newGameFunc      NewGameFunc
	newBotFunc       NewBotFunc
	matchMaker       MatchMaker
	maxPlayersInRoom int
}

func NewLobby(newGameFunc NewGameFunc, newBotFunc NewBotFunc, matchMaker MatchMaker, maxPlayersInRoom int) *Lobby {
	return &Lobby{
		broadcast:             make(chan interface{}),
		register:              make(chan ClientSender),
		unregister:            make(chan ClientSender),
		clients:               make(map[uint64]ClientPlayer),
		clientCommands:        make(chan *ClientCommand),
		games:                 make([]GameEventsDispatcher, 0),
		roomsCreatedByClients: make(map[ClientPlayer]*Room),
		clientsJoinedRooms:    make(map[ClientPlayer]*Room),
		newGameFunc:           newGameFunc,
		newBotFunc:            newBotFunc,
		matchMaker:            matchMaker,
		maxPlayersInRoom:      maxPlayersInRoom,
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
					client.SendEvent(event)
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
				client.CloseConnection()
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

func (l *Lobby) joinLobbyCommand(c ClientPlayer, nickname string) {
	c.SetNickname(nickname)

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
	for _, room := range l.roomsCreatedByClients {
		roomInList := room.toRoomInList()
		roomsInList = append(roomsInList, roomInList)
	}

	event := &ClientJoinedEvent{
		YourId:       c.Id(),
		YourNickname: c.Nickname(),
		Clients:      clientsInList,
		Rooms:        roomsInList,
	}
	c.SendEvent(event)
}

func (l *Lobby) onClientLeft(client ClientPlayer) {
	room := l.clientsJoinedRooms[client]
	if room != nil {
		l.onLeftRoom(client, room)
	}
	leftEvent := &ClientLeftEvent{
		Id: client.Id(),
	}
	l.broadcastEvent(leftEvent)
}

func (l *Lobby) createNewRoomCommand(c ClientPlayer) {
	_, roomExists := l.roomsCreatedByClients[c]
	if roomExists {
		errEvent := &ClientCommandError{errorYouCanCreateOneRoomOnly}
		c.SendEvent(errEvent)
		return
	}

	oldRoomJoined := l.clientsJoinedRooms[c]
	if oldRoomJoined != nil {
		l.onLeftRoom(c, oldRoomJoined)
	}

	atomic.AddUint64(&lastRoomId, 1)
	lastRoomIdSafe := atomic.LoadUint64(&lastRoomId)

	room := newRoom(lastRoomIdSafe, c, l)
	l.roomsCreatedByClients[c] = room

	event := &ClientCreatedRoomEvent{room.toRoomInList()}
	l.broadcastEvent(event)

	roomJoinedEvent := &RoomJoinedEvent{room.toRoomInfo()}
	c.SendEvent(roomJoinedEvent)
}

func (l *Lobby) getRoomById(roomId uint64) (room *Room, err error) {
	for _, r := range l.roomsCreatedByClients {
		if r.Id() == roomId {
			return r, nil
		}
	}
	return nil, fmt.Errorf("room not found by id = %d", roomId)
}

func (l *Lobby) onLeftRoom(c ClientPlayer, room *Room) {
	changedOwner, roomBecameEmpty := room.removeClient(c)
	delete(l.clientsJoinedRooms, c)
	if roomBecameEmpty {
		roomInListRemovedEvent := &RoomInListRemovedEvent{room.Id()}
		l.broadcastEvent(roomInListRemovedEvent)
		l.roomsCreatedByClients[c] = nil
		delete(l.roomsCreatedByClients, c)
		return
	}
	if changedOwner {
		l.roomsCreatedByClients[room.owner.client] = room
		delete(l.roomsCreatedByClients, c)
	}
	roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
	l.broadcastEvent(roomInListUpdatedEvent)
}

func (l *Lobby) joinRoomCommand(c ClientPlayer, roomId uint64) {
	oldRoomJoined := l.clientsJoinedRooms[c]
	if oldRoomJoined != nil && oldRoomJoined.Id() == roomId {
		return
	}
	if oldRoomJoined != nil {
		l.onLeftRoom(c, oldRoomJoined)
	}
	room, err := l.getRoomById(roomId)
	if err == nil {
		l.clientsJoinedRooms[c] = room
		room.addClient(c)
		roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
		l.broadcastEvent(roomInListUpdatedEvent)
	} else {
		errEvent := &ClientCommandError{errorRoomDoesNotExist}
		c.SendEvent(errEvent)
	}
}

func (l *Lobby) makeMatch(c ClientPlayer) {
	oldRoomJoined := l.clientsJoinedRooms[c]
	if oldRoomJoined != nil {
		l.onLeftRoom(c, oldRoomJoined)
	}
	l.createNewRoomCommand(c)
	room := l.roomsCreatedByClients[c]
	l.matchMaker.MakeMatch(
		c,
		func(clients []ClientPlayer) {
			room.onStartGameCommand(c)
		},
		func() {},
		func() ClientPlayer {
			return room.createBot()
		},
	)
}

func (l *Lobby) onClientCommand(cc *ClientCommand) {
	if cc.Type == ClientCommandTypeLobby {
		if cc.SubType == ClientCommandLobbySubTypeJoin {
			var nickname string
			if err := json.Unmarshal(cc.Data, &nickname); err != nil {
				return
			}
			l.joinLobbyCommand(cc.client, nickname)
		} else if cc.SubType == ClientCommandLobbySubTypeCreateRoom {
			l.createNewRoomCommand(cc.client)
		} else if cc.SubType == ClientCommandLobbySubTypeJoinRoom {
			var roomId uint64
			if err := json.Unmarshal(cc.Data, &roomId); err != nil {
				return
			}
			l.joinRoomCommand(cc.client, roomId)
		} else if cc.SubType == ClientCommandLobbySubTypeMakeMatch {
			l.makeMatch(cc.client)
		}
	} else if cc.Type == ClientCommandTypeRoom {
		if l.clientsJoinedRooms[cc.client] == nil {
			return
		}
		l.clientsJoinedRooms[cc.client].onClientCommand(cc)
	} else if cc.Type == ClientCommandTypeGame {
		l.dispatchGameCommand(cc)
	}
}

func (l *Lobby) dispatchGameCommand(cc *ClientCommand) {
	if l.clientsJoinedRooms[cc.client] == nil {
		return
	}
	if l.clientsJoinedRooms[cc.client].game == nil {
		return
	}
	l.clientsJoinedRooms[cc.client].game.DispatchGameCommand(cc.client, cc.SubType, cc.Data)
}

func (l *Lobby) sendRoomUpdate(room *Room) {
	roomInListUpdatedEvent := &RoomInListUpdatedEvent{room.toRoomInList()}
	l.broadcastEvent(roomInListUpdatedEvent)
}
