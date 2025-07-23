package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "embed"
)

//go:embed fallback.json
var FallbackMetadata []byte

type FipMetadata struct {
	DelayToRefresh uint
	Now            struct {
		StartTime uint
		EndTime   uint
		FirstLine struct {
			Title string
		}
		SecondLine struct {
			Title string
		}
		Song struct {
			Id      string
			Year    uint
			Release struct {
				Title     string
				Label     string
				Reference string
			}
		}
		Visuals struct {
			Card struct {
				Src string
			}
		}
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
		// API sent strange data, fallback to safe values and retry in 3 seconds
		log.Printf("Error unmarshalling metadata, falling back to safe values, API payload: %v", string(jsonRes))

		err = json.Unmarshal(FallbackMetadata, &metadata)
		if err != nil {
			log.Fatalf("error unmarshalling fallback metadata: %v", err.Error())
		}
	}

	return &metadata
}

func (fm *FipMetadata) Duration() time.Duration {
	return time.Duration(fm.Now.EndTime-fm.Now.StartTime) * time.Millisecond
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
	if fm.Now.EndTime == 0 || fm.Now.StartTime == 0 {
		return 0
	}

	duration := fm.Now.EndTime - fm.Now.StartTime
	position := uint(time.Now().UnixMilli()) - fm.Now.StartTime

	progress := float64(position) / float64(duration)

	return progress
}

func (fm *FipMetadata) ValueOfOneSecond() float64 {
	var duration float64

	// Between songs, sometimes, API sends a generic name without startTime/endTime
	if fm.Now.EndTime == 0 || fm.Now.StartTime == 0 {
		duration = float64(fm.DelayToRefresh) / 1000
	} else {
		duration = float64(fm.Now.EndTime - fm.Now.StartTime)
	}

	return (100 / duration) / 100
}
