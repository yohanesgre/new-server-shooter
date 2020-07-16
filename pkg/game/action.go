package game

type ActionShootResponse struct {
	Id       int32
	PlayerId int
	Pos_x    float64
	Pos_y    float64
}

func NewActionShootResponse(id int32, playerId int, pos_x, pos_y float64) *ActionShootResponse {
	return &ActionShootResponse{
		id, playerId, pos_x, pos_y,
	}
}

func (h *ActionShootResponse) CheckCulled(pos_x, pos_y, fov float64) bool {
	dist := (h.Pos_x-pos_x)*(h.Pos_x-pos_x) + (h.Pos_y-pos_y)*(h.Pos_y-pos_y)
	if dist <= fov*fov {
		return true
	} else {
		return false
	}
}
