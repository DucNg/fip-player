package main

import (
	"fmt"
)

func main() {
	metadata := fetchMetadata()
	fmt.Println(metadata.Now.FirstLine.Title + metadata.Now.SecondLine.Title)

	mpv := &MPV{}
	mpvChan, _ := mpv.initialize()
	mpv.sendCommand([]string{"loadfile", "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance"})

	for event := range mpvChan {
		switch event {
		case STATE_PLAYING:
			fmt.Println("Playing...")
		}
	}
}
