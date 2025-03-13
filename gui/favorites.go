package gui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

// filterFavoriteRadios returns a new slice containing only radios that are marked as favorites
func filterFavoriteRadios(items []list.Item) []list.Item {
	var favorites []list.Item

	for _, itm := range items {
		if itm.(item).favorite {
			favorites = append(favorites, itm)
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
		itemToToggle.title = itemToToggle.title[:len(itemToToggle.title)-len(" ❤️")]
	} else {
		itemToToggle.favorite = true
		itemToToggle.title = fmt.Sprintf("%v ❤️", itemToToggle.title)
	}

	m.list.SetItem(m.list.Index(), itemToToggle)
}
