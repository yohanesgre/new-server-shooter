package loadbalancer

import "github.com/vmihailenco/msgpack/v5"

type EndpointType int

const (
	GET_LOBBIES      EndpointType = 1
	CREATE_LOBBY     EndpointType = 2
	CONNECT_LOBBY    EndpointType = 3
	DISCONNECT_LOBBY EndpointType = 4
	START_GAME       EndpointType = 5
)

type Request struct {
	Endpoint int
	Payload  []string
}

func NewRequest(_endpoint int, _payload []string) *Request {
	r := new(Request)
	r.Endpoint = _endpoint
	r.Payload = _payload
	return r
}

func UnmarshalRequest(_b []byte) Request {
	var result Request
	err := msgpack.Unmarshal(_b, &result)
	if err != nil {
		panic(err)
	}
	return result
}
