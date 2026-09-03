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

	cfg, err := config.Load()
	if err != nil {
		// Ainda não há uma tela de erro dedicada; por enquanto seguimos
		// com a configuração padrão e apenas registramos no log.
		fyne.LogError("falha ao carregar configuração, usando padrão", err)
		cfg = config.Default()
	}

	w := a.NewWindow(appTitle)
	w.Resize(fyne.NewSize(winWidth, winHeight))
	w.SetMaster()

	w.SetContent(screens.NewHome(w, cfg))

	w.ShowAndRun()
}
