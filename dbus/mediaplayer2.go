package dbus

import (
	"errors"

	"github.com/DucNg/fip-player/player"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

type MediaPlayer2 struct {
	*Instance

	mpv *player.MPV
}

func (m *MediaPlayer2) properties() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"CanQuit":      newProp(false, nil),        // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:CanQuit
		"CanRaise":     newProp(false, nil),        // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:CanRaise
		"HasTrackList": newProp(false, nil),        // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:HasTrackList
		"Identity":     newProp("fip-player", nil), // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:Identity
		// Empty because we can't add arbitary files in...
		"SupportedUriSchemes": newProp([]string{}, nil), // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:SupportedUriSchemes
		"SupportedMimeTypes":  newProp([]string{}, nil), // https://specifications.freedesktop.org/mpris-spec/latest/Media_Player.html#Property:SupportedMimeTypes
	}
}

func (m *MediaPlayer2) playerProps() map[string]*prop.Prop {
	return map[string]*prop.Prop{
		"PlaybackStatus": newProp("Playing", nil),
		"Rate":           newProp(1.0, notImplemented),
		"Metadata":       newProp(MetadataMap{}, nil),
		"Volume":         newProp(float64(100), nil),
		"Position":       newProp(int64(0), nil),
		"MinimumRate":    newProp(1.0, nil),
		"MaximumRate":    newProp(1.0, nil),
		"CanGoNext":      newProp(false, nil),
		"CanGoPrevious":  newProp(false, nil),
		"CanPlay":        newProp(true, nil),
		"CanPause":       newProp(false, nil),
		"CanSeek":        newProp(false, nil),
		"CanControl":     newProp(false, nil),
	}
}

// Creates a new property.
func newProp(value interface{}, cb func(*prop.Change) *dbus.Error) *prop.Prop {
	return &prop.Prop{
		Value:    value,
		Writable: true,
		Emit:     prop.EmitTrue,
		Callback: cb,
	}
}

func notImplemented(c *prop.Change) *dbus.Error {
	return dbus.MakeFailedError(errors.New("not implemented"))
}

// TimeInUs is time in microseconds.
// https://specifications.freedesktop.org/mpris-spec/latest/Player_Interface.html#Simple-Type:Time_In_Us
type TimeInUs int64
