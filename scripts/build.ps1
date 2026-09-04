# Script de build do app para Windows.
#
# Requer:
#   - Go instalado
#   - fyne CLI: go install fyne.io/tools/cmd/fyne@latest
#
# O "fyne package" gera automaticamente o manifest do executável. Se
# encontrar um arquivo "<nome-do-exe>.exe.manifest" na mesma pasta do
# .exe final, ele usa esse manifest customizado em vez do padrão. Por
# isso copiamos build/windows/app.manifest para o nome esperado antes
# de empacotar.
#
# O manifest exige execução como Administrador
# (requireAdministrator), necessário para o app poder
# Iniciar/Parar/Reiniciar serviços do Windows.
#
# Uso:
#   .\scripts\build.ps1

$ErrorActionPreference = "Stop"

Write-Host "Baixando dependências..."
go mod tidy

$exeName = "CCL-Health-Check"
$exeDir = "cmd/CCL-Health-Check"
$manifestTarget = Join-Path $exeDir "$exeName.exe.manifest"

Write-Host "Copiando manifest customizado (requireAdministrator)..."
Copy-Item -Path "build/windows/app.manifest" -Destination $manifestTarget -Force

Write-Host "Empacotando aplicação para Windows..."
Push-Location $exeDir
fyne package -os windows -icon ../../build/windows/icon.ico -name $exeName
Pop-Location

Write-Host "Build concluído. Executável gerado em $exeDir/"
Write-Host "O Windows vai pedir elevação (UAC) sempre que o .exe for iniciado."
