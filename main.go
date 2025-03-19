package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/DucNg/fip-player/cache"
	"github.com/DucNg/fip-player/dbus"
	"github.com/DucNg/fip-player/gui"
	"github.com/DucNg/fip-player/player"
)

func main() {
	logFile := enableDebugLogs()
	if logFile != nil {
		defer logFile.Close()
	}

	mpv := &player.MPV{}
	mpv.Initialize()

	ins := dbus.CreateDbusInstance(mpv)
	defer ins.CloseConnection()

	IDOnClose := gui.Render(ins, mpv, cache.GetLastRadioID())
	cache.WriteLastRadioID(IDOnClose)
}

func enableDebugLogs() *os.File {
	debug := flag.Bool("d", false, "enables debug output to /tmp")
	flag.Parse()

	if !*debug {
		log.SetOutput(io.Discard)
		return nil
	}

	log.SetFlags(log.Lshortfile) // Enable line number on error

	logFile, err := os.CreateTemp("", "fip-player-log")
	if err != nil {
		log.Fatalln(err)
	}

	log.SetOutput(logFile)

	return logFile
}
