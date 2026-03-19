package metadata

import (
	"encoding/json"
	"fmt"
)

// fieldIdx returns the integer index stored at obj[field], or -1 if missing/null.
func fieldIdx(obj map[string]json.RawMessage, field string) int {
	v, ok := obj[field]
	if !ok {
		return -1
	}
	var idx int
	if err := json.Unmarshal(v, &idx); err != nil {
		return -1
	}
	return idx
}

// resolveStr returns the string at arr[idx], or "" if out of bounds or null.
func resolveStr(arr []json.RawMessage, idx int) string {
	if idx < 0 || idx >= len(arr) {
		return ""
	}
	var s string
	if err := json.Unmarshal(arr[idx], &s); err != nil {
		return ""
	}
	return s
}

// resolveUint returns the uint at arr[idx], or 0 if out of bounds or not a number.
func resolveUint(arr []json.RawMessage, idx int) uint {
	if idx < 0 || idx >= len(arr) {
		return 0
	}
	var n uint
	if err := json.Unmarshal(arr[idx], &n); err != nil {
		return 0
	}
	return n
}

// resolveObj returns arr[idx] parsed as a JSON object, or nil on failure.
func resolveObj(arr []json.RawMessage, idx int) map[string]json.RawMessage {
	if idx < 0 || idx >= len(arr) {
		return nil
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(arr[idx], &obj); err != nil {
		return nil
	}
	return obj
}

func parseAPIResponse(data []byte) (*FipMetadata, error) {
	var outer struct {
		Type   string `json:"type"`
		Result string `json:"result"`
	}
	if err := json.Unmarshal(data, &outer); err != nil {
		return nil, fmt.Errorf("unmarshal outer: %w", err)
	}

	var arr []json.RawMessage
	if err := json.Unmarshal([]byte(outer.Result), &arr); err != nil {
		return nil, fmt.Errorf("unmarshal inner array: %w", err)
	}

	root := resolveObj(arr, 0)
	if root == nil {
		return nil, fmt.Errorf("root element is not an object")
	}

	var fm FipMetadata

	fm.DelayToRefresh = resolveUint(arr, fieldIdx(root, "delayToRefresh"))

	now := resolveObj(arr, fieldIdx(root, "now"))
	if now == nil {
		return nil, fmt.Errorf("now element is not an object")
	}

	fm.Now.StartTime = resolveUint(arr, fieldIdx(now, "startTime"))
	fm.Now.EndTime = resolveUint(arr, fieldIdx(now, "endTime"))
	fm.Now.Title = resolveStr(arr, fieldIdx(now, "firstLine"))
	fm.Now.Artist = resolveStr(arr, fieldIdx(now, "secondLine"))

	cover := resolveObj(arr, fieldIdx(now, "cover"))
	if cover != nil {
		fm.Now.Cover.Id = resolveStr(arr, fieldIdx(cover, "id"))
	}

	song := resolveObj(arr, fieldIdx(now, "song"))
	if song != nil {
		fm.Now.Song.Id = resolveStr(arr, fieldIdx(song, "id"))
		fm.Now.Song.Year = resolveUint(arr, fieldIdx(song, "year"))

		release := resolveObj(arr, fieldIdx(song, "release"))
		if release != nil {
			fm.Now.Song.Release.Title = resolveStr(arr, fieldIdx(release, "title"))
			fm.Now.Song.Release.Label = resolveStr(arr, fieldIdx(release, "label"))
			fm.Now.Song.Release.Reference = resolveStr(arr, fieldIdx(release, "reference"))
		}
	}

	return &fm, nil
}

func fallbackFipMetadata() *FipMetadata {
	var fm FipMetadata
	fm.DelayToRefresh = 60000
	fm.Now.Title = "Error reading metadata"
	fm.Now.Artist = "Error reading metadata"
	return &fm
}
