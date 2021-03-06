package game

import (
	"math/rand"
	"time"

	"github.com/yohanesgre/new-server-shooter/pkg/udpnetwork"
)

//Random Direction
func RandDirection(min, max int) int {
	return min + rand.Int()*(max-min)
}

//Random float64
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

//Queue Enqueue
func Enqueue(queue []Request, element Request) []Request {
	queue = append(queue, element) // Simply append to enqueue.
	// fmt.Println("Enqueued:", element)
	return queue
}

//Queue Dequeue
func Dequeue(queue []Request) []Request {
	// element := queue[0] // The first element is the one to be dequeued.
	// fmt.Println("Dequeued:", element)
	return queue[1:] // Slice off the element once it is dequeued.
}

//Queue Peek
func Peek(queue []Request) Request {
	element := queue[0]
	return element
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Lerp(s, e, t float64) float64 { return s + (e-s)*t }

func RemoveListActionResponseElementAt(s []ActionShootResponse, i int) []ActionShootResponse {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func RemoveConnElementAt(s []*udpnetwork.Connection, i int) []*udpnetwork.Connection {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func CastInterfaceToInt32(val interface{}) int32 {
	switch i := val.(type) {
	case int8:
		return int32(i)
	case int16:
		return int32(i)
	case int32:
		return int32(i)
	default:
		return 0
	}
}
