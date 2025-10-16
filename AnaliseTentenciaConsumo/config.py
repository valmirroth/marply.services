# ==========================================================
# ARQUIVO 1: config.py
# Configurações do sistema
# ==========================================================

SQL_SERVER = "192.168.1.28"
SQL_DATABASE = "Marini_PRD"
SQL_USER = "Octopus"
SQL_PASSWORD = "A45182008069199"

EMPRESAS = {
    'matriz': '1',
    'filial': '5',
    'consolidado': '1,5'
}

SECRET_KEY = 'sua-chave-secreta-aqui-mude-em-producao'
PORT = 8096