package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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
	}
}

func fetchMetadata() FipMetadata {
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
