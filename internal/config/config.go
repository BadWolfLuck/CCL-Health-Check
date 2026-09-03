// Package config cuida da leitura e escrita das preferências do usuário
// em um arquivo JSON local, já que o projeto não utiliza banco de dados.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// fileName é o nome do arquivo de configuração salvo no diretório
// de configuração do usuário (ex: %AppData% no Windows).
const fileName = "config.json"

// Config representa as preferências persistentes do usuário.
// Adicione novos campos aqui conforme o app crescer.
type Config struct {
	Theme        string `json:"theme"`         // "light", "dark" ou "system"
	WindowWidth  int    `json:"window_width"`
	WindowHeight int    `json:"window_height"`
}

// Default retorna a configuração padrão usada quando ainda não existe
// um arquivo salvo ou quando a leitura falha.
func Default() *Config {
	return &Config{
		Theme:        "system",
		WindowWidth:  900,
		WindowHeight: 600,
	}
}

// path retorna o caminho completo do arquivo de configuração,
// dentro do diretório padrão de configuração do sistema operacional.
func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(dir, "meu-projeto")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, fileName), nil
}

// Load lê a configuração do disco. Se o arquivo ainda não existir,
// retorna a configuração padrão sem erro.
func Load() (*Config, error) {
	p, err := path()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save grava a configuração atual no disco em formato JSON legível.
func (c *Config) Save() error {
	p, err := path()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0o644)
}
