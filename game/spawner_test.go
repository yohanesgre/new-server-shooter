package game

import "testing"

func TestSeedingSpanwner(t *testing.T) {
	SeedSpawnerPlayer(5)
	for {
		for i := 0; i < len(list_spawner_player); i++ {
			t.Logf("%+v", list_spawner_player[i])
			spawner := list_spawner_player[i]
			if !spawner.Filled {
				spawner.Filled = true
				break
			}
		}
		if list_spawner_player[4].Filled == true {
			break
		}
	}
}
