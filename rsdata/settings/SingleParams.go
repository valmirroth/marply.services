package settings

import (
	"sync"

	"vroth.rsdata/models"
)

var once sync.Once

var parametros *models.Parametros

// função singleton
func GetInstange() *models.Parametros {
	if parametros == nil {
		once.Do(
			func() {
				parametros = LoadConfigFile()
			})
	} else {
		//fmt.Println("Parâmtros em memória...")
	}
	return parametros

}
