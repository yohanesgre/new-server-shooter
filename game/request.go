package game

import (
	"github.com/vmihailenco/msgpack/v5"
)

type EndpointType int

const (
	JOIN       EndpointType = 1
	MOVE       EndpointType = 2
	SHOOT      EndpointType = 3
	SHOOT_DONE EndpointType = 4
	COLIDED    EndpointType = 5
)

type Request struct {
	Endpoint EndpointType
	Payload  interface{}
}

type RequestMove struct {
	Id        int
	Direction Direction
	Rotation  float64
	State     PlayerState
}

type RequestJoin struct {
	Name string
	FOV  float64
}

type RequestShoot struct {
	Id int
}

type RequestShootDone struct {
	Id int
}

type RequestBulletColided struct {
	HittedId int
	Pos_x    float64
	Pos_y    float64
	WeaponId int
}

type RequestReload struct {
	Id int
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

func (r *Request) PayloadToRequestMove() *RequestMove {
	p := r.Payload.(map[string]interface{})
	return &RequestMove{
		int(p["Id"].(int8)),
		Direction(p["Direction"].(int8)),
		p["Rotation"].(float64),
		PlayerState(p["State"].(int8)),
	}
}

func (r *Request) PayloadToRequestJoin() *RequestJoin {
	p := r.Payload.(map[string]interface{})
	return &RequestJoin{
		p["Name"].(string),
		p["FOV"].(float64),
	}
}

func (r *Request) PayloadToRequestBulletColided() *RequestBulletColided {
	p := r.Payload.(map[string]interface{})
	return &RequestBulletColided{
		int(p["HittedId"].(int8)),
		p["Pos_x"].(float64),
		p["Pos_y"].(float64),
		int(p["WeaponId"].(int8)),
	}
}

func (r *Request) PayloadToRequestShoot() *RequestShoot {
	p := r.Payload.(map[string]interface{})
	return &RequestShoot{
		int(p["Id"].(int8)),
	}
}

func (r *Request) PayloadToRequestShootDone() *RequestShootDone {
	p := r.Payload.(map[string]interface{})
	return &RequestShootDone{
		int(p["Id"].(int8)),
	}
}

func (r *Request) PayloadToRequestReload() *RequestReload {
	p := r.Payload.(map[string]interface{})
	return &RequestReload{
		int(p["Id"].(int8)),
	}
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
