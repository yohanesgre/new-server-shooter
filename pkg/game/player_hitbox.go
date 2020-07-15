package game

import (
	"container/list"
	"math"
)

const (
	HitboxWidth  float64 = 0.5803897
	HitboxHeight float64 = 0.6792674
	OffsetX      float64 = 0.01059404
	OffsetY      float64 = 0.1583269
)

type PlayerHitBox struct {
	Id       int
	Pos_x    float64
	Pos_y    float64
	Height   float64
	Width    float64
	Rotation float64
}

func NewPlayerHitBox(player *Player) *PlayerHitBox {
	p := new(PlayerHitBox)
	p.Id = player.Id
	p.Pos_x = player.Pos_x + OffsetX
	p.Pos_y = player.Pos_y + OffsetY
	p.Height = HitboxHeight
	p.Width = HitboxWidth
	p.Rotation = player.Rotation
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
				bullet = _b
			}
			if distY <= (h.Height / 2) {
				bullet = _b
			}
			dx := distX - h.Width/2
			dy := distY - h.Height/2
			hit = dx*dx+dy*dy <= (FindBulletType(_b.Bullet_type).Radius * FindBulletType(_b.Bullet_type).Radius)
			dmg = FindBulletType(_b.Id).Damage
			bullet = _b
			break
		}
	}
	return hit, dmg, bullet
}

func (h *PlayerHitBox) CheckHit(id int, pos_x, pos_y float64) bool {

	// temporary variables to set edges for testing
	testX := pos_x
	testY := pos_y

	// which edge is closest?
	if pos_x < h.Pos_x { // test left edge
		testX = h.Pos_x
	} else if pos_x > (h.Pos_x + h.Width) { // right edge
		testX = h.Pos_x + h.Width
	}

	if pos_x < h.Pos_y { // top edge
		testY = h.Pos_y
	} else if pos_y > (h.Pos_y + h.Height) { // bottom edge
		testY = h.Pos_y + h.Height
	}

	// get distance from closest edges
	distX := pos_x - testX
	distY := pos_y - testY
	distance := math.Sqrt((distX * distX) + (distY * distY))

	// if the distance is less than the radius, collision!
	if distance <= 2 {
		return true
	}
	return false
}
