package dbus

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/DucNg/fip-player/metadata"
	"github.com/DucNg/fip-player/player"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type MetadataMap map[string]interface{}

type Instance struct {
	props *prop.Properties

	name string
}

func RunDbusListener(mpv *player.MPV) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	ins := &Instance{
		name: fmt.Sprintf("org.mpris.MediaPlayer2.fipPlayer.instance%d", os.Getpid()),
	}
	mp2 := &MediaPlayer2{ins: ins, mpv: mpv}

	err = conn.Export(mp2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2")
	if err != nil {
		log.Fatalln(err)
	}

	err = conn.Export(mp2, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")
	if err != nil {
		log.Fatalln(err)
	}

	err = conn.Export(introspect.NewIntrospectable(IntrospectNode()), "/org/mpris/MediaPlayer2", "org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Fatalln(err)
	}

	ins.props, err = prop.Export(conn, "/org/mpris/MediaPlayer2", map[string]map[string]*prop.Prop{
		"org.mpris.MediaPlayer2":        mp2.properties(),
		"org.mpris.MediaPlayer2.Player": mp2.playerProps(),
	})
	if err != nil {
		log.Fatalln(err)
	}

	reply, err := conn.RequestName("org.mpris.MediaPlayer2.fip-player", dbus.NameFlagReplaceExisting)
	if err != nil {
		log.Fatalln(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalln("Name already taken")
	}
	fmt.Println("D-Bus listening")

	go func() {
		UpdateMetadata(ins)
	}()

	select {}
}

func GetMetadata() (MetadataMap, time.Duration) {
	fm := metadata.FetchMetadata()

	id := strings.ReplaceAll(fm.Now.Song.Id, "-", "")

	var trackId string
	if id == "" {
		trackId = "/org/mpris/MediaPlayer2/TrackList/NoTrack"
	} else {
		trackId = "/org/mpris/MediaPlayer2/" + id
	}

	m := &MetadataMap{
		"mpris:trackid": dbus.ObjectPath(trackId),
		"mpris:length":  fm.Duration(),

		"xesam:title":       fm.Now.FirstLine.Title,
		"xesam:artist":      fm.Now.SecondLine.Title,
		"xesam:albumArtist": fm.Now.SecondLine.Title,
	}

	return *m, fm.Delay()
}

func UpdateMetadata(ins *Instance) {
	metadata, delayToRefresh := GetMetadata()

	dbusErr := ins.props.Set("org.mpris.MediaPlayer2.Player", "Metadata", dbus.MakeVariant(metadata))
	if dbusErr != nil {
		log.Println(dbusErr, metadata)
	}

	time.Sleep(delayToRefresh)
	UpdateMetadata(ins)
}
