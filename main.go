package main

import (
	"fmt"
	"log"

	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/metadata"
	"github.com/DucNg/fip-player/player"
)

func main() {
	log.SetFlags(log.Lshortfile) // Enable line number on error

	metadata := metadata.FetchMetadata()
	fmt.Println(metadata.Now.FirstLine.Title + metadata.Now.SecondLine.Title)

	mpv := &player.MPV{}
	mpvChan, _ := mpv.Initialize()
	mpv.SendCommand([]string{"loadfile", "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance"})

	go dbus.RunDbusListener(mpv)

	for event := range mpvChan {
		switch event {
		case player.STATE_PLAYING:
			fmt.Println("Playing...")
		}
	}
}
