package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
)

var data = [][2]string{
	{"A01", "Apple"},
	{"B02", "Banana"},
	{"C03", "Cherry"},
}

func encrypt(data []byte, encKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data []byte, encKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	decData, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return decData, nil
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func SaveToFile(filename string, data [][3]string, encKey []byte) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	encData, err := encrypt(jsonData, encKey)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, encData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func LoadJSONData(filename string) [][3]string {
	jsonContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file with error: %v", err)
	}
	var jsonData [][3]string
	err = json.Unmarshal(jsonContent, &jsonData)
	if err != nil {
		log.Fatalf("Failed to parse file with error: %v", err)
	}
	return jsonData
}

func LoadEncryptedData(filename string, decKey []byte) ([][3]string, error) {
	encData, readErr := os.ReadFile(filename)
	if readErr != nil {
		return nil, readErr
	}
	jsonData, decErr := decrypt(encData, decKey)
	if decErr != nil {
		return nil, decErr
	}
	var content [][3]string
	parsingErr := json.Unmarshal(jsonData, &content)
	if parsingErr != nil {
		return nil, parsingErr
	}
	return content, nil
}

func ConvertSliceToListItem(data [][3]string) []list.Item {
	items := make([]list.Item, len(data))
	for i, r := range data {
		pw := ""
		if len(r) > 2 {
			pw = r[2]
		}
		items[i] = item{title: "Title: "+r[0], id: "ID: "+r[1], password: pw}
	}
	return items
}
