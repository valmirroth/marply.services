package service

import (
	"gorm.io/gorm"
	rsdatamodel "vroth.rsdata/RSDataModel"
	"vroth.rsdata/models"
	"vroth.rsdata/repositories"
)

type EntradaNFService struct {
	DB   *gorm.DB
	repo *repositories.EntradaNFRepository
}

// Agora o construtor recebe o DB e o repositório e retorna ponteiro
func NewEntradaNFService(db *gorm.DB, repos *repositories.EntradaNFRepository) *EntradaNFService {
	return &EntradaNFService{
		DB:   db,
		repo: repos,
	}
}

func (s *EntradaNFService) GetListaItensConsultarBaixa() ([]*models.RSDataLancamento, error) {

	return s.repo.GetListaItensConsultarBaixa()
}

func (s *EntradaNFService) GetNextLancRSData(itemEPI string) ([]*models.RSDataLancamento, error) {
	return s.repo.GetNextLancRSData(itemEPI)
}

func (s *EntradaNFService) InserirLogMovimentoRsData(itemepe *models.RSDataLancamento, ret *rsdatamodel.EnvelopeEstRetorno) (*models.RSDataLancamento, error) {
	return s.repo.InserirLogMovimentoRsData(itemepe, ret)
}

func (s *EntradaNFService) LimparDadosBaixaEstoqueSetor(data string) error {
	erro := s.repo.LimparDadosBaixaEstoqueSetor(data)
	return erro
}

func (s *EntradaNFService) GravarMoviementacaoEPISetor(Movi models.EntregaEPISetor) (*models.EntregaEPISetor, error) {
	mov, erro := s.repo.GravarMoviementacaoEPISetor(&Movi)

	return mov, erro
}
