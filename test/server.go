package main

import (
	"fmt"
	"net"
	"time"

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
	ticker := time.NewTicker(time.Millisecond * 1000)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				world.Timestamp = world.Timestamp + 1
			case <-quit:
				fmt.Println("ticker stopped")
				return
			}
		}
	}()
	world = game.NewWorld(5)
	server.Start()
	fmt.Println("server started")
	world.StartWorld()
	select {}
}

func clientConnect(conn *server.Connection, data []byte) {
	fmt.Println("client connection with:", data)
	world.AddConn(conn)
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
	game.Mult.Write(u)
}
