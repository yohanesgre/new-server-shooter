package game

import (
	"math/rand"
	"time"
)

type Agent struct {
	Id       int
	Pos_x    float64
	Pos_y    float64
	Rotation float64
	Hp       float64
	State    PlayerState
	Hitbox   *AgentHitBox
	Ticker   *time.Ticker
	Done     chan bool
}

func NewAgent(id int, pos_x, pos_y float64) *Agent {
	ticker := time.NewTicker(15 * time.Millisecond)
	agent := new(Agent)
	agent.Id = id
	agent.Pos_x = pos_x
	agent.Pos_y = pos_y
	agent.Hp = 100
	agent.State = 1
	agent.Ticker = ticker
	agent.Done = make(chan bool)
	agent.Hitbox = NewAgentHitBox(agent)
	rand.Seed(time.Now().UnixNano())
	go func() {
		for {
			select {
			case <-agent.Done:
				agent.Ticker.Stop()
				break
			case <-ticker.C:
				agent.MoveAgent(Direction(RandDirection(1, 4)), RandFloat64(0, 360))
				agent.Hitbox.UpdateAgentHitBox(agent)
			}
		}
	}()
	return agent
}

func (a *Agent) MoveAgent(_d Direction, _a float64) {
	switch _d {
	case Top:
		a.Pos_y = Lerp(a.Pos_y, a.Pos_y+speed, 16.6666667)
		a.UpdateState(Walking)
	case Bottom:
		a.Pos_y = Lerp(a.Pos_y, a.Pos_y-speed, 16.6666667)
		a.UpdateState(Walking)
	case Right:
		a.Pos_x = Lerp(a.Pos_x, a.Pos_x+speed, 16.6666667)
		a.UpdateState(Walking)
	case Left:
		a.Pos_x = Lerp(a.Pos_x, a.Pos_x-speed, 16.6666667)
		a.UpdateState(Walking)
	case 0:
		a.UpdateState(Idling)
	}
	a.Rotation = _a
}

func (a *Agent) UpdateState(_state PlayerState) {
	a.State = _state
}

func (h *Agent) CheckCulled(pos_x, pos_y, fov float64) bool {
	dist := (h.Pos_x-pos_x)*(h.Pos_x-pos_x) + (h.Pos_y-pos_y)*(h.Pos_y-pos_y)
	if dist <= fov*fov {
		return true
	} else {
		return false
	}
}
