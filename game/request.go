package game

import (
	"github.com/vmihailenco/msgpack/v5"
)

type EndpointType int

const (
	JOIN    EndpointType = 1
	MOVE    EndpointType = 2
	SHOOT   EndpointType = 3
	RESPAWN EndpointType = 4
)

type Request struct {
	Endpoint EndpointType
	Payload  interface{}
}

func NewRequest(_endpoint EndpointType, _payload interface{}) *Request {
	r := new(Request)
	r.Endpoint = _endpoint
	r.Payload = _payload
	return r
}

func (r *Request) MarshalRequest() []byte {
	b, err := msgpack.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}

func UnmarshalRequest(_b []byte) Request {
	var result Request
	err := msgpack.Unmarshal(_b, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func (r *Request) PayloadToPlayer() *Player {
	p := r.Payload.(map[string]interface{})
	return &Player{
		int(p["Id"].(int8)),
		p["Name"].(string),
		p["Pos_x"].(float64),
		p["Pos_y"].(float64),
		p["Rotation"].(float64),
		p["FOV"].(float64),
		p["Hp"].(float64),
		int(p["Ammo"].(int8)),
		int(p["WeaponOwned"].(int8)),
		PlayerState(p["State"].(int8)),
	}
}

func (r *Request) PayloadToBullet() Bullet {
	p := r.Payload.(map[string]interface{})
	return Bullet{
		int(p["Id"].(int64)),
		int(p["Owner_id"].(int64)),
		int(p["Bullet_type"].(int64)),
		float64(p["Pos_x"].(float64)),
		float64(p["Pos_y"].(float64)),
		float64(p["Rotation"].(float64)),
		float64(p["Distance"].(float64)),
	}
}

func (r *Request) PayloadToWeaponDrop() WeaponDrop {
	p := r.Payload.(map[string]interface{})
	return WeaponDrop{
		int(p["Id"].(int64)),
		int(p["Type_id"].(int64)),
		float64(p["Pos_x"].(float64)),
		float64(p["Pos_y"].(float64)),
		int(p["Owner_id"].(float64)),
	}
}