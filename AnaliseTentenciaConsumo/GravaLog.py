import pandas as pd
from database import QUERY_TEMPLATE, get_conn

def inserir_log_acesso(mes_atual, empresa_filter, num_meses=10, username=None):
    """
    Registra log de acesso aos dados no banco
    
    Parâmetros:
    -----------
    mes_atual : str
        Indica se considera mês atual ('sim' ou 'nao')
    empresa_filter : str
        Filtro de empresa aplicado ('matriz', 'filial', 'consolidado')
    num_meses : int
        Número de meses considerados na análise (padrão: 10)
    username : str
        Nome do usuário que executou a consulta (opcional)
    """
    from datetime import datetime
    
    query_log = """
    INSERT INTO CST_LOG_ACESSO_ANALISE_CUSTOS (
        data_hora_acesso,
        usuario,
        mes_atual,
        empresa_filter,
        num_meses
    ) VALUES (?, ?, ?, ?, ?)
    """
    
    try:
        with get_conn() as conn:
            cursor = conn.cursor()
            cursor.execute(query_log, (
                datetime.now(),
                username or 'desconhecido',
                mes_atual,
                empresa_filter,
                num_meses
            ))
            conn.commit()
            print(f"Log registrado com sucesso para usuário: {username}")
            
    except Exception as e:
        print(f"Erro ao registrar log: {str(e)}")