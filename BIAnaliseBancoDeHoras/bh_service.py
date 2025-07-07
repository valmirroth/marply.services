# comando para empacotar com PyInstaller:
# Certifique-se de que o ambiente virtual está ativado e que todas as dependências estão instaladas.
# Use o comando abaixo para criar um executável único sem console, coletando todos os arquivos necessários do pandas.
# Ajuste o caminho do log conforme necessário.
# Se necessário, adicione outras dependências ao comando de coleta.
# Exemplo de comando:
# python -m PyInstaller --onefile --noconsole --collect-all pandas bh_service.py

import time
import logging
import os
from datetime import datetime, timedelta
import pandas as pd   # garante coleta no PyInstaller (--collect-all pandas)
import pyodbc         # garante coleta do driver
import subprocess
from pathlib import Path
from bh_process_to_sql_full import JobBH  # importa a função principal do script

# Configurações via variáveis de ambiente
LOG_PATH = os.getenv("BH_LOG_PATH", r"C:\MarplyServices\BI - Calculo saldo banco de horas\bh_service.log")
RUN_HOUR = int(os.getenv("BH_RUN_HOUR", "3"))   # 3 = 03h
RUN_MIN  = int(os.getenv("BH_RUN_MIN", "0"))    # 0 = 00min
PYTHON_EXE = os.getenv("BH_PYTHON_EXE", None)   # se None, usa o Python atual
#SCRIPT_NAME = os.getenv("BH_SCRIPT_NAME", "bh_process_to_sql_full.py")
RUN_ON_START = os.getenv("BH_RUN_ON_START", "1")  # "1" executa imediatamente ao iniciar o serviço

# Logging
os.makedirs(os.path.dirname(LOG_PATH), exist_ok=True)
logging.basicConfig(filename=LOG_PATH, level=logging.INFO,
                    format="%(asctime)s %(levelname)s %(message)s")

def next_run_at(hour: int, minute: int) -> datetime:
    now = datetime.now()
    run = now.replace(hour=hour, minute=minute, second=0, microsecond=0)
    if run <= now:
        run = run + timedelta(days=1)
    return run

def run_job() -> int:
    """Executa o script principal como subprocesso, capturando stdout/stderr nos logs."""
    exe = PYTHON_EXE or os.sys.executable
    #script_path = Path(__file__).with_name(SCRIPT_NAME)
    JobBH()  # Chama a função diretamente
    logging.info("Iniciando job: %s", " ".join(cmd))
    try:
        proc = subprocess.run(cmd, capture_output=True, text=True)
    except Exception as e:
        logging.exception("Falha ao iniciar o job: %s", e)
        return -1
    if proc.stdout:
        logging.info("Job stdout:\n%s", proc.stdout)
    if proc.stderr:
        logging.error("Job stderr:\n%s", proc.stderr)
    logging.info("Job finalizado com código %s", proc.returncode)
    return proc.returncode

def main():
    try:
        if RUN_ON_START == "1":
            logging.info("Serviço iniciado. Executando job imediatamente (RUN_ON_START=1).")
            run_job()
    except Exception as e:
        logging.exception("Erro ao executar job inicial: %s", e)
    
    logging.info("Serviço iniciado. Executando job imediatamente (RUN_ON_START=1).")
    run_job()
    logging.info("Serviço em loop (NSSM). Agendado diariamente às %02d:%02d.", RUN_HOUR, RUN_MIN)
    while True:
        nr = next_run_at(RUN_HOUR, RUN_MIN)
        wait_sec = (nr - datetime.now()).total_seconds()
        if wait_sec > 0:
            logging.info("Aguardando até %s (%.0f segundos) para executar o job.", nr.isoformat(), wait_sec)
            slept = 0.0
            while slept < wait_sec:
                chunk = min(300.0, max(0.0, wait_sec - slept))  # 5 min
                time.sleep(chunk)
                slept += chunk
        try:
            run_job()
        except Exception as e:
            logging.exception("Erro ao executar job agendado: %s", e)

if __name__ == "__main__":
    main()
