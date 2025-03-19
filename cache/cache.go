package cache

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

func GetLastRadioID() int {
	indexBytes, err := os.ReadFile(lastRadioIDPath())
	if err != nil {
		return 0
	}

	index, err := strconv.Atoi(string(indexBytes))
	if err != nil {
		return 0
	}

	return index
}

func WriteLastRadioID(IDOnClose int) {
	err := os.WriteFile(lastRadioIDPath(), []byte(fmt.Sprintf("%v", IDOnClose)), 0666)
	if err != nil {
		log.Printf("failed to write last buffer at %q: %s\n", lastRadioIDPath(), err)
	}
}

func GetFavoritesRadioIDs() []int {
	favoritesFile, err := os.Open(favoritesIDPath())
	if err != nil {
		log.Printf("Failed to open favorites file: %v\n", err)
		return nil
	}
	defer favoritesFile.Close()

	scanner := bufio.NewScanner(favoritesFile)

	favoritesRadioIndexes := []int{}
	for scanner.Scan() {
		radioIndex, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("invalid radio ID: %v\n", err)
			continue
		}

		favoritesRadioIndexes = append(favoritesRadioIndexes, radioIndex)
	}

	return favoritesRadioIndexes
}

func AddRadioIDToFavorites(ID int) {
	favoritesFile, err := os.OpenFile(favoritesIDPath(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Failed to open favorites file: %v\n", err)
		return
	}

	defer favoritesFile.Close()

	if _, err = favoritesFile.WriteString(strconv.Itoa(ID) + "\n"); err != nil {
		log.Printf("Failed to write favorites file: %v\n", err)
		return
	}
}

func RemoveRadioIDFromFavorites(ID int) {
	favoritesFile, err := os.Open(favoritesIDPath())
	if err != nil {
		log.Printf("Failed to write favorites file: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(favoritesFile)

	favoritesRadioIndexes := []int{}
	for scanner.Scan() {
		radioIndex, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("invalid radio ID: %v\n", err)
			continue
		}

		if radioIndex != ID {
			favoritesRadioIndexes = append(favoritesRadioIndexes, radioIndex)
		}
	}

	favoritesFile.Close()

	err = os.Remove(favoritesIDPath())
	if err != nil {
		log.Printf("Failed to remove favorites file: %v\n", err)
		return
	}

	favoritesFile, err = os.Create(favoritesIDPath())
	if err != nil {
		log.Printf("Failed to write favorites file: %v\n", err)
		return
	}
	defer favoritesFile.Close()

	for _, radioIndex := range favoritesRadioIndexes {
		favoritesFile.WriteString(strconv.Itoa(radioIndex) + "\n")
	}
}

func cachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	cache := path.Join(cacheDir, "fip-player")
	err = os.MkdirAll(cache, 0755)
	if err != nil {
		panic(err)
	}
	return cache
}

func lastRadioIDPath() string {
	return path.Join(cachePath(), "lastradioindex.txt")
}

func favoritesIDPath() string {
	return path.Join(cachePath(), "favoritesindex.txt")
}
