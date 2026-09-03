// Command CCL-Health-Check é o ponto de entrada da aplicação desktop.
// Toda a lógica de inicialização fica em internal/app; este arquivo
// deve permanecer o mais enxuto possível.
package main

import (
	"github.com/BadWolfLuck/CCL-Health-Check/internal/app"
)

func main() {
	app.Run()
}
