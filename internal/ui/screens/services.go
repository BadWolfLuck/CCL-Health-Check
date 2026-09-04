// Package screens contém as telas completas da aplicação.
package screens

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
// necessária: a linha na tela, os botões e o polling são criados
// automaticamente para cada item desta lista.
var monitoredServices = []domain.MonitoredService{
	{
		DisplayName:        "FortiNAC Agent",
		WindowsServiceName: "BNPagent",
	},
	{
		DisplayName:        "Qualys Agent",
		WindowsServiceName: "QualysAgent",
	},
	{
		DisplayName:        "Rapid7 Agent",
		WindowsServiceName: "ir_agent",
	},
	{
		DisplayName:        "Spooler de Impressão",
		WindowsServiceName: "Spooler",
	},
	{
		DisplayName:        "Windows Update",
		WindowsServiceName: "wuauserv",
	},
	{
		DisplayName:        "Trend Micro Unauthorized Change Prevention Service",
		WindowsServiceName: "TMBMServer",
	},
	{
		DisplayName:        "TrendAI Application Control Service",
		WindowsServiceName: "TMiACAgentSvc",
	},
	{
		DisplayName:        "TrendAI Endpoint Basecamp",
		WindowsServiceName: "Trend Micro Endpoint Basecamp",
	},
	{
		DisplayName:        "TrendAI Network Service",
		WindowsServiceName: "tm_netsrv",
	},
	{
		DisplayName:        "TrendAI™ Vulnerability Protection Service",
		WindowsServiceName: "iVPAgent",
	},
}

// serviceRow agrupa os componentes de UI de uma linha da lista, para
// que o loop de polling e os botões de controle saibam qual bolinha,
// label e conjunto de botões atualizar.
type serviceRow struct {
	service  domain.MonitoredService
	dot      *uiwidgets.StatusDot
	label    *widget.Label
	controls *uiwidgets.ServiceControls
	window   fyne.Window
}

// NewServices monta a tela de monitoramento de serviços. Cada serviço
// em monitoredServices vira uma linha com bolinha de status + nome +
// descrição textual do status + botões de controle. Um polling em
// background atualiza todas as linhas a cada pollInterval.
func NewServices(w fyne.Window) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Monitoramento de serviços",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	rows := make([]*serviceRow, 0, len(monitoredServices))
	rowsContainer := container.NewVBox()

	for _, svc := range monitoredServices {
		row := newServiceRow(w, svc)
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

// newServiceRow cria os widgets de uma linha (bolinha + nome + status
// + botões) para o serviço informado. window é necessária para exibir
// diálogos de erro em caso de falha ao controlar o serviço.
func newServiceRow(w fyne.Window, svc domain.MonitoredService) *serviceRow {
	row := &serviceRow{
		service: svc,
		dot:     uiwidgets.NewStatusDot(),
		label:   widget.NewLabel(domain.StatusPending.String()),
		window:  w,
	}

	row.controls = uiwidgets.NewServiceControls(
		func() { row.runAction("iniciar", servicemon.Start) },
		func() { row.runAction("parar", servicemon.Stop) },
		func() { row.runAction("reiniciar", servicemon.Restart) },
	)

	return row
}

// render monta o layout horizontal de uma linha: bolinha, nome do
// serviço, descrição textual do status e os botões de controle.
func (r *serviceRow) render() fyne.CanvasObject {
	name := widget.NewLabel(r.service.DisplayName)
	name.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewHBox(
		r.dot,
		name,
		widget.NewLabel("—"),
		r.label,
		layoutSpacer(),
		r.controls,
	)
}

// layoutSpacer cria um espaço flexível entre o status e os botões,
// empurrando os botões para a direita da linha.
func layoutSpacer() fyne.CanvasObject {
	return widget.NewLabel(" ")
}

// runAction executa uma ação de controle (Start/Stop/Restart) em uma
// goroutine separada, para não travar a interface durante a chamada
// bloqueante ao Service Control Manager. Em caso de erro, um diálogo
// é exibido ao usuário na thread de UI. Ao final, o status da linha é
// atualizado imediatamente, sem esperar o próximo ciclo de polling.
func (r *serviceRow) runAction(actionLabel string, action func(string) error) {
	// Desabilita os botões imediatamente para evitar cliques
	// duplicados enquanto a ação está em andamento.
	fyne.Do(func() {
		r.controls.SetStatus(domain.StatusRestarting)
		r.label.SetText("Executando: " + actionLabel + "...")
	})

	go func() {
		err := action(r.service.WindowsServiceName)
		if err != nil {
			fyne.LogError(fmt.Sprintf("%s serviço %q", actionLabel, r.service.WindowsServiceName), err)
			fyne.Do(func() {
				dialog.ShowError(
					fmt.Errorf("não foi possível %s o serviço %q: %w", actionLabel, r.service.DisplayName, err),
					r.window,
				)
			})
		}

		// Independentemente de sucesso ou falha, consultamos o status
		// real para refletir o estado atual do serviço na tela.
		status, statusErr := servicemon.QueryStatus(r.service.WindowsServiceName)
		if statusErr != nil {
			fyne.LogError(fmt.Sprintf("consultar serviço %q após %s", r.service.WindowsServiceName, actionLabel), statusErr)
			status = domain.StatusStopped
		}

		fyne.Do(func() {
			r.dot.SetStatus(status)
			r.label.SetText(status.String())
			r.controls.SetStatus(status)
		})
	}()
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
				row.controls.SetStatus(status)
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
