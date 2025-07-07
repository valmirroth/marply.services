package repositories

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	rsdatamodel "vroth.rsdata/RSDataModel"
	"vroth.rsdata/models"
	"vroth.rsdata/models/entities"
)

type EntradaNFRepository struct {
	DB *gorm.DB
}

func NewEntradaNFRepository(db *gorm.DB) *EntradaNFRepository {
	return &EntradaNFRepository{
		DB: db,
	}
}

func (repo *EntradaNFRepository) GetLastNfKorp(produto string) error {
	return nil
}

func (repo *EntradaNFRepository) GetLastNfRSData(produto string) error {
	return nil
}

func (repo *EntradaNFRepository) GetNextLancRSData(itemEPI string) ([]*models.RSDataLancamento, error) {
	var listNf []*models.RSDataLancamento

	trans := repo.DB.Raw(GetListaPendencia).Scan(&listNf)

	if trans.RowsAffected == 0 {
		return nil, nil
	}

	if trans.Error != nil {
		return nil, trans.Error
	}
	return listNf, nil
}

func (repo *EntradaNFRepository) IntOrNull(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func StringToNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func (repo *EntradaNFRepository) InserirLogMovimentoRsData(itemepi *models.RSDataLancamento, ret *rsdatamodel.EnvelopeEstRetorno) (*models.RSDataLancamento, error) {
	var log models.RSDataLLog

	log.QtdMovimento = itemepi.QtdMovimento
	log.Descricao = StringToNullString(itemepi.Descricao)
	log.Item = StringToNullString(itemepi.Item)
	log.RecnoHF = itemepi.RecnoHF
	log.DataMovimento = StringToNullString(itemepi.DataMovimento)
	log.RecnoHistlise = itemepi.RecnoHistlise
	log.GradeTamanho = StringToNullString(itemepi.TamanhoGrade)
	var detalhes string
	if ret != nil {
		log.RetornoRSData = StringToNullString(ret.Body.InsertEstoqueResp.RetornoMsg)
		detalhes = ret.Body.InsertEstoqueResp.RsdataReturn.Mensagens.Mensagem.TxDescricao
	}
	log.Cnpj = StringToNullString(itemepi.Cnpj)
	log.LocalEstoque = StringToNullString(itemepi.LocalEstoque)
	fmt.Println(detalhes)
	log.RetornoDetalhado = StringToNullString(detalhes)
	trans := repo.DB.Select("LOCAL_ESTOQUE", "CNPJ", "RECNO_HISTLISE", "ITEM", "DESCRICAO", "QTDEMOV", "GRADE_TAMANHO", "DATA_MOV", "RECNO_HISTLISE_FOR", "RETORNORSDATA", "RETORNO_DETALHADO").Create(&log).Scan(&log)

	if trans.Error != nil {
		return nil, trans.Error
	}
	return itemepi, nil
}

func (repo *EntradaNFRepository) GetListaItensConsultarBaixa() ([]*models.RSDataLancamento, error) {
	var listNf []*models.RSDataLancamento

	trans := repo.DB.Raw(GetListaItens).Scan(&listNf)

	if trans.RowsAffected == 0 {
		return nil, nil
	}

	if trans.Error != nil {
		return nil, trans.Error
	}
	return listNf, nil
}

func (repo *EntradaNFRepository) LimparDadosBaixaEstoqueSetor(data string) error {
	rs := repo.DB.Exec("delete from CST_RSDATA_ENTREGA_SETOR where DataMovimento = ?", data)
	if rs.Error != nil {
		return rs.Error
	}
	return nil
}

func IntOrNull(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func (repo *EntradaNFRepository) getCodigoItemKorpGrade(Movi *models.EntregaEPISetor) (*models.RSDataLancamento, error) {
	var listNf []*models.RSDataLancamento
	if Movi.Tamanho == "Padrão" {
		return nil, nil
	}
	trans := repo.DB.Raw(GetItemGrade, &Movi.CodIntegracaoEpi, &Movi.Tamanho).Scan(&listNf)

	if trans.RowsAffected == 0 {
		return nil, nil
	}

	if trans.Error != nil {
		return nil, trans.Error
	}
	if len(listNf) > 0 {
		return listNf[0], nil
	}
	return nil, nil
}

func (repo *EntradaNFRepository) GravarMoviementacaoEPISetor(Movi *models.EntregaEPISetor) (*models.EntregaEPISetor, error) {
	var mov entities.EntregaEPISetor

	if Movi.NomeEpi == "BOTA FLORESTAL C/AÇO PRETA nº 44 CA 47412" {
		fmt.Println("BOTA FLORESTAL C/AÇO PRETA nº 44 CA 47412")
	}
	item, err := repo.getCodigoItemKorpGrade(Movi)
	if err != nil {
		return nil, err
	}

	if item != nil {
		mov.CodIntegracaoEpi = StringToNullString(item.Item)
		Movi.CodIntegracaoEpi = item.Item
	} else {
		mov.CodIntegracaoEpi = StringToNullString(Movi.CodIntegracaoEpi)

	}
	mov.Tamanho = StringToNullString(Movi.Tamanho)
	mov.CodIntegracaoEpi = StringToNullString(Movi.CodIntegracaoEpi)
	mov.DataMovimento = StringToNullString(Movi.DataMovimento)
	mov.EmpresaRecno = IntOrNull(int64(Movi.EmpresaRecno))
	mov.Quantidade = StringToNullString(Movi.Quantidade)
	mov.Local = StringToNullString(Movi.Local)
	mov.IdSetor = StringToNullString(Movi.IdSetor)
	mov.NomeSetor = StringToNullString(Movi.NomeSetor)
	mov.TipoEpi = StringToNullString(Movi.TipoEpi)
	mov.Cnpj = StringToNullString(Movi.Cnpj)
	mov.NomeEpi = StringToNullString(Movi.NomeEpi)
	trans := repo.DB.Create(&mov).Scan(&mov)

	if trans.Error != nil {
		return nil, trans.Error
	}
	return Movi, nil

}
