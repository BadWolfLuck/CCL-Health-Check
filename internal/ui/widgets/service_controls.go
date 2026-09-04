// Package widgets reúne componentes de UI reutilizáveis entre telas.
package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/domain"
)

// ServiceControls agrupa os três botões de controle de um serviço:
// Iniciar, Parar e Reiniciar. O widget não executa nenhuma ação por
// conta própria — ele apenas dispara os callbacks informados na
// criação e habilita/desabilita os botões conforme o status atual,
// via SetStatus.
type ServiceControls struct {
	widget.BaseWidget

	startBtn   *widget.Button
	stopBtn    *widget.Button
	restartBtn *widget.Button

	container *fyne.Container
}

// NewServiceControls cria os três botões de controle. Os callbacks
// onStart, onStop e onRestart são chamados quando o respectivo botão
// é clicado; a tela que criar este widget é responsável por executar
// a ação real (via internal/servicemon) dentro desses callbacks,
// tipicamente em uma goroutine para não travar a UI.
func NewServiceControls(onStart, onStop, onRestart func()) *ServiceControls {
	c := &ServiceControls{
		startBtn:   widget.NewButton("Iniciar", onStart),
		stopBtn:    widget.NewButton("Parar", onStop),
		restartBtn: widget.NewButton("Reiniciar", onRestart),
	}

	c.container = container.NewHBox(c.startBtn, c.stopBtn, c.restartBtn)
	c.ExtendBaseWidget(c)

	// Estado inicial: enquanto o status ainda não foi consultado
	// (StatusPending), nenhuma ação faz sentido.
	c.SetStatus(domain.StatusPending)

	return c
}

// SetStatus habilita/desabilita cada botão de acordo com o status
// atual do serviço, evitando ações sem sentido — por exemplo, não faz
// sentido oferecer "Iniciar" para um serviço que já está rodando, nem
// "Parar" ou "Reiniciar" para um serviço que não existe/está parado.
//
// Chame este método sempre que uma nova consulta de status retornar,
// no mesmo lugar em que StatusDot.SetStatus é chamado.
func (c *ServiceControls) SetStatus(status domain.ServiceStatus) {
	switch status {
	case domain.StatusRunning:
		c.startBtn.Disable()
		c.stopBtn.Enable()
		c.restartBtn.Enable()

	case domain.StatusStopped:
		c.startBtn.Enable()
		c.stopBtn.Disable()
		c.restartBtn.Disable()

	case domain.StatusRestarting:
		// Durante uma transição, evitamos empilhar novos comandos até
		// o próximo status conhecido chegar.
		c.startBtn.Disable()
		c.stopBtn.Disable()
		c.restartBtn.Disable()

	default: // domain.StatusPending
		c.startBtn.Disable()
		c.stopBtn.Disable()
		c.restartBtn.Disable()
	}
}

// CreateRenderer implementa fyne.Widget.
func (c *ServiceControls) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}
