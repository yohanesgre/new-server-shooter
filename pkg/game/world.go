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
	max_player           int
	List_conn            list.List
	list_player_hitbox   list.List
	list_player          list.List
	list_bullet          list.List
	list_weapon          list.List
	list_action_shoot    []ActionShootResponse
	game_loop            *gloop.Loop
	Timestamp            float32
	currTime30Hz         int64
	initTime             int64
	deltaTime100Hz       float64
	currTime100Hz        int64
	lastTime100Hz        int64
	action_shoot_counter int
}

func NewWorld(max_player int) *World {
	w := new(World)
	w.max_player = max_player
	w.list_player.Init()
	w.list_bullet.Init()
	w.list_weapon.Init()
	w.list_action_shoot = make([]ActionShootResponse, 0, 20)
	SeedSpawnerPlayer(max_player)
	Mult = NewMultiplexer(120)
	w.action_shoot_counter = 0
	return w
}

func (w *World) Destroy() {
	w = nil
}

//Request Handler
//Handle all request sent by Players
func (w *World) RequestHandler(_r Request) {
	// if w.list_player.Len() != 0 {
	// 	w.ForceAllPlayersIdle()
	// }
	// if w.list_action_shoot.Len() != 0 {
	// 	fmt.Println("List Actions: ", w.ListActionShootToArray())
	// }
	switch _r.Endpoint {
	case JOIN:
		p := _r.PayloadToRequestJoin()
		for i := 0; i < len(list_spawner_player); i++ {
			spawner := list_spawner_player[i]
			if !spawner.Filled {
				spawner.Filled = true
				p_ := NewPlayer(w.list_player.Len()+1, p.Name, spawner.Pos_x, spawner.Pos_y, 0.0, p.FOV)
				h_ := NewPlayerHitBox(p_)
				w.list_player.PushBack(p_)
				w.list_player_hitbox.PushBack(h_)
				snap := w.GenerateSnapshot(seq_counter)
				for temp := w.List_conn.Front(); temp != nil; temp = temp.Next() {
					temp.Value.(*udpnetwork.Connection).SendReliableOrdered(snap)
				}
				break
			}
		}
	case MOVE:
		r := _r.PayloadToRequestMove()
		p := w.FindPlayerInListById(r.Id)
		h := w.FindPlayerHitBoxInList(p)
		mutex.Lock()
		p.Move(r.Direction, r.Rotation, w.deltaTime100Hz)
		h.UpdatePlayerHitBox(p)
		mutex.Unlock()
	case SHOOT:
		w.action_shoot_counter = w.action_shoot_counter + 1
		r := _r.PayloadToRequestShoot()
		p := w.FindPlayerInListById(r.Id)
		p.UpdateState(Shooting)
		w.list_action_shoot = append(w.list_action_shoot, *NewActionShootResponse(w.action_shoot_counter, r.Id))
	case SHOOT_DONE:
		r := _r.PayloadToRequestShootDone()
		mutex.Lock()
		for i := 0; i < len(w.list_action_shoot); i++ {
			if w.list_action_shoot[i].Id == r.Id {
				w.list_action_shoot = RemoveListActionResponseElementAt(w.list_action_shoot, i)
			}
		}
		mutex.Unlock()
	case COLIDED:
		r := _r.PayloadToRequestBulletColided()
		hitbox := w.FindPlayerHitBoxInListByHittedId(r.HittedId)
		if hitbox.CheckHit(r.HittedId, r.Pos_x, r.Pos_y) {
			dmg := FindBulletType(FindWeaponType(r.WeaponId).Bullet_id).Damage
			// fmt.Println("Dmg: ", dmg)
			w.FindPlayerInListById(r.HittedId).HitPlayer(dmg)
			// fmt.Println("Hp: ", w.FindPlayerInListById(r.HittedId).Hp)
		}
	}
}

//Find Player Pointer in List Of Players' Pointer
func (w *World) ForceAllPlayersIdle() {
	for temp := w.list_player.Front(); temp != nil; temp = temp.Next() {
		_p := temp.Value.(*Player)
		_p.UpdateState(Idling)
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

//Find PlayerHitBox Pointer in List Of PlayerHitBox' Pointer with HittedId
func (w *World) FindPlayerHitBoxInListByHittedId(_hittedId int) *PlayerHitBox {
	var result *PlayerHitBox
	for temp := w.list_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		_h := temp.Value.(*PlayerHitBox)
		if _h.Id == _hittedId {
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

func (w *World) ListHitboxToArray() []PlayerHitBox {
	result := make([]PlayerHitBox, 0, w.list_player_hitbox.Len())
	for temp := w.list_player_hitbox.Front(); temp != nil; temp = temp.Next() {
		result = append(result, *temp.Value.(*PlayerHitBox))
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
	// var wg sync.WaitGroup
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
		for temp := w.List_conn.Front(); temp != nil; temp = temp.Next() {
			c := temp.Value.(*udpnetwork.Connection)
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

func (w *World) AddConn(conn *udpnetwork.Connection) {
	w.List_conn.PushBack(conn)
}

func (w *World) GenerateSnapshot(seq int32) []byte {
	n := NewSnapshot(seq, w.Timestamp, w.ListPlayerToArray(), w.ListHitboxToArray(), w.list_action_shoot)
	b := n.MarshalSnapshot()
	// fmt.Println("Snapshot: ", n)
	return b
}

// if w.list_bullet.Len() != 0 {
// 	go func() {
// 		for tempHitBox := w.list_player_hitbox.Front(); tempHitBox != nil; tempHitBox = tempHitBox.Next() {
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