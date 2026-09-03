// Package domain contém as regras de negócio da aplicação, isoladas
// de qualquer detalhe de interface gráfica (Fyne) ou de sistema
// operacional (Windows Service Manager).
package domain

// ServiceStatus representa o estado observável de um serviço monitorado.
// Este tipo é o contrato entre a camada que consulta o Windows
// (internal/servicemon) e a camada de UI (internal/ui) — a UI nunca
// deve conhecer detalhes da API do Windows, só este enum.
type ServiceStatus int

const (
	// StatusPending indica que o status ainda não foi consultado,
	// ou a consulta está em andamento. Cor: cinza.
	StatusPending ServiceStatus = iota

	// StatusRunning indica que o serviço está em execução normal.
	// Cor: verde.
	StatusRunning

	// StatusStopped indica que o serviço não está rodando, foi
	// parado, ou não existe/não foi encontrado no Service Manager.
	// Cor: vermelho.
	StatusStopped

	// StatusRestarting indica que o serviço está em uma transição de
	// estado (iniciando, parando, ou reiniciando). Cor: amarelo.
	StatusRestarting
)

// String retorna uma descrição legível do status, usada em labels de UI.
func (s ServiceStatus) String() string {
	switch s {
	case StatusRunning:
		return "Em execução"
	case StatusStopped:
		return "Parado / não encontrado"
	case StatusRestarting:
		return "Reiniciando"
	default:
		return "Pendente"
	}
}

// MonitoredService descreve um serviço do Windows a ser monitorado.
// Para adicionar um novo serviço à tela de monitoramento, basta criar
// um novo valor deste tipo — veja internal/ui/screens/services.go.
type MonitoredService struct {
	// DisplayName é o nome amigável mostrado na interface.
	DisplayName string

	// WindowsServiceName é o nome interno do serviço no Service
	// Control Manager do Windows (o mesmo usado em "services.msc",
	// coluna "Nome do serviço", ou via `sc query <nome>`).
	WindowsServiceName string
}
