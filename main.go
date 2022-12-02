package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/gui"
	"github.com/DucNg/fip-player/player"
)

func main() {
	enableDebugLogs()

	mpv := &player.MPV{}
	mpv.Initialize()
	mpv.SendCommand([]string{"loadfile", "https://stream.radiofrance.fr/fip/fip.m3u8?id=radiofrance"})

	go dbus.RunDbusListener(mpv)

	gui.Render()

	// for event := range mpvChan {
	// 	switch event {
	// 	case player.STATE_PLAYING:
	// 		fmt.Println("Playing...")
	// 	case player.STATE_PAUSED:
	// 		fmt.Println("Pausing...")
	// 	}
	// }
}

func enableDebugLogs() {
	debug := flag.Bool("d", false, "enables debug output to /tmp")
	flag.Parse()

	if !*debug {
		log.SetOutput(io.Discard)
		return
	}

	log.SetFlags(log.Lshortfile) // Enable line number on error

	logFile, err := os.CreateTemp("", "fip-player-log")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logFile)
}
