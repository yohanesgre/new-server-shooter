package game

import (
	"github.com/yohanesgre/new-server-shooter/pkg/udpnetwork"
)

type PlayerState int
type Direction int

const (
	Idling   PlayerState = 1
	Walking  PlayerState = 2
	Shooting PlayerState = 3
	Top      Direction   = 1
	Bottom   Direction   = 2
	Right    Direction   = 3
	Left     Direction   = 4
	speed                = 2
)

type Player struct {
	Id          int
	Name        string
	Pos_x       float64
	Pos_y       float64
	Rotation    float64
	FOV         float64
	Hp          float64
	Ammo        int
	WeaponOwned int
	State       PlayerState
	Conn        *udpnetwork.Connection
}

func NewPlayer(_id int, _name string, _pos_x, _pos_y, _rotation, _fov float64, conn *udpnetwork.Connection) *Player {
	p := &Player{
		_id, _name, _pos_x, _pos_y, _rotation, _fov, 100, 0, 1, 1, conn,
	}
	return p
}

func (p *Player) UpdatePlayer(_p *Player) {
	p.Pos_x = _p.Pos_x
	p.Pos_y = _p.Pos_y
	p.Rotation = _p.Rotation
	p.State = 1
}

func (p *Player) Move(_d Direction, _a float64, dtime float64) {
	switch _d {
	case Top:
		p.Pos_y = Lerp(p.Pos_y, p.Pos_y+speed, dtime)
		p.UpdateState(Walking)
	case Bottom:
		p.Pos_y = Lerp(p.Pos_y, p.Pos_y-speed, dtime)
		p.UpdateState(Walking)
	case Right:
		p.Pos_x = Lerp(p.Pos_x, p.Pos_x+speed, dtime)
		p.UpdateState(Walking)
	case Left:
		p.Pos_x = Lerp(p.Pos_x, p.Pos_x-speed, dtime)
		p.UpdateState(Walking)
	case 0:
		p.UpdateState(Idling)
	}
	p.Rotation = _a
}

func (p *Player) PickWeapon(_w WeaponDrop) {
	p.WeaponOwned = _w.Id
	p.Ammo = FindWeaponType(p.WeaponOwned).Ammo
}

func (p *Player) Reload() {
	p.Ammo = FindWeaponType(p.WeaponOwned).Ammo
}

func (p *Player) Shoot() {
	p.Ammo = p.Ammo - 1
}

func (p *Player) HitPlayer(_dmg float64) {
	p.Hp = p.Hp - _dmg
}

func (p *Player) Destroy() {
	p.Destroy()
}

func (p *Player) UpdateState(_state PlayerState) {
	p.State = _state
}

func (h *Player) CheckCulled(pos_x, pos_y, fov float64) bool {
	dist := (h.Pos_x-pos_x)*(h.Pos_x-pos_x) + (h.Pos_y-pos_y)*(h.Pos_y-pos_y)
	if dist <= fov*fov {
		return true
	} else {
		return false
	}
}
