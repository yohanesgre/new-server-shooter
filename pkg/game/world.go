package game

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/erinpentecost/gloop"
	"github.com/yohanesgre/new-server-shooter/pkg/udpnetwork"
)

var (
	Mult        *Multiplexer
	seq_counter int32
	mutex       = &sync.Mutex{}
)

const Hz30Delay time.Duration = time.Duration(int64(time.Second) / 60)
const Hz200Delay time.Duration = time.Duration(int64(time.Second) / 120)

type World struct {
	id                   int
	Max_player           int
	List_conn            list.List
	List_player_hitbox   list.List
	List_player          list.List
	list_bullet          list.List
	list_weapon          list.List
	list_action_shoot    list.List
	game_loop            *gloop.Loop
	start_game           bool
	Timestamp            float32
	currTime30Hz         int64
	initTime             int64
	deltaTime100Hz       float64
	currTime100Hz        int64
	lastTime100Hz        int64
	action_shoot_counter int
	isNetworkBindCulling bool
}

func NewWorld(Max_player int, culling bool) *World {
	w := new(World)
	w.Max_player = Max_player
	w.List_player.Init()
	w.list_bullet.Init()
	w.list_weapon.Init()
	w.list_action_shoot.Init()
	w.List_conn.Init()
	SeedSpawnerPlayer(Max_player)
	Mult = NewMultiplexer(120)
	w.action_shoot_counter = 0
	w.isNetworkBindCulling = culling
	return w
}

func (w *World) Destroy() {
	w = nil
}

//Request Handler
//Handle all request sent by Players
func (w *World) RequestHandler(_r Request) {
	// if w.List_player.Len() != 0 {
	// 	w.ForceAllPlayersIdle()
	// }
	// if w.list_action_shoot.Len() != 0 {
	// 	fmt.Println("List Actions: ", w.ListActionShootToArray())
	// }
	switch _r.Endpoint {
	case JOIN:
		fmt.Println("Len: ", w.List_player.Len())
		for i := 0; i < len(list_spawner_player); i++ {
			spawner := list_spawner_player[i]
			if !spawner.Filled {
				spawner.Filled = true
				r := _r.Payload.(*RequestJoin)
				p := w.FindPlayerInListByConn(r.Conn)
				// fmt.Println("Player Before: ", p)
				p.Name = r.Name
				p.FOV = r.FOV
				p.Pos_x = spawner.Pos_x
				p.Pos_y = spawner.Pos_y
				// fmt.Println("Player After: ", p)
				h := w.FindPlayerHitBoxInList(p)
				// fmt.Println("Hitbox Before: ", h)
				h.UpdatePlayerHitBox(p)
				// fmt.Println("Hitbox After: ", h)
				break
			}
		}
		c_name := 0
		for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
			c := temp.Value.(*Player)
			if c.Name != "" {
				c_name++
			}
		}
		if c_name == w.Max_player {
			snap := w.GenerateSnapshotReliable(seq_counter)
			for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
				c := temp.Value.(*Player)
				c.Conn.SendReliableOrdered(snap)
			}

			w.start_game = true
		}
	case MOVE:
		r := _r.PayloadToRequestMove()
		p := w.FindPlayerInListById(r.Id)
		h := w.FindPlayerHitBoxInList(p)
		// mutex.Lock()
		p.Move(r.Direction, r.Rotation, w.deltaTime100Hz)
		h.UpdatePlayerHitBox(p)
		// mutex.Unlock()
	case SHOOT:
		if w.action_shoot_counter >= 125 {
			w.action_shoot_counter = 0
		}
		w.action_shoot_counter = w.action_shoot_counter + 1
		r := _r.PayloadToRequestShoot()
		p := w.FindPlayerInListById(r.Id)
		p.UpdateState(Shooting)
		w.list_action_shoot.PushBack(NewActionShootResponse(w.action_shoot_counter, r.Id, r.Pos_x, r.Pos_y))
	case SHOOT_DONE:
		r := _r.PayloadToRequestShootDone()
		for temp := w.list_action_shoot.Front(); temp != nil; temp = temp.Next() {
			a := temp.Value.(*ActionShootResponse)
			if a.Id == r.Id {
				w.list_action_shoot.Remove(temp)
			}
		}
	case COLIDED:
		r := _r.PayloadToRequestBulletColided()
		// fmt.Printf("payload: %#v\n", r)
		// fmt.Printf("hitted id: %#v\n", r.HittedId)
		hitbox := w.FindPlayerHitBoxInListByHittedId(r.HittedId)
		// fmt.Printf("hitbox: %#v\n", hitbox)
		if hitbox.CheckHit(r.HittedId, r.Pos_x, r.Pos_y) {
			dmg := FindBulletType(FindWeaponType(r.WeaponId).Bullet_id).Damage
			// fmt.Println("Dmg: ", dmg)
			w.FindPlayerInListById(r.HittedId).HitPlayer(dmg)
			// fmt.Println("Hp: ", w.FindPlayerInListById(r.HittedId).Hp)
		}
	}
}

