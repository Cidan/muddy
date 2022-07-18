package atlas

import (
	"fmt"
	"sync"
)

var (
	atlas *Atlas
	once  sync.Once
)

type Atlas struct {
	worldMap map[string]Room
	lock     sync.RWMutex
}

func Get() *Atlas {
	once.Do(func() {
		atlas = &Atlas{
			worldMap: make(map[string]Room),
		}
	})
	return atlas
}

func (a *Atlas) AddRoom(r Room) {
	c := r.Coordinates()
	d := fmt.Sprintf("%d_%d_%d", c.X, c.Y, c.Z)
	a.lock.Lock()
	defer a.lock.Unlock()
	a.worldMap[d] = r
}

func (a *Atlas) Room(x, y, z int32) Room {
	d := fmt.Sprintf("%d_%d_%d", x, y, z)
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.worldMap[d]
}
