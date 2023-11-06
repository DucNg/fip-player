package dbus

import (
	"log"
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
	conn  *dbus.Conn
}

func CreateDbusInstance(mpv *player.MPV) *Instance {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalln(err)
	}

	ins := &Instance{
		conn: conn,
	}
	mp2 := &MediaPlayer2{ins: ins, mpv: mpv}

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

	reply, err := conn.RequestName("org.mpris.MediaPlayer2.fipplayer", dbus.NameFlagReplaceExisting)
	if err != nil {
		log.Fatalln(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Fatalln("Name already taken")
	}

	return ins
}

func GetMetadataMap(fm *metadata.FipMetadata) MetadataMap {
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
		"mpris:artUrl":  fm.Now.Cover.Src,

		"xesam:album":          fm.Now.Song.Release.Title,
		"xesam:albumArtist":    []string{fm.Now.SecondLine},
		"xesam:artist":         []string{fm.Now.SecondLine},
		"xesam:contentCreated": fm.ContentCreated(),
		"xesam:title":          fm.Now.FirstLine,
	}

	return *m
}

func UpdateMetadata(ins *Instance, fm *metadata.FipMetadata) {
	metadata := GetMetadataMap(fm)

	dbusErr := ins.props.Set(
		"org.mpris.MediaPlayer2.Player",
		"Metadata",
		dbus.MakeVariant(metadata),
	)
	if dbusErr != nil {
		log.Println(dbusErr, metadata)
	}

	position := time.Duration(fm.Now.NowTime-fm.Media.StartTime) * time.Millisecond
	dbusErr = ins.props.Set(
		"org.mpris.MediaPlayer2.Player",
		"Position",
		dbus.MakeVariant(int64(position)),
	)
	if dbusErr != nil {
		log.Println(dbusErr)
	}
}

func IncrementPosition(ins *Instance) {
	previousPosition, dbusErr := ins.props.Get(
		"org.mpris.MediaPlayer2.Player",
		"Position",
	)
	if dbusErr != nil {
		log.Println(dbusErr)
	}

	newPosition := previousPosition.Value().(int64) + 1e+6

	dbusErr = ins.props.Set(
		"org.mpris.MediaPlayer2.Player",
		"Position",
		dbus.MakeVariant(int64(newPosition)),
	)
	if dbusErr != nil {
		log.Println(dbusErr)
	}
}

func (ins *Instance) CloseConnection() {
	ins.conn.Close()
}
