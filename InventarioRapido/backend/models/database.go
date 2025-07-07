package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

type Database struct {
	DB *sql.DB
}

// NewDatabase cria uma nova instância da conexão com o banco
func NewDatabase() (*Database, error) {
	server := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	// String de conexão para SQL Server
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=true;TrustServerCertificate=true",
		server, user, password, port, database)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %v", err)
	}

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao ping no banco: %v", err)
	}

	log.Println("Conexão com SQL Server estabelecida com sucesso!")

	return &Database{DB: db}, nil
}

// Close fecha a conexão com o banco
func (d *Database) Close() error {
	return d.DB.Close()
}

// InserirContagem insere uma nova contagem no banco
func (d *Database) InserirContagem(contagem ContagemRequest) (*ContagemLocal, error) {
	// SQL Server: use @p1..@pN e OUTPUT para capturar os valores gerados
	const query = `
		INSERT INTO CST_CONTAGEM_LOCAL 
			(CODIGO_ITEM, DESCRICAO_ITEM, LOCAL, QUANTIDADE, VOLUMES, USUARIO_CONTAGEM)
		OUTPUT INSERTED.ID, INSERTED.DATA_CONTAGEM, INSERTED.STATUS
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6);
	`

	// Converter quantidade (string -> float64)
	quantidade, err := strconv.ParseFloat(contagem.Quantidade, 64)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter quantidade: %w", err)
	}

	// Usuário padrão se não informado
	usuario := contagem.UsuarioContagem
	if strings.TrimSpace(usuario) == "" {
		usuario = "Sistema"
	}

	// Contexto (opcional, mas recomendado)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Variáveis que receberão o OUTPUT
	var (
		id           string // ajuste para o tipo real da coluna (ex.: int64 se for INT)
		dataContagem time.Time
		status       string
	)

	// Executa o INSERT e lê o OUTPUT
	if err := d.DB.
		QueryRowContext(ctx, query,
			contagem.CodigoItem,
			contagem.DescricaoItem,
			contagem.Local,
			quantidade,
			contagem.Volumes,
			usuario,
		).
		Scan(&id, &dataContagem, &status); err != nil {
		return nil, fmt.Errorf("erro ao inserir contagem: %w", err)
	}

	// Monta o retorno
	resultado := &ContagemLocal{
		ID:              id,
		CodigoItem:      contagem.CodigoItem,
		DescricaoItem:   contagem.DescricaoItem,
		Local:           contagem.Local,
		Quantidade:      quantidade,
		Volumes:         contagem.Volumes,
		DataContagem:    dataContagem,
		UsuarioContagem: usuario,
		Status:          status,
	}

	return resultado, nil
}

// ListarContagens retorna todas as contagens
func (d *Database) ListarContagens() ([]ContagemLocal, error) {
	query := `
                SELECT cast(ID as varchar(150)) as ID, CODIGO_ITEM, DESCRICAO_ITEM, LOCAL, QUANTIDADE, 
                           VOLUMES, DATA_CONTAGEM, USUARIO_CONTAGEM, STATUS
                FROM CST_CONTAGEM_LOCAL
                WHERE STATUS = 'ATIVO'
                ORDER BY DATA_CONTAGEM DESC
        `

	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar contagens: %v", err)
	}
	defer rows.Close()

	var contagens []ContagemLocal

	for rows.Next() {
		var c ContagemLocal
		var rawID interface{}
		err := rows.Scan(
			&rawID,
			&c.CodigoItem,
			&c.DescricaoItem,
			&c.Local,
			&c.Quantidade,
			&c.Volumes,
			&c.DataContagem,
			&c.UsuarioContagem,
			&c.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear registro: %v", err)
		}

		// Converter GUID para string
		c.ID = fmt.Sprintf("%s", rawID)
		contagens = append(contagens, c)
	}

	return contagens, nil
}

// BuscarContagemPorID busca uma contagem específica
func (d *Database) BuscarContagemPorID(id string) (*ContagemLocal, error) {
	query := `
                SELECT ID, CODIGO_ITEM, DESCRICAO_ITEM, LOCAL, QUANTIDADE,
                           VOLUMES, DATA_CONTAGEM, USUARIO_CONTAGEM, STATUS
                FROM CST_CONTAGEM_LOCAL
                WHERE ID = ? AND STATUS = 'ATIVO'
        `

	var c ContagemLocal
	var rawID interface{}
	err := d.DB.QueryRow(query, id).Scan(
		&rawID,
		&c.CodigoItem,
		&c.DescricaoItem,
		&c.Local,
		&c.Quantidade,
		&c.Volumes,
		&c.DataContagem,
		&c.UsuarioContagem,
		&c.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contagem não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar contagem: %v", err)
	}

	// Converter GUID para string
	c.ID = fmt.Sprintf("%s", rawID)

	return &c, nil
}

