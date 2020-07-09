package game

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/erinpentecost/gloop"
	"github.com/yohanesgre/new-server-shooter/server"
)

var (
	Mult        *Multiplexer
	seq_counter int32
)

const Hz30Delay time.Duration = time.Duration(int64(time.Second) / 30)

type World struct {
	id                 int
	max_player         int
	list_conn          list.List
	list_player_hitbox list.List
	list_player        list.List
	list_bullet        list.List
	list_weapon        list.List
	game_loop          *gloop.Loop
	Timestamp          float32
	currTime           int64
	initTime           int64
}

func NewWorld(max_player int) *World {
	w := new(World)
	w.max_player = max_player
	w.list_player.Init()
	w.list_bullet.Init()
	w.list_weapon.Init()
	SeedSpawnerPlayer(max_player)
	Mult = NewMultiplexer(50)
	return w
}

func (w *World) Destroy() {
	w.Destroy()
}

//Request Handler
//Handle all request sent by Players
func (w *World) RequestHandler(_r Request) {
	switch _r.Endpoint {
	case JOIN:
		p := _r.PayloadToPlayer()
		for _, spawner := range list_spawner_player {
			if !spawner.Filled {
				p_ := NewPlayer(w.list_player.Len()+1, p.Name, spawner.Pos_x, spawner.Pos_y, 0.0, p.FOV)
				h_ := NewPlayerHitBox(p_, 10, 10)
				w.list_player.PushBack(p_)
				w.list_player_hitbox.PushBack(h_)
				w.list_conn.Back().Value.(*server.Connection).SendReliableOrdered(w.GenerateSnapshot(seq_counter))
				fmt.Printf("Player Joined: %v\n", p.Name)
				break
			}
		}
	case MOVE:
		p := _r.PayloadToPlayer()
		p_ := w.FindPlayerInList(p)
		h_ := w.FindPlayerHitBoxInList(p_)
		p_.UpdatePlayer(p)
		h_.UpdatePlayerHitBox(p_)

	case SHOOT:
		p := _r.PayloadToPlayer()
		p_ := w.FindPlayerInList(p)
		p_.UpdatePlayer(p)
		w.SpawnBullet(p_)
	}
}

//Find Player Pointer in List Of Players' Pointer
func (w *World) FindPlayerInList(_player *Player) *Player {
	var result *Player
	for temp := w.list_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		if _p.Id == _player.Id {
			result = _p
		}
	}
	return result
}

//Find Player Pointer in List Of Players' Pointer
func (w *World) FindPlayerInListById(id int) *Player {
	var result *Player
	for temp := w.list_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		if _p.Id == id {
			result = _p
		}
	}
	return result
}

//Find Bullet Pointer in List Of Bullets' Pointer
func (w *World) FindBulletInList(_bullet *Bullet) *Bullet {
	var result *Bullet
	for temp := w.list_bullet.Front(); temp != nil; temp = temp.Next() {
		_b := temp.Value.(*Bullet)
		if _b.Id == _bullet.Id {
			result = _b
		}
	}
	return result
}

//Find WeaponDrop Pointer in List Of WeaponDrop' Pointer
func (w *World) FindWeaponDropInList(_weapon *WeaponDrop) *WeaponDrop {
	var result *WeaponDrop
	for temp := w.list_weapon.Front(); temp != nil; temp = temp.Next() {
		_w := temp.Value.(*WeaponDrop)
		if _w.Id == _weapon.Id {
			result = _w
		}
	}
	return result
}

//Find PlayerHitBox Pointer in List Of PlayerHitBox' Pointer with Player
func (w *World) FindPlayerHitBoxInList(_player *Player) *PlayerHitBox {
	var result *PlayerHitBox
	for temp := w.list_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		_h := temp.Value.(*PlayerHitBox)
		if _h.Id == _player.Id {
			result = _h
		}
	}
	return result
}

func (w *World) ListPlayerToArray() []Player {
	result := make([]Player, 0, w.max_player)
	for temp := w.list_player.Front(); temp != nil; temp = temp.Next() {
		result = append(result, *temp.Value.(*Player))
	}
	return result
}

