// Package screens contém as telas completas da aplicação.
package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/domain"
	"github.com/BadWolfLuck/CCL-Health-Check/internal/servicemon"
	uiwidgets "github.com/BadWolfLuck/CCL-Health-Check/internal/ui/widgets"
)

// pollInterval define de quanto em quanto tempo cada serviço
// monitorado tem seu status consultado novamente.
const pollInterval = 5 * time.Second

// monitoredServices é a lista de serviços exibidos nesta tela.
//
// PARA REPLICAR PARA OUTRO SERVIÇO: basta adicionar uma nova linha
// aqui com o DisplayName (nome que aparece na tela) e o
// WindowsServiceName (nome interno do serviço no Windows — o mesmo
// que aparece em services.msc na coluna "Nome do serviço", ou via
// `sc query <nome>` no terminal). Nenhuma outra alteração de código é
// necessária: a linha na tela e o polling são criados automaticamente
// para cada item desta lista.
var monitoredServices = []domain.MonitoredService{
	{
		DisplayName:        "Spooler de Impressão",
		WindowsServiceName: "Spooler",
	},
	{
		DisplayName:        "Windows Update",
		WindowsServiceName: "wuauserv",
	},
}

// serviceRow agrupa os componentes de UI de uma linha da lista, para
// que o loop de polling saiba qual bolinha e qual label atualizar.
type serviceRow struct {
	service domain.MonitoredService
	dot     *uiwidgets.StatusDot
	label   *widget.Label
}

// NewServices monta a tela de monitoramento de serviços. Cada serviço
// em monitoredServices vira uma linha com bolinha de status + nome +
// descrição textual do status atual. Um polling em background atualiza
// todas as linhas a cada pollInterval.
func NewServices(w fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Monitoramento de serviços",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	rows := make([]*serviceRow, 0, len(monitoredServices))
	rowsContainer := container.NewVBox()

	for _, svc := range monitoredServices {
		row := newServiceRow(svc)
		rows = append(rows, row)
		rowsContainer.Add(row.render())
	}

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		rowsContainer,
	)

	// startPolling roda em uma goroutine separada; cada atualização de
	// widget é despachada de volta para a thread de UI do Fyne via
	// fyne.Do, que é a forma segura de atualizar widgets a partir de
	// outra goroutine.
	startPolling(w, rows)

	return content
}

// newServiceRow cria os widgets de uma linha (bolinha + nome + status)
// para o serviço informado.
func newServiceRow(svc domain.MonitoredService) *serviceRow {
	return &serviceRow{
		service: svc,
		dot:     uiwidgets.NewStatusDot(),
		label:   widget.NewLabel(domain.StatusPending.String()),
	}
}

// render monta o layout horizontal de uma linha: bolinha, nome do
// serviço e descrição textual do status.
func (r *serviceRow) render() fyne.CanvasObject {
	name := widget.NewLabel(r.service.DisplayName)
	name.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewHBox(
		r.dot,
		name,
		widget.NewLabel("—"),
		r.label,
	)
}

// startPolling dispara uma goroutine que consulta todos os serviços
// imediatamente e depois a cada pollInterval, até a janela ser fechada.
func startPolling(w fyne.Window, rows []*serviceRow) {
	ticker := time.NewTicker(pollInterval)

	pollOnce := func() {
		for _, row := range rows {
			row := row // evita captura por referência da variável de loop
			status, err := servicemon.QueryStatus(row.service.WindowsServiceName)
			if err != nil {
				fyne.LogError(fmt.Sprintf("consultar serviço %q", row.service.WindowsServiceName), err)
				status = domain.StatusStopped
			}

			// fyne.Do garante que a atualização do widget aconteça na
			// goroutine de UI, que é obrigatório no Fyne 2.6+.
			fyne.Do(func() {
				row.dot.SetStatus(status)
				row.label.SetText(status.String())
			})
		}
	}

	go func() {
		pollOnce() // primeira consulta imediata, sem esperar o primeiro tick
		for range ticker.C {
			pollOnce()
		}
	}()

	// Encerra o ticker quando a janela for fechada, evitando goroutine
	// órfã consumindo recursos após o app ser fechado.
	w.SetOnClosed(func() {
		ticker.Stop()
	})
}
