package game

type ActionShootResponse struct {
	Id       int
	PlayerId int
}

func NewActionShootResponse(id, playerId int) *ActionShootResponse {
	return &ActionShootResponse{
		id, playerId,
	}
}
