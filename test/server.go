package main

import (
	"fmt"
	"net"

	"github.com/yohanesgre/new-server-shooter/game"
	server "github.com/yohanesgre/new-server-shooter/server"
)

var (
	world *game.World
	mult  *game.Multiplexer
)

func main() {
	server := server.NewServer(":10001")

	server.ClientConnect = clientConnect
	server.ClientDisconnect = clientDisconnect
	server.ClientTimeout = clientTimeout
	server.ClientValidation = validateClient
	server.PacketHandler = handleServerPacket
	server.Start()
	fmt.Println("server started")
	select {}
}

func clientConnect(conn *server.Connection, data []byte) {
	fmt.Println("client connection with:", data)
	if world == nil {
		world = game.NewWorld(5)
		world.AddConn(conn)
		world.StartWorld()
		fmt.Println("World started")
	}
	if data[0] != 0 {
		conn.Disconnect([]byte("not allowed"))
	}
}

func clientDisconnect(conn *server.Connection, data []byte) {
	fmt.Println("client disconnect")
}

func clientTimeout(conn *server.Connection, data []byte) {
	fmt.Println("client timeout")
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	return len(data) == 3
}

func handleServerPacket(conn *server.Connection, data []byte, channel server.Channel) {
	u := game.UnmarshalRequest(data)
	if u.Endpoint == 1 {
		fmt.Println("Data: ", u)
	}
	game.Mult.Write(u)
}
