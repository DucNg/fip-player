package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type FipMetadata struct {
	DelayToRefresh uint
	Now            struct {
		FirstLine struct  {
			Title string
		}
		SecondLine struct {
			Title string
		}
		Song       struct {
			Id      string
			Year    uint
			Release struct {
				Title     string
				Label     string
				Reference string
			}
		}
		Cover struct {
			Src string
		}
		NowTime uint
	}
	Media struct {
		StartTime uint
		EndTime   uint
	}
}

func FetchMetadata(url string) *FipMetadata {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer res.Body.Close()

	jsonRes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var metadata FipMetadata
	err = json.Unmarshal(jsonRes, &metadata)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Trying to debug FIP API sending strange data
	log.Printf("Delay to refresh %d\n", metadata.DelayToRefresh)
	if metadata.DelayToRefresh >= 1000000 {
		log.Println(string(jsonRes))
	}

	return &metadata
}

func (fm *FipMetadata) Duration() time.Duration {
	return time.Duration(fm.Media.EndTime-fm.Media.StartTime) * time.Millisecond
}

func (fm *FipMetadata) Delay() time.Duration {
	return time.Duration(fm.DelayToRefresh) * time.Millisecond
}

func (fm *FipMetadata) ContentCreated() string {
	// Between songs, sometimes, API sends a generic name without year
	if fm.Now.Song.Year == 0 {
		return ""
	}

	parsedYear, err := time.Parse("2006", fmt.Sprintf("%v", fm.Now.Song.Year))
	if err != nil {
		log.Println(err)
	}

	return parsedYear.Format(time.RFC3339)
}

func (fm *FipMetadata) ProgressPercent() float64 {
	// Between songs, sometimes, API sends a generic name without startTime/endTime
	// To get a more accurate progress bar in this case, progress needs to be set to 0%
	if fm.Media.EndTime == 0 || fm.Media.StartTime == 0 {
		return 0
	}

	duration := fm.Media.EndTime - fm.Media.StartTime
	position := fm.Now.NowTime - fm.Media.StartTime

	progress := float64(position) / float64(duration)

	return progress
}

func (fm *FipMetadata) ValueOfOneSecond() float64 {
	var duration float64

	// Between songs, sometimes, API sends a generic name without startTime/endTime
	if fm.Media.EndTime == 0 || fm.Media.StartTime == 0 {
		duration = float64(fm.DelayToRefresh) / 1000
	} else {
		duration = float64(fm.Media.EndTime - fm.Media.StartTime)
	}

	return (100 / duration) / 100
}