func (w *World) ForceAllPlayersIdle() {
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		_p.UpdateState(Idling)
	}
}

//Find Player Pointer in List Of Players' Pointer
func (w *World) FindPlayerInList(_player *Player) *Player {
	var result *Player
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
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
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		if _p.Id == id {
			result = _p
		}
	}
	return result
}

func (w *World) FindPlayerInListByConn(conn *udpnetwork.Connection) *Player {
	var result *Player
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		if _p.Conn == conn {
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
	for temp := w.List_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		_h := temp.Value.(*PlayerHitBox)
		if _h.Id == _player.Id {
			result = _h
		}
	}
	return result
}

//Find PlayerHitBox Pointer in List Of PlayerHitBox' Pointer with HittedId
func (w *World) FindPlayerHitBoxInListByHittedId(_hittedId int) *PlayerHitBox {
	var result *PlayerHitBox
	for temp := w.List_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		_h := temp.Value.(*PlayerHitBox)
		if _h.Id == _hittedId {
			result = _h
		}
	}
	return result
}

func (w *World) DestroyHitboxInListByPlayer(player *Player) {
	for temp := w.List_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*PlayerHitBox)
		if _p.Id == player.Id {
			w.List_player_hitbox.Remove(temp)
		}
	}
}

func (w *World) DestroyPlayerInListByConn(conn *udpnetwork.Connection) {
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		if _p.Conn == conn {
			w.List_player.Remove(temp)
		}
	}
}

func (w *World) ListPlayerToArray() []Player {
	result := make([]Player, 0, w.Max_player)
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
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

func (w *World) ListHitboxToArray() []PlayerHitBox {
	result := make([]PlayerHitBox, 0, w.List_player_hitbox.Len())
	for temp := w.List_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		result = append(result, *temp.Value.(*PlayerHitBox))
	}
	return result
}

func (w *World) ListActionShootToArray() []ActionShootResponse {
	result := make([]ActionShootResponse, 0, w.list_action_shoot.Len())
	for temp := w.list_action_shoot.Front(); temp != nil; temp = temp.Next() {
		result = append(result, *temp.Value.(*ActionShootResponse))
	}
	return result
}

func (w *World) ListConnToArray() []*udpnetwork.Connection {
	result := make([]*udpnetwork.Connection, 0, w.List_conn.Len())
	for temp := w.List_conn.Front(); temp != nil; temp = temp.Next() {
		result = append(result, temp.Value.(*udpnetwork.Connection))
	}
	return result
}

// func (w *World) ListActionShootToArray() []ActionShootResponse {
// 	result := make([]ActionShootResponse, 0, w.list_action_shoot.Len())
// 	for temp := w.list_action_shoot.Front(); temp != nil; temp = temp.Next() {
// 		result = append(result, *temp.Value.(*ActionShootResponse))
// 	}
// 	return result
// }

func (w *World) SpawnBullet(_player *Player) {
	_w := FindWeaponType(_player.WeaponOwned)
	b := NewBullet(w.list_bullet.Len()+1, _player.Id, _w.Bullet_id, _player.Pos_x, _player.Pos_y, _player.Rotation)
	w.list_bullet.PushBack(b)
}

func (w *World) DestroyBullet(_bullet *Bullet) {
	for temp := w.list_bullet.Front(); temp != nil; temp = temp.Next() {
		if temp.Value.(*Bullet) == _bullet {
			w.list_bullet.Remove(temp)
		}
	}
}

func (w *World) StartWorld() {
	var wg sync.WaitGroup
	w.initTime = MakeTimestamp()
	w.lastTime100Hz = w.initTime
	loop, _ := gloop.NewLoop(nil, nil, Hz200Delay, Hz30Delay)
	w.game_loop = loop
	render := func(step time.Duration) error {
		w.currTime100Hz = MakeTimestamp()
		w.deltaTime100Hz = float64(w.currTime100Hz-w.lastTime100Hz) / 1000
		w.lastTime100Hz = w.currTime100Hz

		return nil
	}
	simulate := func(step time.Duration) error {
		seq_counter = seq_counter + 1
		w.currTime30Hz = MakeTimestamp()
		w.Timestamp = float32(w.currTime30Hz-w.initTime) / 1000
		if Mult.GetStream().GetLen() != 0 {
			_ar := Mult.ReadAll()
			for _, r := range *_ar {
				w.RequestHandler(r)
			}
		}
		if w.start_game {
			if w.isNetworkBindCulling {
				for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
					p := temp.Value.(*Player)
					wg.Add(1)
					// fmt.Println("Before Filtered Snapshot: ", p)
					go w.sendFilteredSnapshot(seq_counter, p, &wg)
				}
				wg.Wait()
			} else {
				for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
					p := temp.Value.(*Player)
					wg.Add(1)
					// fmt.Println("Before Filtered Snapshot: ", p)
					go w.sendSnapshot(seq_counter, p, &wg)
				}
			}
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
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		fmt.Printf("%v\n", temp.Value.(*Player))
	}
	fmt.Print("--------------------------------------------\n")
	fmt.Println("Game World Stop")
	w.game_loop.Done()
}

