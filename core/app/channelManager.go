package app

import (
	"math/rand"
	"sync"
	"time"
)

var cancelChannels channelSync

type channelSync struct {
	sync.RWMutex
	Channels map[string](chan int)
}

func AddCancelChannel() (key string, c chan int) {
	cancelChannels.Lock()

	if cancelChannels.Channels == nil {
		cancelChannels.Channels = make(map[string](chan int))
	}

	c = make(chan int)
	key = RandomString(20)
	cancelChannels.Channels[key] = c
	cancelChannels.Unlock()

	return
}

func CancelChannel(key string) {
	cancelChannels.Lock()

	if cancelChannels.Channels == nil {
		return
	}

	cancelChannels.Channels[key] <- 0
	delete(cancelChannels.Channels, key)
	cancelChannels.Unlock()

	return
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
