package game

import (
	"fmt"
	"math"
	"sync"
)

type BulletType struct {
	Id     int
	Damage float64
	Speed  float64
	Radius float64
}

var ListOfBulletType = []*BulletType{
	&BulletType{
		1,
		10,
		5,
		1,
	},
	&BulletType{
		2,
		30,
		3,
		1,
	},
	&BulletType{
		3,
		50,
		1,
		1,
	},
}

type Bullet struct {
	Id          int
	Owner_id    int
	Bullet_type int
	Pos_x       float64
	Pos_y       float64
	Rotation    float64
	Distance    float64
}

func NewBullet(_id int, _ownerId int, _bullet_type int, _pos_x float64, _pos_y float64, _rotation float64) *Bullet {
	b := new(Bullet)
	b.Id = _id
	b.Owner_id = _ownerId
	b.Bullet_type = _bullet_type
	b.Pos_x = _pos_x
	b.Pos_y = _pos_y
	b.Rotation = _rotation
	b.Distance = 0
	return b
}

func FindBulletType(_id int) *BulletType {
	var b *BulletType
	for i := 0; i < len(ListOfBulletType); i++ {
		if ListOfBulletType[i].Id == _id {
			b = ListOfBulletType[i]
		}
	}
	return b
}

func (b *Bullet) Destroy() {
	b.Destroy()
}

func (b *Bullet) MoveBullet(wg *sync.WaitGroup) {
	lastX := b.Pos_x
	lastY := b.Pos_y
	b.Pos_x = b.Pos_x + FindBulletType(b.Bullet_type).Speed*math.Cos(b.Rotation)*1/30 //*float64(time.Duration(int64(time.Second)/30))
	b.Pos_y = b.Pos_y + FindBulletType(b.Bullet_type).Speed*math.Sin(b.Rotation)*1/30 //*float64(time.Duration(int64(time.Second)/30))
	difX := b.Pos_x - lastX
	difY := b.Pos_y - lastY
	length := math.Sqrt(math.Pow(difX, 2) + math.Pow(difY, 2))
	b.Distance = b.Distance + length
	if b.Id == 1 {
		fmt.Println("Bullet Move: ", b.Id, " | ", b.Pos_x, ", ", b.Pos_y)
	}
	wg.Done()
}
