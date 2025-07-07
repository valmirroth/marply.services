package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"vroth.rsdata/models"
	"vroth.rsdata/service"
)

var log *service.Logger

func init() {
	// Inicializa logger para settings
	var err error
	log, err = service.NewLogger("ResultSync", "settings", service.ERROR)
	if err != nil {
		panic(fmt.Errorf("não foi possível inicializar o logger: %w", err))
	}
}

// LoadConfigFile carrega o arquivo de configuração de conexão e API REST.
func LoadConfigFile() *models.Parametros {
	file, err := os.Getwd()
	file = file + "\\data\\globalconfig.json"

	filepath, err := getFilePath()
	if err != nil {
		log.Error("Erro ao obter path do config: %v", err)
	}

	var cfg models.Parametros
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error("Falha ao ler arquivo %s: %v", filepath, err)
	}

	if err := json.Unmarshal(content, &cfg); err != nil {
		log.Error("Falha ao parsear JSON do config: %v", err)
	}

	return &cfg
}

// getFilePath retorna o caminho completo do arquivo globalconfig.json.
func getFilePath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Error("Erro ao obter diretório atual: %v", err)
		return "", err
	}

	path := filepath.Join(wd, "data", "globalconfig.json")
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Error("Arquivo de configuração não existe: %s", path)
		return path, err
	}
	if err != nil {
		log.Error("Erro ao verificar arquivo %s: %v", path, err)
		return path, err
	}

	if info.IsDir() {
		err = fmt.Errorf("esperado arquivo, mas %s é diretório", path)
		log.Error("%v", err)
		return path, err
	}

	return path, nil
}
