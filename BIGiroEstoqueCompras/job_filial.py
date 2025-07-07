import logging
import os
import pandas as pd
import pyodbc
from datetime import datetime

LOG_PATH = os.getenv("GIRO_LOG_PATH", r"C:\MarplyServices\BI - AnaliseGiroEstoque\logs\giro_nssm.log")
os.makedirs(os.path.dirname(LOG_PATH), exist_ok=True)
logging.basicConfig(filename=LOG_PATH, level=logging.INFO,
                    format="%(asctime)s %(levelname)s %(message)s")

def converter_formato_data(data_str):
    formato_origem = "%Y%m%d"  # Formato da data de entrada
    formato_destino = "%d/%m/%Y"  # Formato desejado na saída
    data_objeto = datetime.strptime(data_str, formato_origem)
    return data_objeto.strftime(formato_destino)

def calcular_intervalo(dia_inicio, dia_fim):
    dataconvertida = converter_formato_data(dia_inicio)
    formato = "%d/%m/%Y"  # Define o formato das datas
    data_inicio = datetime.strptime(dataconvertida, formato)
    data_fim = datetime.strptime(dia_fim, formato)
    intervalo = data_fim - data_inicio
    return intervalo.days

# Função para calcular o intervalo médio (em dias) entre uma lista de datas
def compute_avg_interval(dates):
    if len(dates) < 2:
        return None
    # As datas já devem estar em ordem crescente
    diffs = [(dates[i+1] - dates[i]).days for i in range(len(dates)-1)]
    return sum(diffs) / len(diffs)

def contar_notas(item: str) -> int:
    # Conexão com o SQL Server (ajuste conforme seu ambiente)
    conn_str = (
        "DRIVER={ODBC Driver 17 for SQL Server};"
        "SERVER=192.168.1.28;"          # Substitua pelo endereço do seu servidor
        "DATABASE=Marini_PRD;"          # Substitua pelo nome do seu banco de dados
        "UID=Octopus;"                  # Substitua pelo seu usuário
        "PWD=A45182008069199;"          # Substitua pela sua senha
    )
    conn = None
    conn = pyodbc.connect(conn_str)
    cursor = conn.cursor()

    query = """
    SELECT COUNT(R_E_C_N_O_) AS NUMNOTAS 
    FROM (
        SELECT hf.R_E_C_N_O_ 
        FROM HISTLISE_FOR hf 
            INNER JOIN HISTLISE h ON h.RECNO_NOTA = hf.R_E_C_N_O_
        WHERE hf.DTENT BETWEEN '20250101' AND (SELECT CONVERT(varchar, DATEADD(day,-1,GETDATE()), 112))
            AND hf.EFETUADOENTR = 's'
            AND h.ITEM = ?
            and hf.EMPRESA_RECNO = 5
        GROUP BY hf.R_E_C_N_O_
    ) RESUMO
    """

    cursor.execute(query, (item,))
    result = cursor.fetchone()
    cursor.close()
    conn.close()

    return result[0] if result else 0

