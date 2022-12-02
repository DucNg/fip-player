package main

import (
	"fmt"
	"log"

	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/player"
)

func main() {
	log.SetFlags(log.Lshortfile) // Enable line number on error

	mpv := &player.MPV{}
	mpvChan, _ := mpv.Initialize()
	mpv.SendCommand([]string{"loadfile", "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance"})

	go dbus.RunDbusListener(mpv)

	for event := range mpvChan {
		switch event {
		case player.STATE_PLAYING:
			fmt.Println("Playing...")
		case player.STATE_PAUSED:
			fmt.Println("Pausing...")
		}
	}
}
