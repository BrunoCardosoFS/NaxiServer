package main

import (
	"log"
	"os/exec"

	"fyne.io/systray"
	"github.com/BrunoCardosoFS/NaxiServer/settings"
)

func runSystray() {
	port := settings.GetSettings().Port

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
