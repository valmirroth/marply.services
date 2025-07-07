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

// BuscarContagemAtiva busca a contagem ATIVA pela SIGLA (LOCAIS.SIGLA).
// Se não existir, cria um registro em CST_REGISTRO_CONTAGEM_LOCAL (STATUS='ATIVO')
// para o local correspondente e retorna o CODIGO do local.
// Requer imports: context, database/sql, errors, fmt, time
func (d *Database) BuscarContagemAtiva(empresa string) (string, error) {
	_ = empresa // ainda não utilizado neste SQL

	selectAtiva := `
		SELECT TOP (1) CAST(rc.CODIGO AS VARCHAR(50)) AS Contagem
		FROM CST_REGISTRO_CONTAGEM_LOCAL rc
		WHERE rc.STATUS = 'ATIVO'
			AND EMPRESA_RECNO = '` + empresa + `'
		ORDER BY rc.CODIGO DESC;
	`

	// ATENÇÃO: ajuste as colunas caso haja NOT NULL adicionais.
	insertContagem := `
		INSERT INTO CST_REGISTRO_CONTAGEM_LOCAL (CODIGO, STATUS, EMPRESA_RECNO)
		SELECT ISNULL(MAX(ISNULL(CODIGO,0)),0) + 1, 'ATIVO', '` + empresa + `'
		FROM CST_REGISTRO_CONTAGEM_LOCAL;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var contagem sql.NullString

	// 1) Tenta achar contagem ativa
	if err := d.DB.QueryRowContext(ctx, selectAtiva).Scan(&contagem); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("buscando contagem ativa: %w", err)
	}
	if contagem.Valid && contagem.String != "" {
		return contagem.String, nil
	}

	// 2) Não existe -> cria uma ativa
	if _, err := d.DB.ExecContext(ctx, insertContagem); err != nil {
		return "", fmt.Errorf("inserindo contagem ativa: %w", err)
	}

	// 3) Busca novamente
	contagem = sql.NullString{}
	if err := d.DB.QueryRowContext(ctx, selectAtiva).Scan(&contagem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("nenhum registro encontrado mesmo após a inserção")
		}
		return "", fmt.Errorf("buscando contagem ativa após inserção: %w", err)
	}
	if contagem.Valid && contagem.String != "" {
		return contagem.String, nil
	}

	return "", fmt.Errorf("falha ao obter contagem após a inserção")
}

func (d *Database) BuscarSaldo(codigo, sigla string) (float64, error) {
	const query = `
		SELECT isnull(SUM(el.qtde),0) AS Saldo
		FROM dbo.ESTOQUE_LOCAL el
		INNER JOIN dbo.LOCAIS l ON l.CODIGO = el.LOCAL
		WHERE el.CODIGO = @codigo
		  AND l.SIGLA = @sigla;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var saldo sql.NullFloat64

	err := d.DB.QueryRowContext(ctx, query,
		sql.Named("codigo", codigo),
		sql.Named("sigla", sigla),
	).Scan(&saldo)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("nenhum registro encontrado para CODIGO=%s e SIGLA=%s", codigo, sigla)
		}
		return 0, fmt.Errorf("erro ao buscar saldo: %w", err)
	}

	if saldo.Valid {
		return saldo.Float64, nil
	}

	return 0, nil
}

