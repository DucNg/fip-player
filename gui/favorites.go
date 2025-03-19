package gui

import (
	"fmt"
	"log"

	"github.com/DucNg/fip-player/cache"
	"github.com/charmbracelet/bubbles/list"
)

// filterFavoriteRadios returns a new slice containing only radios that are marked as favorites
func filterFavoriteRadios(items []list.Item) []list.Item {
	var favorites []list.Item

	favoritesIDs := cache.GetFavoritesRadioIDs()
	if favoritesIDs == nil {
		log.Println("failed to read favorites from file")
		return nil
	}

	for _, itm := range items {
		for _, favoritesID := range favoritesIDs {
			if itm.(item).id == favoritesID {
				favorites = append(favorites, itm)
			}
		}
	}

	return favorites
}

// toggleFavoriteList filters the list to show only favorite radios
func (m *model) toggleFavoriteList() {
	if m.isFavoriteModeEnabled {
		m.list.SetItems(getRadiosWithIDs())

		m.isFavoriteModeEnabled = false
	} else {
		// Get the current list items
		allItems := m.list.Items()

		// Filter to get only favorites
		favoriteItems := filterFavoriteRadios(allItems)

		// If there are no favorites, don't change the view
		// if len(favoriteItems) == 0 {
		// 	return
		// }

		// Update the list with only favorites
		m.list.SetItems(favoriteItems)

		m.isFavoriteModeEnabled = true
	}
}

func (m *model) toggleSelectedItemAsFavorite() {
	itemToToggle := m.list.SelectedItem().(item)

	if itemToToggle.favorite {
		itemToToggle.favorite = false
		cache.RemoveRadioIDFromFavorites(itemToToggle.id)
		itemToToggle.title = itemToToggle.title[:len(itemToToggle.title)-len(" ❤️")]
	} else {
		itemToToggle.favorite = true
		cache.AddRadioIDToFavorites(itemToToggle.id)
		itemToToggle.title = fmt.Sprintf("%v ❤️", itemToToggle.title)
	}

	m.list.SetItem(m.list.Index(), itemToToggle)
}
