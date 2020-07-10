package game

import (
	"math/rand"
	"time"
)

var list_spawner_player []*SpawnerPlayer

//Class Spawner Player
type SpawnerPlayer struct {
	Id     int
	Pos_x  float64
	Pos_y  float64
	Filled bool
}

//Seeding Spawner Player in Server Game World
func SeedSpawnerPlayer(_ammount int) {
	counter := 0
	for i := 0; i < _ammount; i++ {
		counter++
		rand.Seed(time.Now().UnixNano())
		var temp *SpawnerPlayer
		if counter == 1 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(1.0, 0.5),
				Pos_y:  RandFloat64(1.0, 0.5),
				Filled: false,
			}
		} else if counter == 2 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-2.0, -1.0),
				Pos_y:  RandFloat64(-2.0, -1.0),
				Filled: false,
			}
		} else if counter == 3 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-2.0, -1.0),
				Pos_y:  RandFloat64(-2.0, -1.0),
				Filled: false,
			}
		} else if counter == 4 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-2.0, -1.0),
				Pos_y:  RandFloat64(-2.0, -1.0),
				Filled: false,
			}
			counter = 0
		}
		list_spawner_player = append(list_spawner_player, temp)
	}
}