// InserirContagem insere uma nova contagem no banco
func (d *Database) InserirContagem(contagem ContagemRequest) (*ContagemLocal, error) {
	// SQL Server: use @p1..@pN e OUTPUT para capturar os valores gerados
	const query = `
		INSERT INTO CST_CONTAGEM_LOCAL 
			(CODIGO_ITEM, DESCRICAO_ITEM, LOCAL, QUANTIDADE, VOLUMES, USUARIO_CONTAGEM, SALDO_ESTOQUE, NROCONTAGEM, CODLOCAL, EMPRESA_RECNO)
		OUTPUT INSERTED.ID, INSERTED.DATA_CONTAGEM, INSERTED.STATUS
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9, @p10);
	`

	// Converter quantidade (string -> float64)
	quantidade, err := strconv.ParseFloat(contagem.Quantidade, 64)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter quantidade: %w", err)
	}

	volumes, errv := strconv.ParseFloat(contagem.Volumes, 64)
	if errv != nil {
		return nil, fmt.Errorf("erro ao converter volumes: %w", errv)
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

	saldo, errosaldo := d.BuscarSaldo(contagem.CodigoItem, contagem.Local)
	if errosaldo != nil {
		log.Printf("Aviso: não foi possível buscar saldo para item %s no local %s: %v", contagem.CodigoItem, contagem.Local, errosaldo)
		saldo = 0 // ou outro valor padrão
		return nil, fmt.Errorf("erro ao inserir contagem: %w", errosaldo)
	}

	nrocontagem, errocontagem := d.BuscarContagemAtiva(strconv.Itoa(contagem.Empresa))
	if errocontagem != nil {
		log.Printf("Aviso: não foi possível buscar contagem ativa: %v", errocontagem)
		nrocontagem = "" // ou outro valor padrão
		return nil, fmt.Errorf("erro ao inserir contagem: %w", errocontagem)
	}

	xlocal, erro := d.ValidarLocal(contagem.Local, strconv.Itoa(contagem.Empresa))
	if erro != nil {
		log.Printf("Aviso: local %s não encontrado: %v", contagem.Local, erro)
		return nil, fmt.Errorf("erro ao inserir contagem: %w", erro)
	}
	contagem.Local = xlocal.Sigla
	// Executa o INSERT e lê o OUTPUT
	if err := d.DB.
		QueryRowContext(ctx, query,
			contagem.CodigoItem,
			contagem.DescricaoItem,
			contagem.Local,
			quantidade,
			volumes,
			usuario,
			saldo,
			nrocontagem,
			xlocal.Codigo,
			contagem.Empresa,
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
		Volumes:         volumes,
		DataContagem:    dataContagem,
		UsuarioContagem: usuario,
		Status:          status,
		NroContagem:     nrocontagem,
		Empresa:         contagem.Empresa,
	}

	return resultado, nil
}

// ListarContagens retorna todas as contagens
func (d *Database) ListarContagens(empresa string) ([]ContagemLocal, error) {
	query := `
                SELECT cast(ID as varchar(150)) as ID, CODIGO_ITEM, DESCRICAO_ITEM, LOCAL, QUANTIDADE, 
                           VOLUMES, DATA_CONTAGEM, USUARIO_CONTAGEM, STATUS
                FROM CST_CONTAGEM_LOCAL
                WHERE STATUS = 'ATIVO'
					and EMPRESA_RECNO = '` + empresa + `'
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

func (d *Database) BuscarContagemPorID(id string) (*ContagemLocal, error) {
	const query = `
		SELECT
			CAST(ID AS varchar(36)) AS ID,
			CODIGO_ITEM,
			DESCRICAO_ITEM,
			[LOCAL],
			ISNULL(QUANTIDADE, 0),
			ISNULL(VOLUMES, 0),
			DATA_CONTAGEM,
			USUARIO_CONTAGEM,
			STATUS
		FROM dbo.CST_CONTAGEM_LOCAL WITH (NOLOCK)
		WHERE ID = @id AND STATUS = 'ATIVO';
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := d.DB.QueryRowContext(ctx, query, sql.Named("id", id))

	var (
		idStr        string
		qtd          sql.NullFloat64
		vol          sql.NullFloat64
		dataContagem sql.NullTime
		c            ContagemLocal
	)

	if err := row.Scan(
		&idStr,
		&c.CodigoItem,
		&c.DescricaoItem,
		&c.Local,
		&qtd,
		&vol,
		&dataContagem,
		&c.UsuarioContagem,
		&c.Status,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contagem não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar contagem: %w", err)
	}

	c.ID = idStr
	if qtd.Valid {
		c.Quantidade = qtd.Float64
	}
	if vol.Valid {
		c.Volumes = (vol.Float64)
	}
	if dataContagem.Valid {
		c.DataContagem = dataContagem.Time
	}

	return &c, nil
}

