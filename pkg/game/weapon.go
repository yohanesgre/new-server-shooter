package game

type Weapon struct {
	Id        int
	Name      string
	Bullet_id int
	Ammo      int
}

type WeaponDrop struct {
	Id       int
	Type_id  int
	Pos_x    float64
	Pos_y    float64
	Owner_id int
}

var weapon_db = [...]Weapon{
	Weapon{
		Id:        1,
		Name:      "Handgun",
		Bullet_id: 1,
		Ammo:      10,
	},
	Weapon{
		Id:        2,
		Name:      "Rifle",
		Bullet_id: 2,
		Ammo:      20,
	},
	Weapon{
		Id:        4,
		Name:      "Sniper",
		Bullet_id: 3,
		Ammo:      5,
	},
}

func NewWeaponDrop(_id int, _type int, _pos_x float64, _pos_y float64) *WeaponDrop {
	w := new(WeaponDrop)
	w.Id = _id
	w.Type_id = _type
	w.Pos_x = _pos_x
	w.Pos_y = _pos_y
	w.Owner_id = 0
	return w
}

func FindWeaponType(_id int) Weapon {
	var w Weapon
	for i := 0; i < len(weapon_db); i++ {
		if weapon_db[i].Id == _id {
			w = weapon_db[i]
		}
	}
	return w
}
