package main

import (
	_ "embed"
	"io"
	"log"
	"os"
	"path/filepath"

	"fyne.io/systray"
	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/BrunoCardosoFS/NaxiServer/server"
	"github.com/BrunoCardosoFS/NaxiServer/settings"
)

//go:embed icons/icon_success.ico
var iconSuccessData []byte

// // go:embed icons/icon_failure.ico
// var iconFailureData []byte

func main() {
	// Logs
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)

	fileLog, errLog := os.OpenFile(filepath.Join(exPath, "logs.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	runSystray()

	database.InitDB("file:" + settings.GetSettings().DbPath + "/naxistudio.db" + "?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")

	srv := server.NewServer(settings.GetSettings().DbPath)
	if err := srv.Start(settings.GetSettings().Port); err != nil {
		log.Fatal(err)
	}
}

func onExit() {
	database.DB.Close()

	log.Println("Shutting down server")
}
