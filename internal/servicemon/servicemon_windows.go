//go:build windows

// Package servicemon consulta o Windows Service Control Manager (SCM)
// para obter o status atual de serviços do sistema. É a única camada
// que conhece a API do Windows; tudo o que sai daqui já vem traduzido
// para domain.ServiceStatus.
package servicemon

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/domain"
)

// QueryStatus consulta o Service Control Manager do Windows e retorna
// o status atual do serviço identificado por windowsServiceName (o
// nome interno, não o nome de exibição).
//
// Se o serviço não existir, ou se o handle não puder ser aberto por
// qualquer motivo, o retorno é domain.StatusStopped — do ponto de
// vista do monitoramento, "não existe" e "não está rodando" recebem
// o mesmo tratamento (bolinha vermelha).
func QueryStatus(windowsServiceName string) (domain.ServiceStatus, error) {
	m, err := mgr.Connect()
	if err != nil {
		return domain.StatusStopped, fmt.Errorf("conectar ao service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		// Serviço não encontrado ou sem permissão de acesso: tratamos
		// como parado, não como erro fatal da aplicação.
		return domain.StatusStopped, nil
	}
	defer s.Close()

	status, err := s.Query()
	if err != nil {
		return domain.StatusStopped, fmt.Errorf("consultar status do serviço %q: %w", windowsServiceName, err)
	}

	return mapState(status.State), nil
}

// mapState traduz o estado nativo do Windows (svc.State) para o
// domain.ServiceStatus usado pela UI. Esta é a única função que
// precisa mudar se, no futuro, quisermos exibir mais granularidade
// (ex: distinguir "parando" de "reiniciando").
func mapState(state svc.State) domain.ServiceStatus {
	switch state {
	case svc.Running:
		return domain.StatusRunning

	case svc.StartPending, svc.StopPending, svc.ContinuePending, svc.PausePending:
		// Qualquer estado de transição é tratado como "reiniciando"
		// para fins de exibição (bolinha amarela).
		return domain.StatusRestarting

	case svc.Stopped, svc.Paused:
		return domain.StatusStopped

	default:
		return domain.StatusStopped
	}
}
