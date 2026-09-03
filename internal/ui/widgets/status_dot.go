// Package widgets reúne componentes de UI reutilizáveis entre telas.
package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/domain"
)

// Cores fixas para cada status. Centralizadas aqui para que qualquer
// ajuste de paleta (ex: seguir um tema dark/light customizado) seja
// feito em um único lugar.
var (
	colorPending    = color.NRGBA{R: 158, G: 158, B: 158, A: 255} // cinza
	colorRunning    = color.NRGBA{R: 46, G: 160, B: 67, A: 255}   // verde
	colorStopped    = color.NRGBA{R: 218, G: 54, B: 51, A: 255}   // vermelho
	colorRestarting = color.NRGBA{R: 240, G: 173, B: 27, A: 255}  // amarelo
)

// diameter é o tamanho fixo da bolinha em pixels.
const diameter = 16

// StatusDot é um indicador visual circular que representa o
// domain.ServiceStatus de um serviço monitorado. É um componente
// "burro": não consulta nada sozinho, apenas exibe o status que a
// tela (screens) informar via SetStatus.
type StatusDot struct {
	widget.BaseWidget

	circle *canvas.Circle
	status domain.ServiceStatus
}

// NewStatusDot cria uma nova bolinha de status, iniciando como
// domain.StatusPending (cinza) até a primeira consulta ser feita.
func NewStatusDot() *StatusDot {
	d := &StatusDot{
		circle: canvas.NewCircle(colorPending),
		status: domain.StatusPending,
	}
	d.ExtendBaseWidget(d)
	return d
}

// SetStatus atualiza a cor da bolinha conforme o novo status e força
// o redesenho do componente. Chame este método sempre que uma nova
// consulta ao serviço (servicemon.QueryStatus) retornar.
func (d *StatusDot) SetStatus(status domain.ServiceStatus) {
	d.status = status
	d.circle.FillColor = colorFor(status)
	d.circle.Refresh()
}

// colorFor mapeia cada domain.ServiceStatus para sua cor correspondente.
func colorFor(status domain.ServiceStatus) color.Color {
	switch status {
	case domain.StatusRunning:
		return colorRunning
	case domain.StatusStopped:
		return colorStopped
	case domain.StatusRestarting:
		return colorRestarting
	default:
		return colorPending
	}
}

// CreateRenderer implementa fyne.Widget. Fyne chama isto internamente;
// não é necessário (nem recomendado) chamar diretamente.
func (d *StatusDot) CreateRenderer() fyne.WidgetRenderer {
	return &statusDotRenderer{dot: d}
}

// statusDotRenderer controla o layout e a renderização de StatusDot.
type statusDotRenderer struct {
	dot *StatusDot
}

func (r *statusDotRenderer) Layout(size fyne.Size) {
	r.dot.circle.Resize(fyne.NewSize(diameter, diameter))
}

func (r *statusDotRenderer) MinSize() fyne.Size {
	return fyne.NewSize(diameter, diameter)
}

func (r *statusDotRenderer) Refresh() {
	r.dot.circle.Refresh()
}

func (r *statusDotRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.dot.circle}
}

func (r *statusDotRenderer) Destroy() {}
