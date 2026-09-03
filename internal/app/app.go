// Package app inicializa a aplicação Fyne: cria a instância do app,
// a janela principal e monta a tela inicial.
package app

import (
	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/config"
	"github.com/BadWolfLuck/CCL-Health-Check/internal/ui/screens"
)

const (
	appID     = "com.badwolfluck.cclhealthcheck"
	appTitle  = "CCL-Health-Check"
	winWidth  = 900
	winHeight = 600
)

// Run cria e inicia a aplicação. É a única função exportada deste
// pacote e é chamada diretamente por cmd/CCL-Health-Check/main.go.
func Run() {
	a := fyneapp.NewWithID(appID)

	_, err := config.Load()
	if err != nil {
		fyne.LogError("falha ao carregar configuração, usando padrão", err)
	}

	w := a.NewWindow(appTitle)
	w.Resize(fyne.NewSize(winWidth, winHeight))
	w.SetMaster()

	w.SetContent(screens.NewServices(w))

	w.ShowAndRun()
}