func (w *World) ListBulletToArray() []Bullet {
	result := make([]Bullet, 0, w.list_bullet.Len())
	for temp := w.list_bullet.Front(); temp != nil; temp = temp.Next() {
		result = append(result, *temp.Value.(*Bullet))
	}
	return result
}

func (w *World) SpawnBullet(_player *Player) {
	_w := FindWeaponType(_player.WeaponOwned)
	b := NewBullet(w.list_bullet.Len()+1, _player.Id, _w.Bullet_id, _player.Pos_x, _player.Pos_y, _player.Rotation)
	// b := NewBullet(123, 412, 3112, 312, 23132, 23123)
	//fmt.Printf("Bullet: %v\n", b)
	w.list_bullet.PushBack(b)
	fmt.Println("Spawn")
}

func (w *World) DestroyBullet(_bullet *Bullet) {
	/*for temp := w.list_bullet.Front(); temp != nil; temp = temp.Next() {
		if temp.Value.(*Bullet) == _bullet {
			w.list_bullet.Remove(temp).(*Bullet).Destroy()
		}
	}*/
	fmt.Println("Bullet Destroyed: ", _bullet.Id)
}

func (w *World) StartWorld() {
	var wg sync.WaitGroup
	w.initTime = MakeTimestamp()
	loop, _ := gloop.NewLoop(nil, nil, Hz30Delay, Hz30Delay)
	w.game_loop = loop
	render := func(step time.Duration) error {
		return nil
	}
	simulate := func(step time.Duration) error {
		seq_counter = seq_counter + 1
		if w.list_bullet.Len() != 0 {
			for tempBullet := w.list_bullet.Front(); tempBullet != nil; tempBullet = tempBullet.Next() {
				wg.Add(1)
				go tempBullet.Value.(*Bullet).MoveBullet(&wg)
			}
			wg.Wait()
			for tempBullet := w.list_bullet.Front(); tempBullet != nil; tempBullet = tempBullet.Next() {
				bul := tempBullet.Value.(*Bullet)
				go func() {
					if bul.Distance > 100.0 {
						//fmt.Println("Bullet Length: ", bul.Id, " | ", int(bul.Distance))
						w.DestroyBullet(bul)
					}
				}()
			}
			for tempHitBox := w.list_player_hitbox.Front(); tempHitBox != nil; tempHitBox = tempHitBox.Next() {
				hitbox := tempHitBox.Value.(*PlayerHitBox)
				wg.Add(1)
				go func() {
					player := w.FindPlayerInListById(hitbox.Id)
					hit, dmg, bul := hitbox.CheckCollision(w.list_bullet)
					if hit {
						player.HitPlayer(dmg)
						w.DestroyBullet(bul)
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
		if Mult.GetStream().GetLen() != 0 {
			_ar := Mult.ReadAll()
			for _, r := range *_ar {
				w.RequestHandler(r)
			}
		}
		w.currTime = MakeTimestamp()
		w.Timestamp = float32(w.currTime-w.initTime) / 1000
		for temp := w.list_conn.Front(); temp != nil; temp = temp.Next() {
			c := temp.Value.(*server.Connection)
			s := w.GenerateSnapshot(seq_counter)
			c.SendUnreliableOrdered(s)
		}
		return nil
	}
	w.game_loop.Render = render
	w.game_loop.Simulate = simulate
	w.game_loop.Start()
	fmt.Println("Game World Start")
}

func (w *World) StopWorld() {
	fmt.Print("--------------------------------------------\n")
	for temp := w.list_player.Front(); temp != nil; temp = temp.Next() {
		fmt.Printf("%v\n", temp.Value.(*Player))
	}
	fmt.Print("--------------------------------------------\n")
	fmt.Println("Game World Stop")
	w.game_loop.Done()
}

func (w *World) AddConn(conn *server.Connection) {
	w.list_conn.PushBack(conn)
}

func (w *World) GenerateSnapshot(seq int32) []byte {
	n := NewSnapshot(seq, w.Timestamp, w.ListPlayerToArray(), w.ListBulletToArray())
	b := n.MarshalSnapshot()
	fmt.Println("Snapshot: ", n)
	return b
}
