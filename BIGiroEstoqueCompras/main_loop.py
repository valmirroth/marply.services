# comando para empacotar com PyInstaller:
# Certifique-se de que o ambiente virtual está ativado e que todas as dependências estão instaladas.
# Use o comando abaixo para criar um executável único sem console, coletando todos os arquivos necessários do pandas.
# Ajuste o caminho do log conforme necessário.
# Se necessário, adicione outras dependências ao comando de coleta.
# Exemplo de comando:
# python -m PyInstaller --onefile --noconsole --collect-all pandas main_loop.py
import time
import logging
import os
from datetime import datetime, timedelta
from job_giro import run_job
from job_filial import run_job as run_job_filial
import pandas as pd
import pyodbc

LOG_PATH = os.getenv("GIRO_LOG_PATH", r"C:\MarplyServices\BI - AnaliseGiroEstoque\logs\giro_nssm.log")
os.makedirs(os.path.dirname(LOG_PATH), exist_ok=True)
logging.basicConfig(filename=LOG_PATH, level=logging.INFO,
                    format="%(asctime)s %(levelname)s %(message)s")

TARGET_HOUR = int(os.getenv("GIRO_RUN_HOUR", "5"))   # 4 = 04h
TARGET_MIN  = int(os.getenv("GIRO_RUN_MIN", "0"))    # 0 = 00min

def next_run_at(hour: int, minute: int) -> datetime:
    now = datetime.now()
    run = now.replace(hour=hour, minute=minute, second=0, microsecond=0)
    if run <= now:
        run = run + timedelta(days=1)
    return run

def main():
    try:
            logging.info("Serviço em loop (NSSM) iniciado, executando o job imediatamente no inicio do JOB.")
            logging.info("Executando job de análise de giro de estoque para a filial.")
            run_job_filial()
            logging.info("Finalizado job de análise de giro de estoque para a filial.")
            # Executa o job de análise de giro de estoque
            logging.info("Executando job de análise de giro de estoque para a Matriz.")
            run_job()
            logging.info("Finalizado job de análise de giro de estoque para a Matriz.")
    except Exception as e:
            logging.exception("Erro ao executar run_job: %s", e)

    logging.info("Serviço em loop (NSSM) iniciado. Agendado diariamente às %02d:%02d.", TARGET_HOUR, TARGET_MIN)
    while True:
        nr = next_run_at(TARGET_HOUR, TARGET_MIN)
        wait_sec = (nr - datetime.now()).total_seconds()
        if wait_sec > 0:
            logging.info("Aguardando até %s (%.0f segundos) para executar o job.", nr.isoformat(), wait_sec)
            # Dorme em pedaços para permitir logs periódicos longos
            slept = 0
            while slept < wait_sec:
                chunk = min(300.0, wait_sec - slept)  # 5 min chunks
                time.sleep(chunk)
                slept += chunk
        try:
            logging.info("Executando job de análise de giro de estoque para a filial.")
            run_job_filial()
            logging.info("Finalizado job de análise de giro de estoque para a filial.")
            # Executa o job de análise de giro de estoque
            logging.info("Executando job de análise de giro de estoque para a Matriz.")
            run_job()
            logging.info("Finalizado job de análise de giro de estoque para a Matriz.")
        except Exception as e:
            logging.exception("Erro ao executar run_job: %s", e)
        # Após rodar, agenda para o próximo dia e volta ao início do loop

if __name__ == "__main__":
    main()
