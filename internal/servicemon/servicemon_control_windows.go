//go:build windows

package servicemon

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// controlTimeout é o tempo máximo de espera por uma transição de
// estado (ex: aguardar o serviço realmente parar antes de considerar
// a operação como falha). Serviços mais pesados podem demorar alguns
// segundos para finalizar suas rotinas de encerramento.
const controlTimeout = 15 * time.Second

// pollInterval é o intervalo entre verificações de estado enquanto
// aguardamos uma transição (ex: esperar Stopped antes de dar Start
// de novo, no caso do Restart).
const controlPollInterval = 300 * time.Millisecond

// Start inicia o serviço identificado por windowsServiceName.
//
// Requer que o processo esteja rodando com privilégio de
// administrador (ver build/windows/app.manifest); caso contrário, o
// Windows retorna erro de acesso negado.
func Start(windowsServiceName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("conectar ao service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return fmt.Errorf("abrir serviço %q: %w", windowsServiceName, err)
	}
	defer s.Close()

	if err := s.Start(); err != nil {
		return fmt.Errorf("iniciar serviço %q: %w", windowsServiceName, err)
	}

	return nil
}

// Stop envia o comando de parada ao serviço e aguarda até
// controlTimeout pela confirmação de que ele efetivamente parou.
func Stop(windowsServiceName string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("conectar ao service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return fmt.Errorf("abrir serviço %q: %w", windowsServiceName, err)
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("enviar comando de parada ao serviço %q: %w", windowsServiceName, err)
	}

	return waitForState(s, status, svc.Stopped, controlTimeout)
}

// Restart para o serviço e, assim que ele confirmar o estado
// Stopped, inicia novamente. Não existe um comando nativo de
// "restart" no Service Control Manager — por isso encadeamos Stop
// seguido de Start, respeitando o tempo real que o serviço levar
// para parar.
func Restart(windowsServiceName string) error {
	if err := Stop(windowsServiceName); err != nil {
		return fmt.Errorf("parar serviço antes de reiniciar: %w", err)
	}

	if err := Start(windowsServiceName); err != nil {
		return fmt.Errorf("iniciar serviço após parada: %w", err)
	}

	return nil
}

// waitForState aguarda até que o serviço atinja o estado desejado,
// consultando periodicamente seu status. Retorna erro se o timeout
// for atingido antes da transição se completar.
func waitForState(s *mgr.Service, current svc.Status, want svc.State, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for current.State != want {
		if time.Now().After(deadline) {
			return fmt.Errorf("tempo esgotado aguardando o serviço atingir o estado %v (estado atual: %v)", want, current.State)
		}

		time.Sleep(controlPollInterval)

		next, err := s.Query()
		if err != nil {
			return fmt.Errorf("consultar status durante espera de transição: %w", err)
		}
		current = next
	}

	return nil
}