func (w *World) sendFilteredSnapshot(seq int32, p *Player, wg *sync.WaitGroup) {
	s := w.GenerateSnapshot(seq, p)
	p.Conn.SendUnreliableOrdered(s)
	wg.Done()
}

func (w *World) sendSnapshot(seq int32, p *Player, wg *sync.WaitGroup) {
	s := w.GenerateSnapshotReliable(seq)
	p.Conn.SendUnreliableOrdered(s)
	wg.Done()
}

func (w *World) AddConn(conn *udpnetwork.Connection) {
	w.List_conn.PushBack(conn)
}

func (w *World) GenerateSnapshot(seq int32, p *Player) []byte {
	arH, arP := w.generateFilteredPlayerArray(p)
	n := NewSnapshot(seq, w.Timestamp, arP, arH, w.ListActionShootToArray())
	// n := NewSnapshot(seq, w.Timestamp, []Player{Player{1, "TEST12345", 123, 41, 234, 231, 23123, 123, 123, Idling}}, w.ListHitboxToArray(), w.list_action_shoot)
	b := n.MarshalSnapshot()
	// if p.Id == 1 {
	// 	fmt.Println("Snapshot: ", n)
	// }
	return b
}

func (w *World) GenerateSnapshotReliable(seq int32) []byte {
	n := NewSnapshot(seq, w.Timestamp, w.ListPlayerToArray(), w.ListHitboxToArray(), w.ListActionShootToArray())
	// fmt.Println("Snap: ", n)
	b := n.MarshalSnapshot()
	return b
}

// func (w *World) generateFilteredPlayerArray(player *Player) []Player {
// 	var result = make([]Player, 0, w.List_player.Len())
// 	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
// 		p := temp.Value.(*Player)
// 		if p.Pos_x < player.FOV && p.Pos_y < player.FOV {
// 			result = append(result, *p)
// 		}
// 	}
// 	return result
// }

func (w *World) generateFilteredPlayerArray(player *Player) ([]PlayerHitBox, []Player) {
	var arH = make([]PlayerHitBox, 0, w.List_player_hitbox.Len())
	var arP = make([]Player, 0, w.List_player.Len())
	for temp := w.List_player.Front(); temp != nil; temp = temp.Next() {
		p := temp.Value.(*Player)

		if p.CheckCulled(player.Pos_x, player.Pos_y, player.FOV) {
			arH = append(arH, *w.FindPlayerHitBoxInList(p))
			arP = append(arP, *p)
		} else {
			arH = append(arH, PlayerHitBox{})
			arP = append(arP, Player{})
		}
		// } else {
		// 	arH = append(arH, *w.FindPlayerHitBoxInList(p))
		// 	arP = append(arP, *player)
		// }
	}
	// fmt.Println("Player: ", player, "| arP: ", arP)
	return arH, arP
}

func (w *World) generateFilteredActionShootArray(player *Player) []ActionShootResponse {
	var result = make([]ActionShootResponse, 0, w.list_action_shoot.Len())
	for temp := w.list_action_shoot.Front(); temp != nil; temp = temp.Next() {
		p := temp.Value.(*ActionShootResponse)
		if p.CheckCulled(player.Pos_x, player.Pos_y, player.FOV) {
			result = append(result, *p)
		} else {
			result = append(result, ActionShootResponse{})
		}
	}
	return result
}

// if w.list_bullet.Len() != 0 {
// 	go func() {
// 		for tempHitBox := w.List_player_hitbox.Front(); tempHitBox != nil; tempHitBox = tempHitBox.Next() {
// 			// wg.Add(1)
// 			hitbox := tempHitBox.Value.(*PlayerHitBox)
// 			go func() {
// 				player := w.FindPlayerInListById(hitbox.Id)
// 				hit, dmg, bul := hitbox.CheckCollision(w.list_bullet)
// 				if hit {
// 					player.HitPlayer(dmg)
// 					fmt.Println("Hp: ", player.Hp)
// 					w.DestroyBullet(bul)
// 				}
// 				// wg.Done()
// 			}()
// 		}
// 		// wg.Wait()
// 		for tempBullet := w.list_bullet.Front(); tempBullet != nil; tempBullet = tempBullet.Next() {
// 			bul := tempBullet.Value.(*Bullet)
// 			if bul.Distance > FindBulletType(bul.Bullet_type).Range {
// 				w.DestroyBullet(bul)
// 			} else {
// 				go bul.MoveBullet(w.deltaTime100Hz, mutex)
// 			}
// 		}
// 	}()
// }
