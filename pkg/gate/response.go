package gate

import (
	"github.com/vmihailenco/msgpack/v5"
)

type EndpointResponse int

const (
	INDEX_LOBBY      EndpointResponse = 1
	PLAYER_CONNECTED EndpointResponse = 2
	SERVER_START     EndpointResponse = 3
)

type LobbyResponse struct {
	Id         int    `msgpack:"LobbyId"`
	MaxPlayers int    `msgpack:"MaxPlayers"`
	NumPlayers int    `msgpack:"NumPlayers"`
	IsStart    bool   `msgpack:"IsStart"`
	Host       string `msgpack:"Host"`
	Port       int    `msgpack:"Port"`
}

type Response struct {
	Endpoint int     `msgpack:"Endpoint"`
	Lobby    []Lobby `msgpack:"Lobby,asArray"`
}

func (r *Response) MarshalResponse() []byte {
	b, err := msgpack.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}
