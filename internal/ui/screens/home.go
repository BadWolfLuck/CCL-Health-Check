// Package screens contém as telas completas da aplicação. Cada arquivo
// deste pacote monta um fyne.CanvasObject pronto para ser usado como
// conteúdo de uma janela.
package screens

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/config"
)

// NewHome monta a tela inicial da aplicação. Recebe a janela (para
// eventuais diálogos ou navegação futura) e a configuração carregada.
func NewHome(w fyne.Window, cfg *config.Config) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"CCL-Health-Check",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	status := widget.NewLabel(fmt.Sprintf("Tema atual: %s", cfg.Theme))

	counter := 0
	counterLabel := widget.NewLabel("Cliques: 0")
	button := widget.NewButton("Clique aqui", func() {
		counter++
		counterLabel.SetText(fmt.Sprintf("Cliques: %d", counter))
	})

	return container.NewVBox(
		title,
		widget.NewSeparator(),
		status,
		counterLabel,
		button,
	)
}
