package game

import "github.com/vmihailenco/msgpack/v5"

type Snapshot struct {
	Sequence  int
	Timestamp int64
	Players   []Player
	Bullets   []Bullet
}

func NewSnapshot(seq int, timestamp int64, arp []Player, arb []Bullet) *Snapshot {
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
