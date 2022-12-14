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
		FirstLine  string
		SecondLine string
		Song       struct {
			Id   string
			Year uint
		}
		Cover struct {
			Src string
		}
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
	return time.Duration(fm.Media.EndTime-fm.Media.StartTime) * time.Microsecond
}

func (fm *FipMetadata) Delay() time.Duration {
	return time.Duration(fm.DelayToRefresh) * time.Millisecond
}
