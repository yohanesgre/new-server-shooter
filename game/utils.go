package game

import (
	"fmt"
	"math/rand"
)

//Random Float64
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

//Queue Enqueue
func Enqueue(queue []Request, element Request) []Request {
	queue = append(queue, element) // Simply append to enqueue.
	fmt.Println("Enqueued:", element)
	return queue
}

//Queue Dequeue
func Dequeue(queue []Request) []Request {
	element := queue[0] // The first element is the one to be dequeued.
	fmt.Println("Dequeued:", element)
	return queue[1:] // Slice off the element once it is dequeued.
}

//Queue Peek
func Peek(queue []Request) Request {
	element := queue[0]
	return element
}
