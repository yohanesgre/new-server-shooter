package game

import (
	"container/list"
	"math"
)

type AgentHitBox struct {
	Id       int32
	Pos_x    float64
	Pos_y    float64
	Height   float64
	Width    float64
	Rotation float64
}

func NewAgentHitBox(agent *Agent) *AgentHitBox {
	p := new(AgentHitBox)
	p.Id = agent.Id
	p.Pos_x = agent.Pos_x + OffsetX
	p.Pos_y = agent.Pos_y + OffsetY
	p.Height = HitboxHeight
	p.Width = HitboxWidth
	p.Rotation = agent.Rotation
	return p
}

func (h *AgentHitBox) UpdateAgentHitBox(agent *Agent) {
	h.Pos_x = float64(agent.Pos_x)
	h.Pos_y = float64(agent.Pos_y)
}

func (h *AgentHitBox) CheckCollision(list list.List) (bool, float64, *Bullet) {
	var bullet *Bullet
	hit, dmg := false, 0.0
	for temp := list.Front(); temp != nil; temp = temp.Next() {
		_b := temp.Value.(*Bullet)
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
	return hit, dmg, bullet
}

func (h *AgentHitBox) CheckHit(id int32, pos_x, pos_y float64) bool {

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
