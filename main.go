package main

import (
	_ "embed"

	"fyne.io/systray"
)

//go:embed icons/icon.ico
var iconData []byte

func main() {
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

	go func() { //goroutine: tipo uma função assincrona
		for range mSettings.ClickedCh {
			println("Configurações")
		}
	}()
}

func onExit() {
	println("Encerrando servidor")
}
