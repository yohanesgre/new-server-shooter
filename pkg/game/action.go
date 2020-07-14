package game

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
