package lobby

// RoomInList contains short info about room in the lobby.
type RoomInList struct {
	Id         uint64 `json:"id"`
	OwnerId    uint64 `json:"ownerId"`
	Name       string `json:"name"`
	GameStatus string `json:"gameStatus"`
	MembersNum int    `json:"membersNum"`
}

// ClientInList contains short info about client in the lobby.
type ClientInList struct {
	Id       uint64 `json:"id"`
	Nickname string `json:"nickname"`
}

// ClientJoinedEvent contains info for the just connected client.
type ClientJoinedEvent struct {
	YourId       uint64          `json:"yourId"`
	YourNickname string          `json:"yourNickname"`
	Clients      []*ClientInList `json:"clients"`
	Rooms        []*RoomInList   `json:"roomsCreatedByClients"`
}

// ClientLeftEvent contains id of client who left lobby.
type ClientLeftEvent struct {
	Id uint64 `json:"id"`
}

// ClientBroadCastJoinedEvent contains info for other clients when a new client was connected.
type ClientBroadCastJoinedEvent struct {
	Id       uint64 `json:"id"`
	Nickname string `json:"nickname"`
}

// ClientCreatedRoomEvent contains info of created room.
type ClientCreatedRoomEvent struct {
	Room *RoomInList `json:"room"`
}

// RoomInListRemovedEvent contains id of the room which was removed from lobby
type RoomInListRemovedEvent struct {
	RoomId uint64 `json:"roomId"`
}

// RoomInListUpdatedEvent contains info about room which was changed
type RoomInListUpdatedEvent struct {
	Room *RoomInList `json:"room"`
}
