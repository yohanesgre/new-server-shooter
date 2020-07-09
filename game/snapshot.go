package game

import "github.com/vmihailenco/msgpack/v5"

type Snapshot struct {
	Sequence  int32    `msgpack:"Sequence"`
	Timestamp float32  `msgpack:"Timestamp"`
	Players   []Player `msgpack:"Players,asArray"`
	Bullets   []Bullet `msgpack:"Bullets"`
}

func NewSnapshot(seq int32, timestamp float32, arp []Player, arb []Bullet) *Snapshot {
	return &Snapshot{seq, timestamp, arp, arb}
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
