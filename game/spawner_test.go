package game

import "testing"

func TestSeedingSpanwner(t *testing.T) {
	SeedSpawnerPlayer(5)
	for i := 0; i < len(list_spawner_player); i++ {
		t.Logf("%+v", list_spawner_player[i])
	}
	t.Log("test")
}
