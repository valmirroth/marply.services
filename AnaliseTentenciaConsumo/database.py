import pyodbc
from config import SQL_SERVER, SQL_DATABASE, SQL_USER, SQL_PASSWORD

QUERY_TEMPLATE = r"""
with resumo as (
 SELECT * FROM ( SELECT EMPRESA_RECNO, codconta,
   case when 
		case when [Observação] like '%requisição%' then   LTRIM(
			SUBSTRING(
				[Observação],
				CHARINDEX('- ', [Observação], CHARINDEX('Nº', [Observação])) + 2,
				LEN([Observação])
			)
		) else [Observação] end like 'Ajuste Provisão IND 13.Sal -%' then 'Ajuste Provisão IND 13.Sal'
		when [Observação] like 'Valor Provisão 13.Sal -%' then 'Valor Provisão 13.Sal'
		when [Observação] like '1/3 Provisão Férias - %' then '1/3 Provisão Férias'
		when [Observação] like 'Valor Provisão Férias - %' then 'Valor Provisão Férias'
		when [Observação] like 'Ajuste FGTS 13.Sal - %' then 'Ajuste FGTS 13.Sal'
        when [Observação] like 'Baixa Prov.FGTS 13.Sal - %' then 'Baixa Prov.FGTS 13.Sal'
        when [Observação] like 'FGTS - Folha %' then 'FGTS - Folha'
        when [Observação] like 'Provisão FGTS 13.Sal - %' then 'Provisão FGTS 13.Sal'
        when [Observação] like 'FGTS MULTA - Folha %' then 'FGTS MULTA - Folha'
        when [Observação] like 'Ajuste INSS 13_o Sal - %' then 'Ajuste INSS 13_o Sal'
        when [Observação] like 'Baixa Prov.INSS 13_o Sal - %' then 'Baixa Prov.INSS 13º Sal'
        when [Observação] like 'Parte Convênios GPS - Folha%' then 'Parte Convênios GPS'
        when [Observação] like 'Parte Terceiros GPS - Folha %' then 'Parte Terceiros GPS - Folha'
        when [Observação] like 'Provisão INSS 13_o Sal - %' then 'Provisão INSS 13º Sal'
        when [Observação] like 'Provisão INSS Férias - %' then 'Provisão INSS Férias'
        when [Observação] like 'ADICIONAL NOTURNO - Folha%' then 'ADICIONAL NOTURNO - Folha'
        when [Observação] like 'AVISO PRÉVIO INDENIZADO - Folha %' then 'AVISO PRÉVIO INDENIZADO - Folha'
        when [Observação] like 'AVISO PREVIO REAVIDO - Folha%' then 'AVISO PREVIO REAVIDO - Folha'
        when [Observação] like 'BONIFICAÇÂO - Folha %' then 'BONIFICAÇÂO - Folha'
        when [Observação] like 'DSR HORAS EXTRAS - Folha%' then 'DSR HORAS EXTRAS - Folha'
        when [Observação] like 'HORAS ATESTADO - Folha%' then 'HORAS ATESTADO - Folha'
        when [Observação] like 'ESTOURO DO MÊS - Folha%' then 'ESTOURO DO MÊS - Folha'
        when [Observação] like 'ESTOURO MES ANTERIOR - Folha%' then 'ESTOURO MES ANTERIOR - Folha'
        when [Observação] like 'HORAS EXTRAS - Folha%' then 'HORAS EXTRAS - Folha'
        when [Observação] like 'HORAS EXTRAS INDENIZADAS - Folha%' then 'HORAS EXTRAS INDENIZADAS - Folha'
        when [Observação] like 'HORAS FALTAS - Folha%' then 'HORAS FALTAS - Folha'
        when [Observação] like 'HORAS NORMAIS - Folha%' then 'HORAS NORMAIS - Folha'
        when [Observação] like 'SALDO DE SALÁRIO - Folha%' then 'SALDO DE SALÁRIO - Folha'
        when [Observação] like '%requisição%' then   LTRIM(
			SUBSTRING(
				[Observação],
				CHARINDEX('- ', [Observação], CHARINDEX('Nº', [Observação])) + 2,
				LEN([Observação])
			)
		)
		else [Observação]
	end AS [Observação],
	 case when [Observação] like '%requisição%' then 'Baixas Almox' else 'Outos Movs' end as TipoMovimento,
     CST_BI_ANALISE_CUSTEIO_GERAL.Periodo,
     CST_BI_ANALISE_CUSTEIO_GERAL.Calculado
	FROM CST_BI_ANALISE_CUSTEIO_GERAL
    
    UNION ALL 

    SELECT HF.EMPRESA_RECNO, HISTLISE.ITEM AS codconta, E.CODIGO +' - '+ E.DESCRI AS [Observação], 'Baixas Almox' AS TipoMovimento, 
        hf.DTENT as Periodo, HISTLISE.TOTAL Calculado
        FROM HISTLISE 
        inner join estoque E on E.CODIGO = HISTLISE.ITEM    AND E.CATEGORIA = '30'
	    INNER JOIN HISTLISE_FOR HF ON HF.R_E_C_N_O_ = HISTLISE.RECNO_NOTA
	    LEFT JOIN HISREAL HL ON HL.RECNO_HISTLISE = HISTLISE.R_E_C_N_O_
        INNER JOIN HISTLISE_IMP IMP ON IMP.RECNO_ITEM = HISTLISE.R_E_C_N_O_
    WHERE   HL.R_E_C_N_O_ IS NULL
	    AND HISTLISE.EFETUADOENTR = 'S'
        AND HF.DTENT >= '20250101'

    ) FECHA

	WHERE EMPRESA_RECNO IN ({empresa_filter})
)
select Periodo, GRUPOE.DESGRUPO as ContaContabil, e.codigo +' - '+e.descri as [Observação], sum(calculado)as Valor 
from resumo
    left join estoque e on e.codigo =  LTRIM(RTRIM(LEFT(
        [Observação],
        CHARINDEX('-', [Observação]) - 1
    )))  
    left join GRUPOE on GRUPOE.GRUPO = e.FAMILIA
where TipoMovimento = 'Baixas Almox'
  and Periodo between '20250101' and case when 'sim' = '{mes_atual}' then EOMONTH(getdate()) else EOMONTH(dateadd(month,-1,getdate())) end 
  and categoria <> 99 
group by GRUPOE.DESGRUPO, Periodo, e.codigo +' - '+e.descri  
"""

def get_conn():
    """Retorna conexão com SQL Server"""
    return pyodbc.connect(
        "DRIVER={ODBC Driver 17 for SQL Server};"
        f"SERVER={SQL_SERVER};"
        f"DATABASE={SQL_DATABASE};"
        f"UID={SQL_USER};"
        f"PWD={SQL_PASSWORD};"
        "TrustServerCertificate=yes;"
    )

def get_empresa_filter(empresa_tipo):
    """Retorna filtro de empresa"""
    from config import EMPRESAS
    return EMPRESAS.get(empresa_tipo.lower(), '1,5')
