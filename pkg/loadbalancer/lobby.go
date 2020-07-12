package loadbalancer

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
)

type Lobby struct {
	Id         int        `msgpack:"LobbyId"`
	MaxPlayers int        `msgpack:"MaxPlayers"`
	NumPlayers int        `msgpack:"NumPlayers"`
	Name       []string   `msgpack:"Name,asArray`
	IsStart    bool       `msgpack:"IsStart"`
	Host       string     `msgpack:"Host"`
	Port       int        `msgpack:"Port"`
	Conn       []net.Conn `msgpack:",omitempty`
}

const (
	MAX_PLAYER   int    = 5
	SERVER_HOST  string = "yohanesgre.tech"
	INITIAL_PORT int    = 10000
)

var ch chan int = make(chan int)

type Lobbies []*Lobby

var (
	currentId    int
	lobbies      Lobbies         = make(Lobbies, 0, 10)
	lobsResponse []LobbyResponse = make([]LobbyResponse, 0, 10)
)

// func ConvertLobbiesToResponse() []LobbyResponse {
// 	result := make([]LobbyResponse, 0, 10)
// 	for _, l := range lobbies {
// 		go func() {
// 			l.NumPlayers <- l.NumPlayers
// 		}()
// 		n := <-l.NumPlayers
// 		result = append(lobsResponse, LobbyResponse{l.Id, l.MaxPlayers, n, l.IsStart, l.Host, l.Port})
// 	}
// 	return result
// }

func IndexLobbies(conn net.Conn) {
	fmt.Println("Index Lobby")
	r := Response{int(INDEX_LOBBY), MakeLobbiesResponse()}
	resp := r.MarshalResponse()
	_, err := conn.Write(resp[:])
	if err != nil {
		log.Println("Error: ", err)
	}
	fmt.Println("Index Lobby Send")
}

func CreateLobby(conn net.Conn, name string) {
	fmt.Println("Create Lobby")
	currentId += 1
	port := INITIAL_PORT + currentId
	l := &Lobby{currentId, MAX_PLAYER, 0, make([]string, 0, MAX_PLAYER), false, SERVER_HOST, port, make([]net.Conn, 0, MAX_PLAYER)}
	lobbies = append(lobbies, l)
	fmt.Println("create")
	ConnectLobby(currentId, name, conn)
}

func DeleteLobby(id int) error {
	for i, lobby := range lobbies {
		if lobby.Id == id {
			lobbies = append(lobbies[:i], lobbies[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Lobby with id of %d to delete", id)
}

func DisconnectLobby(id int, conn net.Conn, m sync.Mutex) error {
	for _, lobby := range lobbies {
		if lobby.Id == id {
			go func() {
				lobby.NumPlayers--
				ch <- 0
			}()
			<-ch
			DeleteConnMemberInLobby(lobby.Conn, conn)
			for _, c := range lobby.Conn {
				IndexLobbies(c)
			}
		}
	}
	return fmt.Errorf("Could not find Lobby with id of %d to delete", id)
}

func ConnectLobby(id int, name string, conn net.Conn) error {
	fmt.Println("Connect Lobby")
	var wg sync.WaitGroup
	l := FindLobbyById(id)
	fmt.Println("selected lobby: ", l)
	if l.NumPlayers < l.MaxPlayers {
		wg.Add(1)
		go func() {
			l.NumPlayers++
			l.Name = append(l.Name, name)
			l.Conn = append(l.Conn, conn)
			wg.Done()
		}()
	}
	wg.Wait()
	if l.NumPlayers > 1 {
		for _, c := range l.Conn {
			r := Response{int(PLAYER_CONNECTED), MakeLobbiesResponse()}
			c.Write(r.MarshalResponse())
		}
	} else {
		r := Response{int(PLAYER_CONNECTED), MakeLobbiesResponse()}
		conn.Write(r.MarshalResponse())
	}

	return nil
	// r := Response{int(PLAYER_CONNECTED), lobbies}
	// 			conn.Write(r.MarshalResponse())
	// 		} else if lobby.NumPlayers == lobby.MaxPlayers {
	// 			StartServer(lobby.Port, lobby.Id)
	// 			for _, conn := range lobby.Conn {
	// 				r := Response{int(SERVER_START), lobbies}
	// 				conn.Write(r.MarshalResponse())
	// 			}
	// 		}
	// return fmt.Errorf("Could not find Lobby with id of %d to delete", id)
}

func FindLobbyByConnection(conn net.Conn) *Lobby {
	var result *Lobby
	for i, lobby := range lobbies {
		if lobby.Conn[i] == conn {
			result = lobby
		}
	}
	return result
}

func FindLobbyById(id int) *Lobby {
	var result *Lobby
	for _, lobby := range lobbies {
		if lobby.Id == id {
			result = lobby
		}
	}
	return result
}

func MakeLobbiesResponse() []Lobby {
	var resp []Lobby
	if len(lobbies) == 0 {
		resp = make([]Lobby, 0, 1)
	} else {
		for _, l := range lobbies {
			resp = append(resp, *l)
		}
	}
	return resp
}

func DeleteConnMemberInLobby(conns []net.Conn, conn net.Conn) error {
	for i, c := range conns {
		if c == conn {
			conns = append(conns[:i], conns[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Lobby with id of %v to delete", conn)
}

func StartGame(id int) {
	l := FindLobbyById(id)
	fmt.Println("selected lobby: ", l)
	if l.NumPlayers > 1 {
		StartServer(l.Port, l.Id)
		for _, c := range l.Conn {
			r := Response{int(SERVER_START), MakeLobbiesResponse()}
			c.Write(r.MarshalResponse())
		}
	}
}

func StartServer(port, id int) {
	go func() {
		cmd := exec.Command("./server", "-port="+strconv.Itoa(port))
		fmt.Println("Start Server Game")
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case err := <-done:
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					log.Printf("Exit Status: %d", status.ExitStatus())
					DeleteLobby(id)
				}
			} else {
				// log.Fatalf("cmd.Wait: %v", err)
				log.Println("Server Closed")
				DeleteLobby(id)
			}
		}
	}()
}
