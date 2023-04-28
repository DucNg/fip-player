package player

// #include <mpv/client.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -lmpv
//
// /* some helper functions for string arrays */
// char** makeCharArray(int size) {
//     return calloc(sizeof(char*), size);
// }
// void setArrayString(char** a, int i, char* s) {
//     a[i] = s;
// }
import "C"
import (
	"flag"
	"fmt"
	"log"
	"sync"
	"unsafe"
)

// MPV is an implementation of Backend, using libmpv.
type MPV struct {
	handle       *C.mpv_handle
	running      bool
	runningMutex sync.Mutex
	mainloopExit chan struct{}
	mute         bool
}

type State int

const (
	STATE_STOPPED   = 0
	STATE_PLAYING   = 1
	STATE_PAUSED    = 2
	STATE_BUFFERING = 3
	STATE_SEEKING   = 4 // not in the YouTube API
)

var logLibMPV = flag.Bool("log-libmpv", false, "log output of libmpv")

const propPause C.uint64_t = 0

// New creates a new MPV instance and initializes the libmpv player
func (mpv *MPV) Initialize() {
	if mpv.handle != nil || mpv.running {
		panic("already initialized")
	}

	mpv.mainloopExit = make(chan struct{})
	mpv.running = true
	mpv.mute = false
	mpv.handle = C.mpv_create()

	mpv.setOptionFlag("resume-playback", false)
	//mpv.setOptionString("softvol", "yes")
	//mpv.setOptionString("ao", "pulse")
	mpv.setOptionInt("volume", 100)
	mpv.setOptionInt("volume-max", 100)

	// Disable video in three ways.
	mpv.setOptionFlag("video", false)
	mpv.setOptionString("vo", "null")
	mpv.setOptionString("vid", "no")

	// Cache settings assume 128kbps audio stream (16kByte/s).
	// The default is a cache size of 25MB, these are somewhat more sensible
	// cache sizes IMO.
	mpv.setOptionInt("cache-secs", 10) // 10 seconds
	// mpv.setOptionInt("cache-seek-min", 16) // 1 second

	// Some extra debugging information, but don't read from stdin.
	// libmpv has a problem with signal handling, though: when `terminal` is
	// true, Ctrl+C doesn't work correctly anymore and program output is
	// disabled.
	mpv.setOptionFlag("terminal", *logLibMPV)
	mpv.setOptionFlag("input-terminal", false)
	mpv.setOptionFlag("quiet", true)

	mpv.checkError(C.mpv_initialize(mpv.handle))

	propName := C.CString("pause")
	C.mpv_observe_property(mpv.handle, propPause, propName, C.MPV_FORMAT_FLAG)
	C.free(unsafe.Pointer(propName))

	eventChan := make(chan State)

	go mpv.eventHandler(eventChan)
}

// setOptionFlag passes a boolean flag to mpv
func (mpv *MPV) setOptionFlag(key string, value bool) {
	cValue := C.int(0)
	if value {
		cValue = 1
	}

	mpv.setOption(key, C.MPV_FORMAT_FLAG, unsafe.Pointer(&cValue))
}

// setOptionInt passes an integer option to mpv
func (mpv *MPV) setOptionInt(key string, value int) {
	cValue := C.int64_t(value)
	mpv.setOption(key, C.MPV_FORMAT_INT64, unsafe.Pointer(&cValue))
}

// setOptionString passes a string option to mpv
func (mpv *MPV) setOptionString(key, value string) {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	mpv.setOption(key, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue))
}

// setOption is a generic function to pass options to mpv
func (mpv *MPV) setOption(key string, format C.mpv_format, value unsafe.Pointer) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	mpv.checkError(C.mpv_set_option(mpv.handle, cKey, format, value))
}