// AtualizarContagem atualiza uma contagem existente
func (d *Database) AtualizarContagem(id string, contagem ContagemRequest) (*ContagemLocal, error) {
	// SQL Server: use parâmetros nomeados e escape do campo [LOCAL]
	const query = `
		UPDATE dbo.CST_CONTAGEM_LOCAL
		SET
			QUANTIDADE = @quantidade,
			VOLUMES = @volumes,
			USUARIO_CONTAGEM = @usuario
		WHERE ID = @id
		  AND STATUS = 'ATIVO';
	`

	// Converter quantidade aceitando vírgula decimal (pt-BR)
	qtdStr := strings.ReplaceAll(strings.TrimSpace(contagem.Quantidade), ",", ".")
	quantidade, err := strconv.ParseFloat(qtdStr, 64)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter quantidade '%s': %w", contagem.Quantidade, err)
	}

	usuario := strings.TrimSpace(contagem.UsuarioContagem)
	if usuario == "" {
		usuario = "Sistema"
	}

	// Timeout de proteção
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execução com parâmetros nomeados (sqlserver)
	result, execErr := d.DB.ExecContext(
		ctx, query,
		sql.Named("local", contagem.Local),
		sql.Named("quantidade", quantidade),
		sql.Named("volumes", contagem.Volumes),
		sql.Named("usuario", usuario),
		sql.Named("id", id),
	)
	if execErr != nil {
		return nil, fmt.Errorf("erro ao atualizar contagem: %w", execErr)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("contagem não encontrada ou não pode ser atualizada (ID: %s ou STATUS != 'ATIVO')", id)
	}

	// Retorna a versão atualizada
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
			ISNULL(CAST(EMB.CST_METRAGEM AS DECIMAL(18,4)), 0) AS QTDE_EMBALAGEM,
			case when estoque.familia in (3,4,16,166) then 1 else 0 end as CALCULAR_AUTOMATICO
		FROM ESTOQUE
			left join CAD_EMBALAGEM EMB ON EMB.R_E_C_N_O_ = ESTOQUE.RECNO_EMBALAGEM
			WHERE ESTOQUE.CODIGO = @p1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	item := &ItemInfo{Codigo: codigo}

	var (
		desc sql.NullString
		qtd  sql.NullFloat64
		calc sql.NullBool
	)

	err := d.DB.QueryRowContext(ctx, query, sql.Named("p1", codigo)).Scan(&desc, &qtd, &calc)
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

	if calc.Valid {
		item.CalcularAutomatico = calc.Bool
	}

	// Garante defaults se vierem vazios
	if item.Descricao == "" {
		item.Descricao = "Descrição não disponível"
	}

	return item, nil
}

