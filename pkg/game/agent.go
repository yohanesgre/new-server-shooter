package game

import (
	"math"
	"math/rand"
	"time"
)

type AgentSnapshot struct {
	Id       int32
	Pos_x    float64
	Pos_y    float64
	Rotation float64
	Hp       float64
	State    PlayerState
}

func NewAgentSnapshot(agent *Agent) AgentSnapshot {
	return AgentSnapshot{
		agent.Id,
		agent.Pos_x,
		agent.Pos_y,
		agent.Rotation,
		agent.Hp,
		agent.State,
	}
}

type Agent struct {
	Id        int32
	Pos_x     float64
	Pos_y     float64
	Rotation  float64
	Hp        float64
	State     PlayerState
	Hitbox    *AgentHitBox
	Ticker    *time.Ticker
	Counter   int
	LastTime  int64
	DeltaTime float64
	Done      chan bool
}

func NewAgent(id int32, pos_x, pos_y float64) *Agent {
	ticker := time.NewTicker(200 * time.Millisecond)
	agent := new(Agent)
	agent.Id = id
	agent.Pos_x = pos_x
	agent.Pos_y = pos_y
	agent.Hp = 100
	agent.State = 1
	agent.Ticker = ticker
	agent.Done = make(chan bool)
	agent.Hitbox = NewAgentHitBox(agent)
	agent.LastTime = MakeTimestamp()
	agent.Counter = 0
	if agent.Id%2 == 0 {
		go func() {
			for {
				select {
				case <-agent.Done:
					agent.Ticker.Stop()
					break
				case t := <-ticker.C:
					rand.Seed(time.Now().UnixNano())
					currTime := t.UnixNano() / int64(time.Millisecond)
					agent.DeltaTime = float64(currTime-agent.LastTime) / 1000
					agent.LastTime = currTime
					agent.MoveAgent(Direction(rand.Intn(4-1+1)+1), agent.DeltaTime)
					agent.Hitbox.UpdateAgentHitBox(agent)
				}
			}
		}()
	} else {
		go func() {
			for {
				select {
				case <-agent.Done:
					agent.Ticker.Stop()
					break
				case t := <-ticker.C:
					rand.Seed(time.Now().UnixNano())
					currTime := t.UnixNano() / int64(time.Millisecond)
					agent.DeltaTime = float64(currTime-agent.LastTime) / 1000
					agent.LastTime = currTime
					if agent.Counter == 0 {
						agent.MoveAgent(Right, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
					} else if agent.Counter == 1 {
						agent.MoveAgent(Bottom, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
					} else if agent.Counter == 2 {
						agent.MoveAgent(Bottom, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
					} else if agent.Counter == 3 {
						agent.MoveAgent(Left, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
					} else if agent.Counter == 4 {
						agent.MoveAgent(Top, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
					} else if agent.Counter == 5 {
						agent.MoveAgent(Top, agent.DeltaTime)
						agent.Hitbox.UpdateAgentHitBox(agent)
						agent.Counter = 0
					}
					counter++
				}
			}
		}()
	}
	return agent
}

func (a *Agent) HitAgent(_dmg float64) {
	a.Hp = a.Hp - _dmg
}

func (a *Agent) MoveAgent(_d Direction, dt float64) {
	switch _d {
	case Top:
		a.Pos_y = Lerp(a.Pos_y, a.Pos_y+1, dt)
		a.UpdateState(Walking)
	case Bottom:
		a.Pos_y = Lerp(a.Pos_y, a.Pos_y-1, dt)
		a.UpdateState(Walking)
	case Right:
		a.Pos_x = Lerp(a.Pos_x, a.Pos_x+1, dt)
		a.UpdateState(Walking)
	case Left:
		a.Pos_x = Lerp(a.Pos_x, a.Pos_x-1, dt)
		a.UpdateState(Walking)
	case 0:
		a.UpdateState(Idling)
	}
}

func (a *Agent) UpdateState(_state PlayerState) {
	a.State = _state
}

func (h *Agent) CheckCulled(pos_x, pos_y, fov float64) bool {
	return (math.Pow(h.Pos_x-pos_x, 2) + math.Pow(h.Pos_y-pos_y, 2)) <= math.Pow(fov, 2)
}
