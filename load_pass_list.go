package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

var data = [][2]string{
	{"A01", "Apple"},
	{"B02", "Banana"},
	{"C03", "Cherry"},
}

func SaveToFile(filename string, data [][2]string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("JSON encoding failed with error: %v", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write save file with error: %v", err)
	}

	log.Println("Successfully saved to file: " + filename)
}

func LoadJSONData(filename string) [][2]string {
	jsonContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file with error: %v", err)
	}
	var jsonData [][2]string
	err = json.Unmarshal(jsonContent, &jsonData)
	if err != nil {
		log.Fatalf("Failed to parse file with error: %v", err)
	}
	return jsonData
}

func ConvertSliceToListItem(data [][2]string) []list.Item {
	items := make([]list.Item, len(data))
	for i, r := range data {
		items[i] = item{title: r[0], id: r[1]}
	}
	return items
}
