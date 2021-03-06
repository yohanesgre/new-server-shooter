package game

import (
	"testing"
)

func TestWorld(t *testing.T) {
	// var wg sync.WaitGroup
	// w := NewWorld(5)
	// t.Logf("World: %#v", w)
	// w.StartWorld()
	// Demul.Write(*makeRequestJoin(1))
	// time.Sleep(1000 * time.Millisecond)
	// Demul.Write(*makeRequestJoin(2))
	// for i := 0; i < 2; i++ {
	// 	wg.Add(1)
	// 	go worker(Demul, i+1, &wg)
	// }
	// wg.Wait()
	// t.Log("End")
	// w.StopWorld()
	w := NewWorld(1, true)
	w.StartWorld()
	w.spawnAgents(100)
	for temp := w.List_agent.Front(); temp != nil; temp = temp.Next() {
		a := NewAgentSnapshot(temp.Value.(*Agent))
		t.Log("Agent", a.Id, " : ", a.Pos_x, " | ", a.Pos_y)
	}
	w.StopWorld()
}

// func makeRequestJoin(i int) *Request {
// 	p := NewPlayer(i, "Test"+strconv.Itoa(i), 0, 0, 0, 10.0)
// 	r := NewRequest(JOIN, *p)
// 	_m := r.MarshalRequest()
// 	_u := UnmarshalRequest(_m)
// 	return &_u
// }

// func worker(demul *Multiplexer, id int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for i := 0; i < 11; i++ {
// 		p := NewPlayer(id, "Test"+strconv.Itoa(id), 0.0+float64(i), 0.0+float64(i), 0.0+float64(i), 10.0)
// 		r := NewRequest(MOVE, *p)
// 		_m := r.MarshalRequest()
// 		_u := UnmarshalRequest(_m)
// 		demul.Write(_u)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
