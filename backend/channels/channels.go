package channels

import (
	"fmt"
	"sync"
)

type SocketData struct {
	inUse   bool
	channel chan []byte
	id      int
}

type EventSockets struct {
	sockets       []*SocketData
	globalLock    sync.Mutex
	channelsInUse int
	lastId        int
}

func New() *EventSockets {
	defaultSize := 10
	sockets := EventSockets{
		sockets:       make([]*SocketData, 0),
		channelsInUse: 0,
		lastId:        0,
	}

	sockets.increaseSizeUnsafe(0, defaultSize)

	return &sockets
}

func (sockets *EventSockets) FanInMessage(event []byte) {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := sockets.sockets[i]
		if c.inUse {
			fmt.Printf("talking to socket %d\n", i)
			c.channel <- event
		}
	}
}

func (sockets *EventSockets) AddChannel(newChannel chan []byte) int {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()

	sockets.lastId++
	newId := sockets.lastId

	sockets.addChannelUnsafe(newChannel, newId)
	fmt.Printf("created channel of id %d\n", newId)

	sockets.channelsInUse++
	return newId
}

func (sockets *EventSockets) RemoveChannel(id int) {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()

	for _, c := range sockets.sockets {
		if c.id == id {
			c.inUse = false
			close(c.channel)
			c.channel = nil
			c.id = -1
			sockets.channelsInUse--
		}
	}

	sockets.maybeDecreaseSizeUnsafe()
}

func (sockets *EventSockets) increaseSizeUnsafe(start, stop int) {
	for range stop - start {
		sockets.sockets = append(sockets.sockets, &SocketData{
			inUse:   false,
			channel: nil,
			id:      -1,
		})
	}
}

/**
* Only decreases the socket slice size iff we're less than 1/4 of the slots in the slice. We
* decrease by cutting the total number of slots in half and appending all in-use channels to
* the _end_ of the new slice
 */
func (sockets *EventSockets) maybeDecreaseSizeUnsafe() {
	if sockets.channelsInUse > len(sockets.sockets)/4 {
		return
	}

	newSlice := make([]*SocketData, len(sockets.sockets)/2)
	count := 0
	defer fmt.Printf("removed %d slots, %d slots remaining\n", count, len(sockets.sockets))
	for _, c := range sockets.sockets {
		if c.inUse {
			// append starting at the end to make inserts faster when we have new connections,
			// since we scan from the front
			newSlice[len(newSlice)-count-1] = c
			count++
		}
	}

	for i, c := range newSlice {
		if c == nil {
			newSlice[i] = &SocketData{
				inUse:   false,
				channel: nil,
				id:      -1,
			}
		}
	}

	sockets.sockets = newSlice
}

func (sockets *EventSockets) addChannelUnsafe(newChannel chan []byte, newId int) {
	for _, c := range sockets.sockets {
		if !c.inUse {
			c.channel = newChannel
			c.inUse = true
			c.id = newId
			return
		}
	}

	// We could not find an open slot. Increase the size of our slice
	currentSize := len(sockets.sockets)
	newSize := currentSize * 2
	sockets.increaseSizeUnsafe(currentSize, newSize)

	// This channel is empty, we can use it
	sockets.sockets[currentSize].channel = newChannel
	sockets.sockets[currentSize].inUse = true
	sockets.sockets[currentSize].id = newId
}
