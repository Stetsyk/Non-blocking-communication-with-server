package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"time"
)

func send(f func()) { // Non blocking message to server
	go func() {
		f()
	}()
}

func userCallback(n int) {
	fmt.Println("Server accepted function", n)
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	fmt.Println("Server released function", n)
}

type Sender struct {
	queue *list.List
}

func (s *Sender) tmpSendOnce() {
	if s.queue.Len() > 1 {
		s.queue.Remove(s.queue.Front())
		e := s.queue.Front()
		e.Value.(func())()
	}
}

func (s *Sender) sendOnce(f func()) {
	ff := func() {
		send(func() {
			f()
			s.tmpSendOnce()
		})
	}
	s.queue.PushBack(ff)
	if s.queue.Len() == 1 {
		ff()
	}
}

func main() {
	sender := &Sender{queue: list.New()}
	fmt.Println("Hello")
	for i := 0; i < 10; i++ {
		cpI := i
		sender.sendOnce(func() {
			userCallback(cpI)
		})
	}
	fmt.Println("World")
	time.Sleep(120 * time.Millisecond)
}
