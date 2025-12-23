package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"fyne.io/systray"
	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func runSystray(port string) {
	systray.SetIcon(iconSuccessData)
	systray.SetTitle("NaxiServer")
	systray.SetTooltip("NaxiServer Online - localhost" + port)

	systray.AddMenuItem("localhost"+port, "Servidor online")
	systray.AddSeparator()

	mSettings := systray.AddMenuItem("Configurações", "Abrir Configurações")
	mLicense := systray.AddMenuItem("Sobre", "Sobre o NaxiStudio Server")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Encerrar", "Encerrar servidor")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for range mSettings.ClickedCh {
			var args = []string{"/c", "start", "http://localhost" + port + "/settings/"}

			if err := exec.Command("cmd", args...).Start(); err != nil {
				log.Fatalf("Error opening URL: %v", err)
			}
		}
	}()

	go func() {
		for range mLicense.ClickedCh {
			var args = []string{"/c", "start", "http://localhost" + port + "/about/"}

			if err := exec.Command("cmd", args...).Start(); err != nil {
				log.Fatalf("Error opening URL: %v", err)
			}
		}
	}()
}

func getSettings() models.Settings {
	settingsPath := "./settings.json"
	var setting models.Settings

	fileSettings, err := os.Open(settingsPath)
	if err != nil {
		setting = models.Settings{
			DbPath: "../db/",
			Port:   ":8000",
		}

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
