package game

import (
	"container/list"
	"math"
)

type AgentHitBox struct {
	Id       int
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
