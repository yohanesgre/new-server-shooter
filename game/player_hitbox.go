package game

import (
	"container/list"
	"math"
)

type PlayerHitBox struct {
	Id     int
	Pos_x  float64
	Pos_y  float64
	Height float64
	Width  float64
}

func NewPlayerHitBox(player *Player, height, width float64) *PlayerHitBox {
	p := new(PlayerHitBox)
	p.Id = player.Id
	p.Pos_x = float64(player.Pos_x)
	p.Pos_y = float64(player.Pos_y)
	p.Height = height
	p.Width = width
	return p
}

func (h *PlayerHitBox) UpdatePlayerHitBox(player *Player) {
	h.Pos_x = float64(player.Pos_x)
	h.Pos_y = float64(player.Pos_y)
}

func (h *PlayerHitBox) CheckCollision(list list.List) (bool, float64, *Bullet) {
	var bullet *Bullet
	hit, dmg := false, 0.0
	for temp := list.Front(); temp != nil; temp = temp.Next() {
		_b := temp.Value.(*Bullet)
		if _b.Owner_id != h.Id {
			distX := math.Abs(_b.Pos_x - h.Pos_x - h.Height/2)
			distY := math.Abs(_b.Pos_y - h.Pos_y - h.Width/2)
			if distX > (h.Width/2 + FindBulletType(_b.Bullet_type).Radius) {
				bullet = nil
				break
			}
			if distY > (h.Height/2 + FindBulletType(_b.Bullet_type).Radius) {
				bullet = nil
				break
			}
			if distX <= (h.Width / 2) {
				bullet = nil
				break
			}
			if distY <= (h.Height / 2) {
				bullet = nil
				break
			}
			dx := distX - h.Width/2
			dy := distY - h.Height/2
			hit = dx*dx+dy*dy <= (FindBulletType(_b.Bullet_type).Radius * FindBulletType(_b.Bullet_type).Radius)
			dmg = FindBulletType(_b.Id).Damage
			bullet = bullet
			break
		}
	}
	return hit, dmg, bullet
}
