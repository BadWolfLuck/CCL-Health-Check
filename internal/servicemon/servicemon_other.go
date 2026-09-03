//go:build !windows

package servicemon

import (
	"errors"

	"github.com/BadWolfLuck/CCL-Health-Check/internal/domain"
)

// QueryStatus não é suportado fora do Windows. Este stub existe apenas
// para permitir compilar o pacote em outras plataformas (ex: para
// rodar linters ou editores em Linux/Mac); a aplicação em si é para
// Windows 11, conforme o escopo do projeto.
func QueryStatus(windowsServiceName string) (domain.ServiceStatus, error) {
	return domain.StatusStopped, errors.New("consulta de serviços do Windows só é suportada em Windows")
}
