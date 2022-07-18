package construct

import (
	"sync"

	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
	uuid "github.com/satori/go.uuid"
)

type Room struct {
	data    *roomv1.Room
	players map[string]*Player
	lock    sync.RWMutex
}

func NewRoom() *Room {
	r := &Room{
		data: &roomv1.Room{
			Uuid: uuid.NewV4().String(),
		},
		lock: sync.RWMutex{},
	}
	return r
}

func (r *Room) SetCoordinates(c *roomv1.Room_Coordinates) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data.Coordinates = c
}
