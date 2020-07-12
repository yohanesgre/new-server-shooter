package main

import (
	"flag"
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
	var port = flag.Int("port", 10001, "Make new server in a new port")
	flag.Parse()

	server := server.NewServer(port)

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
	// if data[0] != 0 {
	// 	conn.Disconnect([]byte("not allowed"))
	// }
	fmt.Println("client connection with:", data)
	if world == nil {
		world = game.NewWorld(5)
		world.AddConn(conn)
		world.StartWorld()
		fmt.Println("World started")
	} else {
		for temp := world.List_conn.Front(); temp != nil; temp = temp.Next() {
			if conn != temp.Value.(*server.Connection) {
				world.AddConn(conn)
			}
		}
	}
}

func clientDisconnect(conn *server.Connection, data []byte) {
	fmt.Println("client disconnect")
	for temp := world.List_conn.Front(); temp != nil; temp = temp.Next() {
		if conn == temp.Value.(*server.Connection) {
			world.List_conn.Remove(temp)
		}
	}
	if world.List_conn.Len() == 0 {
		world.Destroy()
		world = nil
	}
}

func clientTimeout(conn *server.Connection, data []byte) {
	fmt.Println("client timeout")
}

func validateClient(addr *net.UDPAddr, data []byte) bool {
	return len(data) == 3
}

func handleServerPacket(conn *server.Connection, data []byte, channel server.Channel) {
	u := game.UnmarshalRequest(data)
	// if u.Endpoint == 1 {
	// fmt.Println("Data: ", u)
	// }
	game.Mult.Write(u)
}
