package main

import (
	"fmt"
	"time"

	"github.com/yohanesgre/new-server-shooter/game"
	server "github.com/yohanesgre/new-server-shooter/server"
)

var (
	snapshots []game.Snapshot
	id        int
	player    *game.Player
)

func main() {
	client := server.NewClient("127.0.0.1:10001")

	client.ServerConnect = serverConnect
	client.ServerDisconnect = serverDisconnect
	client.ServerTimeout = serverTimeout
	client.PacketHandler = handleClientPacket
	time.Sleep(3000 * time.Millisecond)
	client.ConnectWithData([]byte{0, 1, 2})
	snapshots = make([]game.Snapshot, 0, 1000000)
	player = game.NewPlayer(0, "Namaaa", 0, 0, 0, 10.0)
	select {}
}

func serverConnect(conn *server.Connection, data []byte) {
	fmt.Println("connected to server")
	r := makeRequestJoin()
	m := r.MarshalRequest()
	conn.SendReliable(m)
}

func serverDisconnect(conn *server.Connection, data []byte) {
	fmt.Println("disconnected from server:", string(data))
}

func serverTimeout(conn *server.Connection, data []byte) {
	fmt.Println("server timeout")
}

func handleClientPacket(conn *server.Connection, data []byte, channel server.Channel) {
	s := game.UnmarshalSnapshot(data)
	fmt.Printf("Id: %d \nSnapshot: %v\n Channel: %+v\n", id, s, channel)
	snapshots = append(snapshots, s)
	fmt.Println("Before If Channel 2")
	if channel == 2 {
		//player = &s.Players[len(s.Players)-1]
		go worker(conn, id)
	}
	fmt.Println("After If Channel 2")
}

func makeRequestJoin() *game.Request {
	r := game.NewRequest(game.JOIN, player)
	_m := r.MarshalRequest()
	_u := game.UnmarshalRequest(_m)
	return &_u
}

func worker(conn *server.Connection, id int) {
	for i := 0; i < 50; i++ {
		//newP := game.NewPlayer(player.Id, player.Name, player.Pos_x+1.0, player.Pos_y+1.0, player.Rotation+1.0, player.FOV)
		// r := game.NewRequest(game.MOVE, newP)
		r := game.NewRequest(game.MOVE, nil)
		_m := r.MarshalRequest()
		conn.SendUnreliableOrdered(_m)
		time.Sleep(200 * time.Millisecond)
	}
	for i := 0; i < 3; i++ {
		// newP := game.NewPlayer(player.Id, player.Name, player.Pos_x+1.0, player.Pos_y+1.0, player.Rotation+1.0, player.FOV)
		// r := game.NewRequest(game.SHOOT, newP)
		r := game.NewRequest(game.MOVE, nil)
		_m := r.MarshalRequest()
		conn.SendUnreliableOrdered(_m)
		time.Sleep(200 * time.Millisecond)
	}
}