func (d *Database) ValidarLocal(sigla, empresa string) (*LocalInfo, error) {
	var query = `
			SELECT TOP 1
				CODIGO,
				LTRIM(RTRIM(SIGLA))      AS SIGLA,
				'End.:'+ LTRIM(RTRIM(SIGLA)) +' local:'+ LTRIM(RTRIM(DESCRICAO)) +' '+ ' Itens: ' +
					isnull((	select STRING_AGG(codigo,' | ') from ESTOQUE_LOCAL 
						where LOCAL = LOCAIS.CODIGO 
						AND ESTOQUE_LOCAL.CODIGO NOT IN (
							SELECT CST_CONTAGEM_LOCAL.CODIGO_ITEM FROM CST_CONTAGEM_LOCAL
								INNER JOIN CST_REGISTRO_CONTAGEM_LOCAL ON CST_REGISTRO_CONTAGEM_LOCAL.CODIGO = CST_CONTAGEM_LOCAL.NROCONTAGEM
							WHERE CST_REGISTRO_CONTAGEM_LOCAL.STATUS = 'ATIVO'
									AND CST_CONTAGEM_LOCAL.CODLOCAL = LOCAIS.CODIGO
						)
					),'') as DESCRICAO
			FROM LOCAIS
			WHERE (SIGLA = @p1 or  DESCRICAO = @p1)
				and EMPRESA_RECNO = '` + empresa + `' ;
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

/*
func (d *Database) FinalizarContagem(empresa string) error {
	const query = `

		insert into CST_CONTAGEM_LOCAL(CODIGO_ITEM, DESCRICAO_ITEM, CODLOCAL, LOCAL, SALDO_ESTOQUE,  VOLUMES, QUANTIDADE, NROCONTAGEM)
		SELECT E.CODIGO, E.DESCRI, LOCAIS.CODIGO, LOCAIS.SIGLA,   SUM(ESTOQUE_LOCAL.QTDE) AS QTDE,
				0 Volumes , 0,
			(SELECT CODIGO FROM CST_REGISTRO_CONTAGEM_LOCAL WHERE STATUS = 'ATIVO')
			FROM ESTOQUE_LOCAL
			INNER JOIN LOCAIS ON LOCAIS.CODIGO = ESTOQUE_LOCAL.LOCAL
			INNER JOIN ESTOQUE E ON E.CODIGO = ESTOQUE_LOCAL.CODIGO
			LEFT JOIN CAD_EMBALAGEM EMB ON EMB.R_E_C_N_O_ = E.RECNO_EMBALAGEM
			LEFT JOIN CST_CONTAGEM_LOCAL ON CST_CONTAGEM_LOCAL.CODIGO_ITEM = E.CODIGO
				AND CST_CONTAGEM_LOCAL.CODLOCAL = LOCAIS.CODIGO
				AND CST_CONTAGEM_LOCAL.STATUS = 'ATIVA'
		WHERE LOCAIS.EMPRESA_RECNO = '1' AND
			LOCAIS.CODIGO NOT IN (1,5,8,10,15,19,3,6,21,11,12)
			AND CST_CONTAGEM_LOCAL.CODIGO_ITEM IS NULL
			and e.CATEGORIA = '99'
			and E.FAMILIA in (3,4,16,166,2,15)
		GROUP BY E.CODIGO, E.DESCRI, LOCAIS.SIGLA, LOCAIS.CODIGO, EMB.CST_METRAGEM;

		insert into CST_HIST_CONTAGEM_LOCAL(CODIGO_ITEM, DESCRICAO_ITEM, CODLOCAL, LOCAL, SALDO_ESTOQUE,  VOLUMES, NROCONTAGEM)
		SELECT E.CODIGO, E.DESCRI, LOCAIS.CODIGO, LOCAIS.SIGLA,   SUM(ESTOQUE_LOCAL.QTDE) AS QTDE,
				0 Volumes ,
			(SELECT CODIGO FROM CST_REGISTRO_CONTAGEM_LOCAL WHERE STATUS = 'ATIVO')
			FROM ESTOQUE_LOCAL
			INNER JOIN LOCAIS ON LOCAIS.CODIGO = ESTOQUE_LOCAL.LOCAL
			INNER JOIN ESTOQUE E ON E.CODIGO = ESTOQUE_LOCAL.CODIGO
		WHERE LOCAIS.EMPRESA_RECNO = '1' AND
			LOCAIS.CODIGO NOT IN (1,5,8,10,15,19,3,6,21,11,12)
			and e.CATEGORIA = '99'
			and E.FAMILIA in (3,4,16,166,2,15)
		GROUP BY E.CODIGO, E.DESCRI, LOCAIS.SIGLA, LOCAIS.CODIGO ;

		UPDATE CST_CONTAGEM_LOCAL
		SET STATUS = 'FINALIZADO'
		WHERE STATUS = 'ATIVO';

		UPDATE CST_REGISTRO_CONTAGEM_LOCAL
			SET FINAL_CONTAGEM = GETDATE(), STATUS = 'FINALIZADA'
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
*/

func (d *Database) FinalizarContagem(empresa string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("iniciando transação: %w", err)
	}
	defer func() { _ = tx.Rollback() }() // rollback é no-op após Commit

	// 1) Captura a contagem ativa
	qAtiva := `
			SELECT TOP (1) CAST(CODIGO AS VARCHAR(50))
			FROM CST_REGISTRO_CONTAGEM_LOCAL
			WHERE STATUS = 'ATIVO'
			AND EMPRESA_RECNO =  '` + empresa + `'
			ORDER BY CODIGO DESC;
		`
	var nroContagem string
	if err := tx.QueryRowContext(ctx, qAtiva).Scan(&nroContagem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("nenhuma contagem ativa encontrada para finalizar")
		}
		return fmt.Errorf("buscando contagem ativa: %w", err)
	}

	// 2) INSERT em CST_CONTAGEM_LOCAL, evitando duplicatas nesta mesma contagem
	const qInsertAtual = `
		INSERT INTO CST_CONTAGEM_LOCAL
			(CODIGO_ITEM, DESCRICAO_ITEM, CODLOCAL, LOCAL, SALDO_ESTOQUE, VOLUMES, QUANTIDADE, NROCONTAGEM, EMPRESA_RECNO)
		SELECT
			E.CODIGO,
			E.DESCRI,
			LOCAIS.CODIGO,
			LOCAIS.SIGLA,
			SUM(ESTOQUE_LOCAL.QTDE) AS QTDE,
			0 AS VOLUMES,
			0 AS QUANTIDADE,
			@NroContagem,
			@Empresa
		FROM ESTOQUE_LOCAL
		INNER JOIN LOCAIS  ON LOCAIS.CODIGO = ESTOQUE_LOCAL.LOCAL
		INNER JOIN ESTOQUE E ON E.CODIGO = ESTOQUE_LOCAL.CODIGO
		LEFT JOIN CST_CONTAGEM_LOCAL C ON
			C.CODIGO_ITEM  = E.CODIGO
			AND C.CODLOCAL = LOCAIS.CODIGO
			AND C.NROCONTAGEM = @NroContagem
			AND C.STATUS = 'ATIVO'
		WHERE
			LOCAIS.EMPRESA_RECNO = @Empresa
			AND LOCAIS.CODIGO NOT IN (1,5,8,10,15,19,3,6,21,11,12)
			AND C.CODIGO_ITEM IS NULL
			AND E.CATEGORIA = '99'
			AND E.FAMILIA IN (3,4,16,166,2,15)
		GROUP BY E.CODIGO, E.DESCRI, LOCAIS.SIGLA, LOCAIS.CODIGO;
	`
	if _, err := tx.ExecContext(
		ctx,
		qInsertAtual,
		sql.Named("NroContagem", nroContagem),
		sql.Named("Empresa", empresa),
	); err != nil {
		return fmt.Errorf("inserindo CST_CONTAGEM_LOCAL: %w", err)
	}

	// 3) INSERT em histórico
	const qInsertHist = `
		INSERT INTO CST_HIST_CONTAGEM_LOCAL
			(CODIGO_ITEM, DESCRICAO_ITEM, CODLOCAL, LOCAL, SALDO_ESTOQUE, VOLUMES, NROCONTAGEM, EMPRESA_RECNO)
		SELECT
			E.CODIGO,
			E.DESCRI,
			LOCAIS.CODIGO,
			LOCAIS.SIGLA,
			SUM(ESTOQUE_LOCAL.QTDE) AS QTDE,
			0 AS VOLUMES,
			@NroContagem,
			@Empresa
		FROM ESTOQUE_LOCAL
		INNER JOIN LOCAIS  ON LOCAIS.CODIGO = ESTOQUE_LOCAL.LOCAL
		INNER JOIN ESTOQUE E ON E.CODIGO = ESTOQUE_LOCAL.CODIGO
		WHERE
			LOCAIS.EMPRESA_RECNO = @Empresa
			AND LOCAIS.CODIGO NOT IN (1,5,8,10,15,19,3,6,21,11,12)
			AND E.CATEGORIA = '99'
			AND E.FAMILIA IN (3,4,16,166,2,15)
		GROUP BY E.CODIGO, E.DESCRI, LOCAIS.SIGLA, LOCAIS.CODIGO;
	`
	if _, err := tx.ExecContext(
		ctx,
		qInsertHist,
		sql.Named("NroContagem", nroContagem),
		sql.Named("Empresa", empresa),
	); err != nil {
		return fmt.Errorf("inserindo CST_HIST_CONTAGEM_LOCAL: %w", err)
	}

	// 4) Finaliza itens ATIVOS somente desta contagem
	const qFinalizaItens = `
		UPDATE CST_CONTAGEM_LOCAL
		SET STATUS = 'FINALIZADO'
		WHERE STATUS = 'ATIVO'
		  AND NROCONTAGEM = @NroContagem;
	`
	resItens, err := tx.ExecContext(ctx, qFinalizaItens, sql.Named("NroContagem", nroContagem))
	if err != nil {
		return fmt.Errorf("finalizando itens da contagem: %w", err)
	}
	itensAfetados, _ := resItens.RowsAffected()

	// 5) Finaliza o registro da contagem ativa
	const qFinalizaRegistro = `
		UPDATE CST_REGISTRO_CONTAGEM_LOCAL
		SET FINAL_CONTAGEM = GETDATE(), STATUS = 'FINALIZADA'
		WHERE CODIGO = @NroContagem
		  AND STATUS = 'ATIVO';
	`
	resReg, err := tx.ExecContext(ctx, qFinalizaRegistro, sql.Named("NroContagem", nroContagem))
	if err != nil {
		return fmt.Errorf("finalizando registro da contagem: %w", err)
	}
	regAfetados, _ := resReg.RowsAffected()

	if itensAfetados == 0 && regAfetados == 0 {
		// Nada alterado — evita commit inútil
		return errors.New("nenhuma linha alterada ao finalizar a contagem")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// VerificarContagemFinalizada verifica se existem contagens finalizadas
func (d *Database) VerificarContagemFinalizada(empresa string) (bool, error) {
	const query = `
		SELECT COUNT(*) 
		FROM CST_CONTAGEM_LOCAL 
		WHERE STATUS = 'FINALIZADOW' 
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
