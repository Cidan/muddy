package construct

import (
	"sync"

	"github.com/Cidan/muddy/atlas"
	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
	uuid "github.com/satori/go.uuid"
)

type Room struct {
	data    *roomv1.Room
	players map[string]atlas.Player
	lock    sync.RWMutex
}

func NewRoom() *Room {
	r := &Room{
		data: &roomv1.Room{
			Uuid: uuid.NewV4().String(),
		},
		lock:    sync.RWMutex{},
		players: make(map[string]atlas.Player),
	}
	return r
}

func (r *Room) SetCoordinates(c *roomv1.Room_Coordinates) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data.Coordinates = c
}

func (r *Room) Coordinates() *roomv1.Room_Coordinates {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.data.Coordinates
}

func (r *Room) AddPlayer(p atlas.Player) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.players[p.Uuid()] = p
}

func (r *Room) RemovePlayer(p atlas.Player) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.players, p.Uuid())
}