// AtualizarContagem atualiza uma contagem existente
func (d *Database) AtualizarContagem(id string, contagem ContagemRequest) (*ContagemLocal, error) {
	query := `
                UPDATE CST_CONTAGEM_LOCAL 
                SET LOCAL = ?, QUANTIDADE = ?, VOLUMES = ?, USUARIO_CONTAGEM = ?
                WHERE ID = ? AND STATUS = 'ATIVO'
        `

	quantidade, err := strconv.ParseFloat(contagem.Quantidade, 64)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter quantidade: %v", err)
	}

	usuario := contagem.UsuarioContagem
	if usuario == "" {
		usuario = "Sistema"
	}

	result, err := d.DB.Exec(query, contagem.Local, quantidade, contagem.Volumes, usuario, id)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar contagem: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("contagem não encontrada ou não pode ser atualizada")
	}

	// Buscar e retornar a contagem atualizada
	return d.BuscarContagemPorID(id)
}

// ExcluirContagem marca uma contagem como inativa (soft delete)
func (d *Database) ExcluirContagem(id string) error {
	// Se ID for UNIQUEIDENTIFIER no SQL Server
	const query = `
		UPDATE CST_CONTAGEM_LOCAL
		SET STATUS = 'INATIVO'
		WHERE ID = CONVERT(uniqueidentifier, @p1)
		  AND STATUS = 'ATIVO';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := d.DB.ExecContext(ctx, query, sql.Named("p1", id))
	if err != nil {
		return fmt.Errorf("erro ao excluir contagem: %w", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if ra == 0 {
		return errors.New("contagem não encontrada ou já inativa")
	}
	return nil
}

// BuscarDescricaoItem busca a descrição do item na tabela ESTOQUE
func (d *Database) BuscarDescricaoItem(codigo string) (*ItemInfo, error) {
	const query = `
		SELECT TOP 1
			LTRIM(RTRIM(DESCRI)) AS DESCRICAO,
			ISNULL(CAST(QTDE_EMBALAGEM AS DECIMAL(18,4)), 0) AS QTDE_EMBALAGEM
		FROM ESTOQUE
		WHERE CODIGO = @p1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := &ItemInfo{Codigo: codigo}

	var (
		desc sql.NullString
		qtd  sql.NullFloat64
	)

	err := d.DB.QueryRowContext(ctx, query, sql.Named("p1", codigo)).Scan(&desc, &qtd)
	if err != nil {
		if err == sql.ErrNoRows {
			// Não achou: retorna defaults amigáveis
			return &ItemInfo{
				Codigo:        codigo,
				Descricao:     "Descrição não encontrada",
				QtdeEmbalagem: 0,
			}, nil
		}
		return nil, fmt.Errorf("erro ao buscar descrição do item: %w", err)
	}

	// Normaliza campos
	if desc.Valid {
		item.Descricao = strings.TrimSpace(desc.String)
	}
	if qtd.Valid {
		item.QtdeEmbalagem = qtd.Float64
	}

	// Garante defaults se vierem vazios
	if item.Descricao == "" {
		item.Descricao = "Descrição não disponível"
	}

	return item, nil
}

func (d *Database) ValidarLocal(sigla string) (*LocalInfo, error) {
	const query = `
		SELECT TOP 1
			CODIGO,
			LTRIM(RTRIM(SIGLA))      AS SIGLA,
			LTRIM(RTRIM(DESCRICAO))  AS DESCRICAO
		FROM LOCAIS
		WHERE SIGLA = @p1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		local       LocalInfo
		siglaNS     sql.NullString
		descricaoNS sql.NullString
	)

	err := d.DB.QueryRowContext(ctx, query, sql.Named("p1", sigla)).
		Scan(&local.Codigo, &siglaNS, &descricaoNS)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("local não encontrado")
		}
		return nil, fmt.Errorf("erro ao validar local: %w", err)
	}

	// Normaliza strings (tratando NULL e espaços)
	if siglaNS.Valid {
		local.Sigla = strings.TrimSpace(siglaNS.String)
	}
	if descricaoNS.Valid {
		local.Descricao = strings.TrimSpace(descricaoNS.String)
	}

	return &local, nil
}

func (d *Database) FinalizarContagem() error {
	const query = `
		UPDATE CST_CONTAGEM_LOCAL
		SET STATUS = 'FINALIZADO'
		WHERE STATUS = 'ATIVO';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := d.DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("erro ao finalizar contagem: %w", err)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if ra == 0 {
		return errors.New("nenhuma contagem ativa encontrada para finalizar")
	}
	return nil
}

// VerificarContagemFinalizada verifica se existem contagens finalizadas
func (d *Database) VerificarContagemFinalizada() (bool, error) {
	const query = `
		SELECT COUNT(*) 
		FROM CST_CONTAGEM_LOCAL 
		WHERE STATUS = 'FINALIZADO' 
		 and CONVERT(VARCHAR,DATA_CONTAGEM,112) = CONVERT(VARCHAR,GETDATE(),112);
	`

	// Define timeout de 5 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int64
	err := d.DB.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar contagens finalizadas: %w", err)
	}

	return count > 0, nil
}
