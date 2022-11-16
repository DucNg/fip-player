package metadata

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type FipMetadata struct {
	DelayToRefresh uint
	Now            struct {
		FirstLine struct {
			Title string
		}
		SecondLine struct {
			Title string
		}
		Song struct {
			Id   string
			Year uint
		}
		StartTime uint
		EndTime   uint
	}
}

func FetchMetadata() FipMetadata {
	res, err := http.Get("https://www.radiofrance.fr/api/v2.0/stations/fip/live")
	if err != nil {
		log.Fatalln(err.Error())
	}
	jsonRes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var metatdata FipMetadata
	err = json.Unmarshal(jsonRes, &metatdata)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return metatdata
}

func (fm *FipMetadata) Duration() time.Duration {
	return time.Duration(fm.Now.EndTime-fm.Now.StartTime) * time.Microsecond
}

func (fm *FipMetadata) Delay() time.Duration {
	return time.Duration(fm.DelayToRefresh) * time.Millisecond
}
