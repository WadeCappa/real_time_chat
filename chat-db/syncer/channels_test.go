package syncer

import (
	"log"
	"testing"
)

func TestGetChannels(t *testing.T) {
	c, err := getChannels("localhost:50055")
	if err != nil {
		t.Errorf("failed: %v", err)
	}
	for channel := range c {
		log.Println(channel)
	}
}
