package mux

import "sync"

type IDGenerator struct {
	sync.Mutex
	nextID int64
}

func (g *IDGenerator) NextID() int64 {
	g.Lock()
	defer g.Unlock()
	id := g.nextID
	g.nextID++
	return id
}
