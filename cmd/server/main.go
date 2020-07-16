package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/yohanesgre/new-server-shooter/pkg/game"
	"github.com/yohanesgre/new-server-shooter/pkg/udpnetwork"
)

var (
	world           *game.World
	mult            *game.Multiplexer
	connectedPlayer int
)

func main() {
	var port = flag.String("port", "10001", "Make new server in a new port")
	var numPlayers = flag.Int("numPlayers", 2, "Insert Connected Players")
	flag.Parse()

	server := udpnetwork.NewServer(":" + *port)

	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket
	server.Start()
	fmt.Println("server started")
	world = game.NewWorld(*numPlayers, false)
	select {}
}

func clientConnect(conn *udpnetwork.Connection, data []byte) {
	// if data[0] != 0 {
	// 	conn.Disconnect([]byte("not allowed"))
	// }
	fmt.Println("client connection with:", data)
	connectedPlayer++
	p := game.NewPlayer(connectedPlayer, "", 0.0, 0.0, 0.0, 0.0, conn)
	h := game.NewPlayerHitBox(p)
	world.List_player.PushBack(p)
	world.List_player_hitbox.PushBack(h)
	conn.SendReliableOrdered([]byte("Welcome"))
	if connectedPlayer == world.Max_player {
		world.StartWorld()
	}
}

func clientDisconnect(conn *udpnetwork.Connection, data []byte) {
	fmt.Println("client disconnect")
	world.DestroyHitboxInListByPlayer(world.FindPlayerInListByConn(conn))
	world.DestroyPlayerInListByConn(conn)
	if world.List_conn.Len() == 0 {
		world.Destroy()
		world = nil
	}
}

func clientTimeout(conn *udpnetwork.Connection, data []byte) {
	fmt.Println("client timeout")
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	return len(data) == 3
}

func handleServerPacket(conn *udpnetwork.Connection, data []byte, channel udpnetwork.Channel) {
	u := game.UnmarshalRequest(data)
	if u.Endpoint == 1 {
		np := u.PayloadToRequestJoin()
		np.Conn = conn
		u.Payload = np
	}
	// fmt.Println("data: ", u)
	game.Mult.Write(u)
}
