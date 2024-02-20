package main

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(weight int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, weight),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.ch
}
