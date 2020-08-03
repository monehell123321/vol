// Spam just handles spam.
package main

import (
	"fmt"
	"strconv"
	"time"
)

var quit = make(chan int)

type SpamInstance struct {
	Channel []struct {
		ID    string
		Delay int
	}
}

func (spam *SpamInstance) Invoke() {
	for index, channel := range spam.Channel {
		fmt.Println("[SPAM] Starting goroutine " + strconv.FormatInt(int64(index), 10) + " for channel ID: " + channel.ID)
		go spam.doSpam(channel.ID, channel.Delay)
	}
}

/* func (spam *SpamInstance) Quit() {
	fmt.Println("[SPAM] Quitting spam routine for channel ID: " + spam.ChannelID)
	close(quit)
} */ // Not working yet

func (spam *SpamInstance) doSpam(id string, delay int) {
	for {
		select {
		case <-quit:
			return
		default:
			client.ChannelMessageSend(id, "test")
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}
}
