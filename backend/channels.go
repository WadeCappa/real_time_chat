package main

import (
	"fmt"
	"sync"
)

type SocketData struct {
	inUse    bool
	uniqueId uint64
	channel  chan []byte
}
type EventSockets struct {
	sockets    []SocketData
	lastId     uint64
	globalLock sync.Mutex
}

func (sockets *EventSockets) FanInMessage(event []byte) {
	eventCopy := event
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if c.inUse {
			fmt.Printf("talking to socket %d\n", c.uniqueId)
			c.channel <- eventCopy
		}
	}
}

func (sockets *EventSockets) AddChannel(newChannel chan []byte) uint64 {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	sockets.lastId++
	defer fmt.Printf("created channel of id %d\n", sockets.lastId)
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if !c.inUse {
			c.channel = newChannel
			c.inUse = true
			c.uniqueId = sockets.lastId
			return sockets.lastId
		}
	}
	// Only if we can't re-use a spot in our slice, do we append a new channel
	sockets.sockets = append(sockets.sockets, SocketData{
		channel:  newChannel,
		uniqueId: sockets.lastId,
		inUse:    true,
	})
	return sockets.lastId
}

func (sockets *EventSockets) RemoveChannel(id uint64) {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := &sockets.sockets[i]
		if c.uniqueId == id {
			close(c.channel)
			c.channel = nil
			c.inUse = false
			fmt.Printf("removed socket %d\n", id)
			return
		}
	}
	fmt.Printf("Could not find socket %d\n", id)
}
