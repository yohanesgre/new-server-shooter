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
				Pos_x:  0.0,
				Pos_y:  0.0,
				Filled: false,
			}
		} else if counter == 2 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(20, 30),
				Pos_y:  RandFloat64(-7, 4),
				Filled: false,
			}
		} else if counter == 3 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-2, 10),
				Pos_y:  RandFloat64(10, 17),
				Filled: false,
			}
		} else if counter == 4 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(14.25, 31),
				Pos_y:  RandFloat64(-19.51, -13.84),
				Filled: false,
			}
		} else if counter == 5 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-30, -17),
				Pos_y:  RandFloat64(-21, -15),
				Filled: false,
			}
		} else if counter == 6 {
			temp = &SpawnerPlayer{
				Id:     i + 1,
				Pos_x:  RandFloat64(-32, -25),
				Pos_y:  RandFloat64(-2, 5),
				Filled: false,
			}
			counter = 0
		}
		list_spawner_player = append(list_spawner_player, temp)
	}
}
