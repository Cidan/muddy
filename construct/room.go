package construct

import (
	"sync"

	"github.com/Cidan/muddy/atlas"
	roomv1 "github.com/Cidan/muddy/gen/proto/go/room/v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/proto"
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

func (r *Room) AddPlayer(uuid string, p atlas.Player) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.players[uuid] = p
}

func (r *Room) RemovePlayer(uuid string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.players, uuid)
}

func (r *Room) Save() error {
	r.lock.RLock()
	defer r.lock.RUnlock()
	data, err := proto.Marshal(r.data)
	if err != nil {
		return err
	}
	// TODO(lobato): save data
	_ = data
	return nil
}

func (r *Room) Send(text string, except atlas.Player, args ...interface{}) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for id, p := range r.players {
		if id == except.Uuid() {
			continue
		}
		p.Send(text, args...)
	}
}

func (r *Room) SetName(text string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.data.Name = text
}

func (r *Room) Name() string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.data.Name
}