def run_job():
    logging.info("Iniciando job de análise de giro de estoque")
    # Configuração da conexão com o SQL Server
    # Dica: se possível, troque para ODBC Driver 18 e ajuste Encrypt/Trust conforme política
    conn_str = (
        "DRIVER={ODBC Driver 17 for SQL Server};"
        "SERVER=192.168.1.28;"          # Substitua pelo endereço do seu servidor
        "DATABASE=Marini_PRD;"          # Substitua pelo nome do seu banco de dados
        "UID=Octopus;"                  # Substitua pelo seu usuário
        "PWD=A45182008069199;"          # Substitua pela sua senha
    )
    conn = None
    try:
        conn = pyodbc.connect(conn_str)
        # Consulta SQL para obter os dados
        query = """ declare @di varchar(10)
            declare @df varchar(10)
            set @di = '20250101'
            set @df = (select CONVERT(varchar,dateadd(day,-1,getdate()),112))
              select  '20241231' as DTRECEB,
                    0 AS Linha, 
                    0 R_E_C_N_O_, 
                    CLASSE.DESCLASSE AS Categoria,
                    ESTOQUE_POSICAO.ITEM as CODIGO, 
                    ESTOQUE.DESCRI as DESCRICAO, GRUPOE.DESGRUPO AS FAMILIA,
                    G.DESCRICAO AS GRUPOESTOQUE,
                    0 as QtdeMovimento,
                    'Compra' as Tipo, 
                    'Compra'   DOCUMEN,
                    ESTOQUE_POSICAO.SALDOREAL as SALDO_EMPRESA, 
                    ESTOQUE_POSICAO.SALDOREAL * ESTOQUE_POSICAO.VALOR as CUSTO_TOTAL,
                    ESTOQUE_POSICAO.SALDOREAL as EstoqueMedio,
                    ESTOQUE_POSICAO.VALOR as CUSTO_MEDIO
            from ESTOQUE_POSICAO 
              
                INNER JOIN ESTOQUE ON ESTOQUE.CODIGO = ESTOQUE_POSICAO.ITEM  and ESTOQUE.CATEGORIA not in ( '99','10','110','30')
                LEFT JOIN CLASSE ON CLASSE.CLASSE = ESTOQUE.CATEGORIA
                LEFT JOIN GRUPOE ON GRUPOE.GRUPO = ESTOQUE.FAMILIA
                left join GRUPO_EST G ON G.CODIGO = ESTOQUE.CODGRUPO
                WHERE   EMPRESA_RECNO = 5 --and ITEM = '000292'
                and ESTOQUE_POSICAO.DATAHORA = '20241231'
                and  ESTOQUE_POSICAO.SALDOREAL > 0
               AND ESTOQUE_POSICAO.ITEM NOT IN (select distinct CODIGO FROM HISREAL WHERE DTRECEB >= @di  and DTRECEB <= @df AND EMPRESA_RECNO = 5 AND MOV_EFETIVA = 'S')
    UNION ALL 
              select   HISREAL.DTRECEB, 
                    ROW_NUMBER() over(PARTITION BY HISREAL.CODIGO order by DTRECEB asc, HISREAL.R_E_C_N_O_ ASC) Linha, 
                    HISREAL.R_E_C_N_O_, 
                    CLASSE.DESCLASSE AS Categoria,
                    HISREAL.CODIGO, 
                    HISREAL.DESCRICAO, GRUPOE.DESGRUPO AS FAMILIA,
                    G.DESCRICAO AS GRUPOESTOQUE,
                    HISREAL.QTRECEB as QtdeMovimento,
                    case when RECNO_HISTLISE IS NOT null then 'Compra' 
                        when FORMA ='s' and (RECNO_PCP_ITENS_REQ IS  NOT null OR RECNO_ITENNOTA IS NOt null) then 'Consumo'
                        when FORMA = 'e' and DOCUMEN like 'devolução %' then 'Devolvido'
                        WHEN DOCUMEN = 'INV.LOCAL' AND FORMA = 'E' THEN 'Devolvido'
                        WHEN DOCUMEN = 'INV.LOCAL' AND FORMA = 'S' THEN 'Consumo'
                        WHEN DOCUMEN = 'AJUSTE DE ESTOQUE' AND FORMA = 'E' THEN 'Devolvido'
                        WHEN DOCUMEN = 'AJUSTE DE ESTOQUE' AND FORMA = 'S' THEN 'Consumo' 
                        end as Tipo, 
                        DOCUMEN,
                    HISREAL.SALDO_EMPRESA, 
                    HISREAL.CUSTO_TOTAL,
                    (
                    select SUM(SALDOREAL)/(DATEDIFF(MONTH,
                        ( SELECT MIN(Z.DATAHORA) FROM   ESTOQUE_POSICAO Z
                          where Z.ITEM = ESTOQUE_POSICAO.ITEM and Z.DATAHORA = EOMONTH(Z.datahora) 
                          AND Z.EMPRESA_RECNO = 5
                          and Z.DATAHORA between   @di and @df  AND Z.SALDOREAL > 0 
                          )
                        ,@df)+1) 
                        from ESTOQUE_POSICAO where ITEM = HISREAL.CODIGO and DATAHORA = EOMONTH(datahora)
                        and DATAHORA between @di and @df 
                        and EMPRESA_RECNO = 5
                        AND SALDOREAL > 0
                    GROUP BY ITEM
                    ) as EstoqueMedio,
                    HISREAL.CUSTO_MEDIO
            from HISREAL 
                inner join ESTOQUE_HIST_MOVIMENTACAO_EFETIVA on ESTOQUE_HIST_MOVIMENTACAO_EFETIVA.RECNO_HISREAL = HISREAL.R_E_C_N_O_ and ESTOQUE_HIST_MOVIMENTACAO_EFETIVA.ORIGEM_MOVIMENTACAO <> 'InspecaoEntrada'
                INNER JOIN ESTOQUE ON ESTOQUE.CODIGO = HISREAL.CODIGO  and ESTOQUE.CATEGORIA not in ( '99','10','110')
                INNER JOIN CLASSE ON CLASSE.CLASSE = ESTOQUE.CATEGORIA
                LEFT JOIN GRUPOE ON GRUPOE.GRUPO = ESTOQUE.FAMILIA
                left join GRUPO_EST G ON G.CODIGO = ESTOQUE.CODGRUPO
                WHERE DTRECEB BETWEEN   @di and @df 
                AND HISREAL.EMPRESA_RECNO = 5
                AND MOV_EFETIVA = 'S'           
                and FORMA in  ('e','s')  
                                and HISREAL.CODIGO in (select codigo from (
                                select ROW_NUMBER()over(partition by codigo order by dtreceb desc, r_e_c_n_o_ desc) linha, z.CODIGO, SALDO_EMPRESA from HISREAL z 
                                where  EMPRESA_RECNO = 5 AND DTRECEB <= @df
                                ) x where linha = 1 and x.SALDO_EMPRESA > 0)
                order by CODIGO,  Linha asc 
                """
        # Carregar os dados do SQL Server
        logging.info("Carregando dados do SQL Server...")
        df = pd.read_sql(query, conn)

        # Converter a coluna DTRECEB para datetime (assumindo o formato AAAAMMDD)
        df['DTRECEB'] = pd.to_datetime(df['DTRECEB'].astype(str), format='%Y%m%d')

        # Ordena os registros por CODIGO e depois por DTRECEB
        df = df.sort_values(by=['CODIGO', 'DTRECEB'])

        # Lista para armazenar a análise por item
        analise_itens = []
        logging.info("Iniciando análise por item...")
        # Agrupa os dados por item (CODIGO)
        for codigo, grupo in df.groupby('CODIGO'):
            grupo = grupo.sort_values(by='Linha',ascending=True)

            # Saldo inicial e final do item
            saldo_inicial = grupo.iloc[0]['SALDO_EMPRESA']
            primomov = grupo.iloc[0]['QtdeMovimento']
            saldo_final = grupo.iloc[-1]['SALDO_EMPRESA']
            if str(grupo.iloc[0]['Tipo']).lower() == 'consumo':
                saldo_inicial = saldo_inicial + primomov
            else:
                saldo_inicial = saldo_inicial - primomov
            
            valcustomedio = grupo.iloc[-1]['CUSTO_MEDIO'] if not grupo.empty else 0

            # Estoque médio
            estoque_medio = grupo.iloc[0]['EstoqueMedio']
          
            # Filtra os movimentos de compra e consumo (ajusta para caixa baixa)
            compras = grupo[(grupo['Tipo'].str.lower() == 'compra') & (grupo['QtdeMovimento'] > 0)]
            consumos = grupo[grupo['Tipo'].str.lower() == 'consumo']
            
            # Volume comprado: soma da quantidade movimentada dos registros de compra
            volume_comprado = compras['QtdeMovimento'].sum() if not compras.empty else 0
            # Número de vezes que o item foi comprado e consumido
            num_compras = contar_notas(codigo) #len(compras)
            num_consumo = len(consumos)
            
            _xconsumo = grupo[grupo['Tipo'].str.lower() == 'consumo'] 

            _ultimoconsumo = grupo[(grupo['Tipo'].str.lower() == 'consumo') & (grupo['Tipo'].str.lower() != 'inv.local') & (grupo['Tipo'].str.lower() != 'ajuste de estoque')]
            _ultimoConsumo = _ultimoconsumo.iloc[-1]['DTRECEB'].strftime('%Y%m%d') if not _xconsumo.empty else '20231231'

            data_atual = datetime.now().strftime("%d/%m/%Y")
            intervalo = calcular_intervalo(_ultimoConsumo, data_atual)

            _consumido = _xconsumo['QtdeMovimento'].sum() if not _xconsumo.empty else 0
            _xdevolucao = grupo[grupo['Tipo'].str.lower() == 'devolvido']
            _devolucaoAlmox = _xdevolucao['QtdeMovimento'].sum() if not _xdevolucao.empty else 0
            _consumido = _consumido - _devolucaoAlmox

            # Giro de estoque (evita divisão por zero)
            if _consumido is None or pd.isna(_consumido) or _consumido <= 0:
                giro = 0
            else:
                if estoque_medio is None or pd.isna(estoque_medio) or estoque_medio <= 0:
                    giro = 0
                else:
                    giro = float(_consumido) / float(estoque_medio)
            
            if _consumido <= 0:
                _consumido = 0

            if giro == 0:
                tipogiro = 'Estoque Parado'
            elif giro < 0.3:
                tipogiro = 'Giro < 0.3 Periodo'
            elif giro < 0.5:
                tipogiro = 'Giro < 0.5 Periodo'
            elif giro < 1:
                tipogiro = 'Giro < 1 Periodo'
            elif giro < 2:
                tipogiro = 'Giro < 2 Periodo'
            elif giro < 3:
                tipogiro = 'Giro < 3 Periodo'
            elif giro < 4:
                tipogiro = 'Giro < 4 Periodo'
            elif giro < 5:
                tipogiro = 'Giro < 5 Periodo'
            else:
                tipogiro = 'Giro > 5 Periodo'

            if (_consumido or 0) > 0 and (saldo_final or 0) == 0 and (estoque_medio or 0) == 0:
                tipogiro = 'Sem estoque'
            # Calcula o intervalo médio entre as datas dos movimentos de compra e consumo
            avg_interval_compra = compute_avg_interval(compras['DTRECEB'].tolist()) if not compras.empty else None
            avg_interval_consumo = compute_avg_interval(consumos['DTRECEB'].tolist()) if not consumos.empty else None
            
            analise_itens.append({
                'CODIGO': codigo,
                'DESCRICAO': f"{codigo}-{grupo.iloc[0]['DESCRICAO']}", 
                'GrupoItem': grupo.iloc[0]['GRUPOESTOQUE'], 
                'Categoria': grupo.iloc[0]['Categoria'],
                'Familia': grupo.iloc[0]['FAMILIA'],
                'Saldo_Inicial': saldo_inicial,
                'Saldo_Final': saldo_final,
                'Custo_Medio': 0 if pd.isna(valcustomedio) else int(valcustomedio),
                'Custo_Total': saldo_final * (0 if pd.isna(valcustomedio) else float(valcustomedio)),
                'Estoque_Medio': 0 if pd.isna(estoque_medio) else int(estoque_medio),
                'ClassificacaoGiro': tipogiro,
                'Giro_Estoque': 0 if pd.isna(giro) else  float(giro),  
                'Registros': len(grupo),
                'Volume_Comprado': 0 if pd.isna(volume_comprado) else int(volume_comprado), 
                'Volume_Consumido': 0 if pd.isna(_consumido) else int(_consumido), 
                'Devolvido_Almox': 0 if pd.isna(_devolucaoAlmox) else int(_devolucaoAlmox),
                'Numero_Compras': int(num_compras),
                'Numero_Consumo': int(num_consumo),
                'Dias_Sem_Giracao': 0 if pd.isna(intervalo) else int(intervalo),  
                'UltimoConsumo': int(_ultimoConsumo), 
                'Intervalo_Medio_Compra': 0 if pd.isna(avg_interval_compra) else int(avg_interval_compra),
                'Intervalo_Medio_Consumo': 0 if pd.isna(avg_interval_consumo) else int(avg_interval_consumo),
            })

        # Cria um DataFrame com as análises por item
        df_analise = pd.DataFrame(analise_itens)
        # Ordena os itens pelo giro de estoque (do menor para o maior)
        df_analise = df_analise.sort_values(by='Giro_Estoque')
        logging.info("Limpando tabela CST_AnaliseGiroEstoque e inserindo novos dados...")
        # Inserir dados no SQL Server
        cursor = conn.cursor()
        delete_sql = "DELETE FROM CST_AnaliseGiroEstoque"
        cursor.execute(delete_sql)
        insert_sql = """
        INSERT INTO CST_AnaliseGiroEstoque (EmpresaRecno,
            CODIGO, DESCRICAO, GrupoItem, Categoria, Familia, 
            Saldo_Inicial, Saldo_Final, Custo_Medio, Custo_Total, Estoque_Medio, 
            ClassificacaoGiro, Giro_Estoque, Registros, Volume_Comprado, Volume_Consumido, 
            Devolvido_Almox, Numero_Compras, Numero_Consumo, Dias_Sem_Giracao, UltimoConsumo,
            Intervalo_Medio_Compra, Intervalo_Medio_Consumo
        ) VALUES (5,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        """
        logging.info("Preparando a tabela CST_AnaliseGiroEstoque para inserir dados...")
        for _, row in df_analise.iterrows():
            cursor.execute(insert_sql, 
                row['CODIGO'], row['DESCRICAO'], row['GrupoItem'], row['Categoria'], row['Familia'],
                row['Saldo_Inicial'], row['Saldo_Final'], row['Custo_Medio'], row['Custo_Total'], row['Estoque_Medio'],
                row['ClassificacaoGiro'], row['Giro_Estoque'], row['Registros'], row['Volume_Comprado'], row['Volume_Consumido'],
                row['Devolvido_Almox'], row['Numero_Compras'], row['Numero_Consumo'], row['Dias_Sem_Giracao'], row['UltimoConsumo'],
                row['Intervalo_Medio_Compra'], row['Intervalo_Medio_Consumo']
            )

        conn.commit()
        cursor.close()
        conn.close()
        logging.info("Dados inseridos com sucesso na tabela CST_AnaliseGiroEstoque")
        logging.info("Job concluído com sucesso")
    except Exception as e:
        logging.exception(f"Falha no job: {e}")
        if conn is not None:
            try:
                conn.close()
            except:
                pass
        raise
