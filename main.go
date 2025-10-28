package main

import (
	_ "embed"
	"io"
	"log"
	"os"
	"os/exec"

	"fyne.io/systray"
	"github.com/BrunoCardosoFS/NaxiServer/server"
)

//go:embed icons/icon.ico
var iconData []byte

func main() {
	// Logs
	fileLog, errLog := os.OpenFile("./logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errLog != nil {
		log.Fatal(errLog)
	}
	defer fileLog.Close()

	multiLog := io.MultiWriter(os.Stdout, fileLog)
	log.SetOutput(multiLog)

	// Systray
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("NaxiServer")
	systray.SetTooltip("NaxiServer")

	mSettings := systray.AddMenuItem("Configurações", "Abrir Configurações")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Encerrar", "Encerrar servidor")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		for range mSettings.ClickedCh {
			var args = []string{"/c", "start", "http://localhost:8000/app/config.html"}

			if err := exec.Command("cmd", args...).Start(); err != nil {
				log.Fatalf("Error opening URL: %v", err)
			}
		}
	}()

	srv := server.NewServer()

	if err := srv.Start(":8000"); err != nil {
		log.Fatal(err)
	}
}

func onExit() {
	println("Encerrando servidor")
}
