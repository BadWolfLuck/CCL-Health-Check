# Script de build do app para Windows.
# Requer: Go instalado, e o utilitário "fyne" (go install fyne.io/tools/cmd/fyne@latest)
#
# Uso:
#   .\scripts\build.ps1

$ErrorActionPreference = "Stop"

Write-Host "Baixando dependências..."
go mod tidy

Write-Host "Empacotando aplicação para Windows..."
Push-Location cmd/meu-projeto
fyne package -os windows -icon ../../build/windows/icon.ico -name "Meu Projeto"
Pop-Location

Write-Host "Build concluído. Executável gerado em cmd/meu-projeto/"
