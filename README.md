# CCL-Health-Check

Repositório oficial: https://github.com/BadWolfLuck/CCL-Health-Check

Aplicação desktop para Windows 11, desenvolvida em Go utilizando [Fyne](https://fyne.io/) para a interface gráfica.

## Requisitos

- Go 1.22 ou superior
- Um compilador C (necessário para o Fyne/CGO). No Windows, recomenda-se o [MSYS2/MinGW-w64](https://www.msys2.org/) ou o TDM-GCC.

## Rodando localmente

```bash
go mod tidy
go run ./cmd/CCL-Health-Check
```

## Estrutura do projeto

```
cmd/CCL-Health-Check/     -> ponto de entrada (main.go)
internal/app/        -> inicialização da aplicação (janela, tema, config)
internal/ui/         -> telas (screens), componentes (widgets) e tema visual
internal/config/     -> leitura/escrita de configuração local (JSON)
internal/domain/     -> regras de negócio, independentes da UI
pkg/                 -> código reutilizável por outros projetos (se houver)
assets/              -> ícones e fontes
build/windows/       -> ícone .ico e artefatos específicos do build Windows
scripts/             -> scripts de build/empacotamento
```

## Gerando o executável (.exe)

```powershell
go install fyne.io/tools/cmd/fyne@latest
.\scripts\build.ps1
```

## Licença

Defina a licença do projeto no arquivo `LICENSE`.
