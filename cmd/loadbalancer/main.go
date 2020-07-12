package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/yohanesgre/new-server-shooter/pkg/loadbalancer"
)

func main() {
	// setLimit()
	fmt.Println("Listening at 9000")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}
		defer conn.Close()
		fmt.Println("Welcome")
		conn.Write([]byte("Welcome"))

		go handleConn(conn)
		connections = append(connections, conn)
	}
}

func handleConn(conn net.Conn) {
	var buf [128]byte

	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Println(err)
			break
		}
		// fmt.Printf("received %x\n", n)

		if len(buf) > 0 {
			fmt.Println("received not 0")
			packet := loadbalancer.UnmarshalRequest(buf[:n])
			switch packet.Endpoint {
			case int(loadbalancer.GET_LOBBIES):
				fmt.Printf("Client %v request to get Lobbies Index\n", conn.RemoteAddr)
				loadbalancer.IndexLobbies(conn)
			case int(loadbalancer.CREATE_LOBBY):
				fmt.Printf("Client %v request to create Lobby\n", conn.RemoteAddr)
				loadbalancer.CreateLobby(conn, packet.Payload[1])
			case int(loadbalancer.CONNECT_LOBBY):
				fmt.Printf("Client %v request to connect to Lobby %d\n", conn.RemoteAddr, packet.Payload)
				i, err := strconv.Atoi(packet.Payload[0])
				if err != nil {
					log.Println(err)
				}
				loadbalancer.ConnectLobby(i, packet.Payload[1], conn)
			case int(loadbalancer.START_GAME):
				fmt.Println("Start Server")
				i, err := strconv.Atoi(packet.Payload[0])
				if err != nil {
					log.Println(err)
				}
				loadbalancer.StartGame(i)
			default:
				conn.Write([]byte{0, 1, 2})
			}
		}
	}
}
