package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strconv"

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

	indexOnClose := gui.Render(ins, mpv, getLastRadioIndex())

	err := os.WriteFile(lastRadioIndexPath(), []byte(fmt.Sprintf("%v", indexOnClose)), 0666)
	if err != nil {
		log.Printf("failed to write last buffer at %q: %s\n", lastRadioIndexPath(), err)
	}
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

func cachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cache := path.Join(cacheDir, "fip-radio")
	err = os.MkdirAll(cache, 0755)
	if err != nil {
		panic(err)
	}
	return cache
}

func lastRadioIndexPath() string {
	return path.Join(cachePath(), "lastradioindex.txt")
}

func getLastRadioIndex() int {
	indexBytes, err := os.ReadFile(lastRadioIndexPath())
	if err != nil {
		return 0
	}

	index, err := strconv.Atoi(string(indexBytes))
	if err != nil {
		return 0
	}

	return index
}
