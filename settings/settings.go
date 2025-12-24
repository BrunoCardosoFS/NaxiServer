package settings

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

var (
	instance *models.Settings
	once     sync.Once
)

func GenerateSecret(n int) (string, error) {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GetSettings() *models.Settings {
	once.Do(func() {
		instance = loadSettingsFromFile()
	})

	return instance
}

func loadSettingsFromFile() *models.Settings {
	setting := &models.Settings{}

	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	settingsPath := filepath.Join(exPath, "settings.json")

	fileSettings, err := os.Open(settingsPath)
	if err != nil {
		token, errGen := GenerateSecret(32)
		if errGen != nil {
			log.Println("Settings: Error generating token.")
			token = "NaxiServerSecretKey"
		}

		setting.DbPath = filepath.Join(exPath, "../db/")
		setting.Port = ":8000"
		setting.JwtKey = token

		if os.IsNotExist(err) {
			log.Println("settings.json file not found. Creating default...")

			newFile, createErr := os.Create(settingsPath)
			if createErr != nil {
				log.Panicf("Failed to create settings.json: %v", createErr)
			}
			defer newFile.Close()

			encoder := json.NewEncoder(newFile)
			if encodeErr := encoder.Encode(setting); encodeErr != nil {
				log.Panicf("Failed to write to settings.json: %v", encodeErr)
			}

			log.Println("settings.json created successfully.")

			return setting
		}

		log.Println(err)
		return setting
	}
	defer fileSettings.Close()

	err = json.NewDecoder(fileSettings).Decode(&setting)
	if err != nil {
		log.Panic(err)
	}

	return setting
}
