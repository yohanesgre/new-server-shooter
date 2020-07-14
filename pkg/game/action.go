package game

import (
	"math"
)

type ActionShootResponse struct {
	Id       int
	PlayerId int
	Pos_x    float64
	Pos_y    float64
}

func NewActionShootResponse(id, playerId int, pos_x, pos_y float64) *ActionShootResponse {
	return &ActionShootResponse{
		id, playerId, pos_x, pos_y,
	}
}

func (a *ActionShootResponse) CheckCulled(pos_x, pos_y, fov float64) bool {
	if math.Pow((pos_x-a.Pos_x), 2)+math.Pow((pos_y-a.Pos_y), 2) < math.Pow(fov, 2) {
		return true
	} else {
		return false
	}
}
