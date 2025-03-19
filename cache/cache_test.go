package cache

import (
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// TestWriteAndGetLastRadioID tests the WriteLastRadioID and GetLastRadioID functions
func TestWriteAndGetLastRadioID(t *testing.T) {
	// Test data
	testID := 42

	// Write the last radio ID using the function from cache.go
	WriteLastRadioID(testID)

	// Verify file was created with correct content
	lastIDPath := lastRadioIDPath()
	content, err := os.ReadFile(lastIDPath)
	if err != nil {
		t.Fatalf("Failed to read last radio ID file: %v", err)
	}

	if string(content) != strconv.Itoa(testID) {
		t.Errorf("Expected file content to be %q, got %q", strconv.Itoa(testID), string(content))
	}

	// Read the ID using the function from cache.go
	readID := GetLastRadioID()
	if readID != testID {
		t.Errorf("Expected GetLastRadioID() to return %d, got %d", testID, readID)
	}
}

// TestAddRadioIDToFavorites tests the AddRadioIDToFavorites function
func TestAddRadioIDToFavorites(t *testing.T) {
	// Clear favorites file before testing
	favPath := favoritesIDPath()
	os.Remove(favPath)

	// Test data
	testID := 123

	// Add radio ID to favorites using the function from cache.go
	AddRadioIDToFavorites(testID)

	// Verify file content
	content, err := os.ReadFile(favPath)
	if err != nil {
		t.Fatalf("Failed to read favorites file: %v", err)
	}

	if !strings.Contains(string(content), strconv.Itoa(testID)) {
		t.Errorf("Expected favorites file to contain %d, got %q", testID, string(content))
	}

	// Read favorites back to verify
	favorites := GetFavoritesRadioIDs()
	found := false
	for _, id := range favorites {
		if id == testID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected GetFavoritesRadioIDs() to contain %d, got %v", testID, favorites)
	}
}

// TestRemoveRadioIDFromFavorites tests the RemoveRadioIDFromFavorites function
func TestRemoveRadioIDFromFavorites(t *testing.T) {
	// Reset favorites file
	favPath := favoritesIDPath()
	os.Remove(favPath)

	// Setup initial content with multiple IDs
	ids := []int{123, 456, 789}
	idToRemove := 456

	// Add all IDs to favorites
	for _, id := range ids {
		AddRadioIDToFavorites(id)
	}

	// Make sure IDs were added
	initialFavorites := GetFavoritesRadioIDs()
	if len(initialFavorites) != len(ids) {
		t.Fatalf("Failed to set up test: expected %d favorites, got %d",
			len(ids), len(initialFavorites))
	}

	// Remove one ID
	RemoveRadioIDFromFavorites(idToRemove)

	// Verify the ID was removed
	updatedFavorites := GetFavoritesRadioIDs()

	// Should have one less ID
	if len(updatedFavorites) != len(ids)-1 {
		t.Errorf("Expected %d favorites after removal, got %d",
			len(ids)-1, len(updatedFavorites))
	}

	// The removed ID should not be present
	for _, id := range updatedFavorites {
		if id == idToRemove {
			t.Errorf("Found removed ID %d in favorites after removal", idToRemove)
		}
	}

	// Check file content directly
	content, err := os.ReadFile(favPath)
	if err != nil {
		t.Fatalf("Failed to read favorites file: %v", err)
	}

	if strings.Contains(string(content), strconv.Itoa(idToRemove)) {
		t.Errorf("Found removed ID %d in favorites file after removal", idToRemove)
	}
}

// TestGetFavoritesRadioIDs tests the GetFavoritesRadioIDs function
func TestGetFavoritesRadioIDs(t *testing.T) {
	// Reset favorites file
	favPath := favoritesIDPath()
	os.Remove(favPath)

	// Create file with known content
	ids := []int{101, 202, 303}
	file, err := os.Create(favPath)
	if err != nil {
		t.Fatalf("Failed to create favorites file: %v", err)
	}

	for _, id := range ids {
		file.WriteString(strconv.Itoa(id) + "\n")
	}
	file.Close()

	// Get favorites using the function
	favorites := GetFavoritesRadioIDs()

	// Verify count
	if len(favorites) != len(ids) {
		t.Errorf("Expected %d favorites, got %d", len(ids), len(favorites))
	}

	// Verify all IDs are present
	for _, expectedID := range ids {
		found := false
		for _, actualID := range favorites {
			if actualID == expectedID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find ID %d in favorites, but it was missing", expectedID)
		}
	}
}

// TestGetLastRadioID_NoFile tests GetLastRadioID when file doesn't exist
func TestGetLastRadioID_NoFile(t *testing.T) {
	// Remove the file if it exists
	lastIDPath := lastRadioIDPath()
	os.Remove(lastIDPath)

	// Test getting last radio ID when file doesn't exist
	id := GetLastRadioID()
	if id != 0 {
		t.Errorf("Expected default ID 0 when file doesn't exist, got %d", id)
	}
}

// TestGetFavoritesRadioIDs_NoFile tests GetFavoritesRadioIDs when file doesn't exist
func TestGetFavoritesRadioIDs_NoFile(t *testing.T) {
	// Remove the file if it exists
	favPath := favoritesIDPath()
	os.Remove(favPath)

	// Test getting favorites when file doesn't exist
	favorites := GetFavoritesRadioIDs()
	if favorites != nil && len(favorites) != 0 {
		t.Errorf("Expected empty favorites list when file doesn't exist, got %v", favorites)
	}
}