// setProperty sets the MPV player property
func (mpv *MPV) setProperty(name, value string) {
	log.Printf("MPV set property: %s=%s\n", name, value)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	// setProperty can take an unbounded time, don't block here using _async
	// TODO: use some form of error handling. Sometimes, it is impossible to
	// know beforehand whether setting a property will cause an error.
	// Importantly, catch the 'property unavailable' error.
	mpv.checkError(C.mpv_set_property_async(mpv.handle, 1, cName, C.MPV_FORMAT_STRING, unsafe.Pointer(&cValue)))
}

// sendCommand sends a command to the libmpv player
func (mpv *MPV) SendCommand(command []string) {
	// Print command, but without the stream
	cmd := make([]string, len(command))
	copy(cmd, command)
	if command[0] == "loadfile" {
		cmd[1] = "<stream>"
	}
	log.Println("MPV command:", cmd)

	cArray := C.makeCharArray(C.int(len(command) + 1))
	if cArray == nil {
		panic("got NULL from calloc")
	}
	defer C.free(unsafe.Pointer(cArray))

	for i, s := range command {
		cStr := C.CString(s)
		C.setArrayString(cArray, C.int(i), cStr)
		defer C.free(unsafe.Pointer(cStr))
	}

	mpv.checkError(C.mpv_command_async(mpv.handle, 0, cArray))
}

func (mpv *MPV) SetVolume(volume float64) {
	log.Printf("MPV set volume: %f\n", volume)

	cName := C.CString("volume")
	defer C.free(unsafe.Pointer(cName))

	mpv.checkError(C.mpv_set_property_async(mpv.handle, 1, cName, C.MPV_FORMAT_DOUBLE, unsafe.Pointer(&volume)))
}

func (mpv *MPV) Pause() {
	mpv.setProperty("pause", "yes")
}

func (mpv *MPV) ToggleMute() {
	mpv.mute = !mpv.mute
	mute := "no"
	if mpv.mute {
		mute = "yes"
	}
	mpv.setProperty("mute", mute)
}

// IsMute return true if mpv is muted
func (mpv *MPV) IsMute() bool { return mpv.mute }

func (mpv *MPV) Resume() {
	mpv.setProperty("pause", "no")
}

// playerEventHandler waits for libmpv player events and sends them on a channel
func (mpv *MPV) eventHandler(eventChan chan State) {
	for {
		// wait until there is an event (negative timeout means infinite timeout)
		// The timeout is 1 second to work around libmpv bug #1372 (mpv_wakeup
		// does not actually wake up mpv_wait_event). It keeps checking every
		// second whether MPV has exited.
		// TODO revert this as soon as the fix for that bug lands in a stable
		// release. Check for the problematic versions and keep the old behavior
		// for older MPV versions.
		event := C.mpv_wait_event(mpv.handle, 1)
		if event.event_id != C.MPV_EVENT_NONE {
			log.Printf("MPV event: %s (%d)\n", C.GoString(C.mpv_event_name(event.event_id)), int(event.event_id))
		}

		if event.error != 0 {
			panic("MPV API error")
		}

		mpv.runningMutex.Lock()
		running := mpv.running
		mpv.runningMutex.Unlock()

		if !running {
			close(eventChan)
			mpv.mainloopExit <- struct{}{}
			return
		}

		switch event.event_id {
		case C.MPV_EVENT_PLAYBACK_RESTART:
			eventChan <- STATE_PLAYING
		case C.MPV_EVENT_END_FILE:
			eventChan <- STATE_STOPPED
		case C.MPV_EVENT_PROPERTY_CHANGE:
			prop := (*C.mpv_event_property)(event.data)
			if event.reply_userdata == propPause {
				if *(*C.int)(prop.data) != 0 {
					eventChan <- STATE_PAUSED
				} else {
					eventChan <- STATE_PLAYING
				}
			}
		}
	}
}

// checkError checks for libmpv errors and panics if it finds one
func (mpv *MPV) checkError(status C.int) {
	if status < 0 {
		// this C string should not be freed (it is static)
		panic(fmt.Sprintf("mpv API error: %s (%d)", C.GoString(C.mpv_error_string(status)), int(status)))
	}
}
