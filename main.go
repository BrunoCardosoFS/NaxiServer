package main

import (
	_ "embed"
	"io"
	"log"
	"os"

	"fyne.io/systray"
	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/BrunoCardosoFS/NaxiServer/server"
)

//go:embed icons/icon_success.ico
var iconSuccessData []byte

// // go:embed icons/icon_failure.ico
// var iconFailureData []byte

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
	settings := getSettings()
	runSystray(settings.Port)

	database.InitDB("file:" + settings.DbPath + "naxistudio.db" + "?_journal_mode=WAL&_busy_timeout=5000")

	srv := server.NewServer(settings.DbPath)
	if err := srv.Start(settings.Port); err != nil {
		log.Fatal(err)
	}
}

func onExit() {
	database.DB.Close()

	log.Println("Shutting down server")
}
