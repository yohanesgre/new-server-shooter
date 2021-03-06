package game

import "github.com/vmihailenco/msgpack/v5"

type Snapshot struct {
	Sequence  int32    `msgpack:"Sequence"`
	Timestamp float32  `msgpack:"Timestamp"`
	Players   []Player `msgpack:"Players,asArray"`
	// Hitboxs      []PlayerHitBox        `msgpack:"Hitboxs,asArray"`
	Agents []AgentSnapshot `msgpack:"Agents,asArray"`
	// AgentHitboxs []AgentHitBox         `msgpack:"AgentHitboxs,asArray"`
	ActionShoots []ActionShootResponse `msgpack:"ActionShoots,asArray"`
	// Bullets   []Bullet       `msgpack:"Bullets,asArray"`
}

func NewSnapshot(seq int32, timestamp float32, arp []Player, ara []AgentSnapshot, aras []ActionShootResponse) *Snapshot {
	return &Snapshot{seq, timestamp, arp, ara, aras}
}

func NewReliableSnapshot(seq int32, timestamp float32, arp []Player, aras []ActionShootResponse) *Snapshot {
	return &Snapshot{seq, timestamp, arp, nil, aras}
}

func (r *Snapshot) MarshalSnapshot() []byte {
	b, err := msgpack.Marshal(r)
	if err != nil {
		panic(err)
	}
	return b
}

func UnmarshalSnapshot(_b []byte) Snapshot {
	var result Snapshot
	err := msgpack.Unmarshal(_b, &result)
	if err != nil {
		panic(err)
	}
	return result
}
