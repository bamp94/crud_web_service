package main

import (
	"bufio"
	"fmt"
	"os"
)

// setGeneralChannel set create channel as general for all messages
func setGeneralChannel(channels ...chan interface{}) (chan interface{}, error) {
	switch len(channels) {
	case 0:
		return nil, fmt.Errorf("Channel list is empty")
	case 1:
		return channels[0], nil
	default:
		gChan := make(chan interface{})
		for i := 0; i < len(channels); i++ {
			go func(generalChannel, currentChannel chan interface{}) {
				for {
					msg, ok := <-currentChannel
					if !ok {
						fmt.Println("current channel closed")
						return
					}
					generalChannel <- msg
				}
			}(gChan, channels[i])
		}
		return gChan, nil
	}
}

func main() {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	channel, err := setGeneralChannel(ch1, ch2)
	if err != nil {
		panic("channel error: " + err.Error())
	}
	go func() {
		for {
			msg, ok := <-channel
			if !ok {
				panic("channel closed")
			}
			fmt.Println(msg)
		}
	}()
	go func() {
		for i := 0; i < 200; i++ {
			ch1 <- "test"
		}
	}()
	go func() {
		for i := 0; i < 200; i++ {
			ch2 <- "test2"
		}
	}()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
