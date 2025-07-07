#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script de teste para validar configurações antes de executar o analisador
"""

import os
from dotenv import load_dotenv
import sys

load_dotenv()

def testar_variaveis_ambiente():
    """Verifica se todas as variáveis estão configuradas"""
    print("\n" + "="*60)
    print("TESTE DE CONFIGURAÇÃO")
    print("="*60)
    
    variaveis = {
        'SQL_SERVER': os.getenv('SQL_SERVER'),
        'SQL_DATABASE': os.getenv('SQL_DATABASE'),
        'SQL_USERNAME': os.getenv('SQL_USERNAME'),
        'SQL_PASSWORD': os.getenv('SQL_PASSWORD'),
        'SMTP_SERVER': os.getenv('SMTP_SERVER'),
        'SMTP_USERNAME': os.getenv('SMTP_USERNAME'),
        'SMTP_PASSWORD': os.getenv('SMTP_PASSWORD'),
        'EMAIL_FROM': os.getenv('EMAIL_FROM'),
        'EMAIL_TO': os.getenv('EMAIL_TO'),
    }
    
    print("\n1. Verificando variáveis de ambiente...")
    todas_ok = True
    for nome, valor in variaveis.items():
        if valor and valor != f'seu_{nome.lower()}' and valor != 'seu-servidor.database.windows.net':
            print(f"   ✓ {nome}: {'*' * 10} (configurado)")
        else:
            print(f"   ✗ {nome}: NÃO CONFIGURADO")
            todas_ok = False
    
    return todas_ok

def testar_bibliotecas():
    """Verifica se as bibliotecas estão instaladas"""
    print("\n2. Verificando bibliotecas Python...")
    bibliotecas = {
        'pyodbc': 'pyodbc',
        'pandas': 'pandas',
        'dotenv': 'python-dotenv'
    }
    
    todas_ok = True
    for modulo, pacote in bibliotecas.items():
        try:
            __import__(modulo)
            print(f"   ✓ {pacote}: instalado")
        except ImportError:
            print(f"   ✗ {pacote}: NÃO INSTALADO - execute: pip install {pacote}")
            todas_ok = False
    
    return todas_ok

def testar_conexao_sql():
    """Testa conexão com SQL Server"""
    print("\n3. Testando conexão SQL Server...")
    try:
        import pyodbc
        
        sql_server = os.getenv('SQL_SERVER')
        sql_database = os.getenv('SQL_DATABASE')
        sql_username = os.getenv('SQL_USERNAME')
        sql_password = os.getenv('SQL_PASSWORD')
        sql_port = os.getenv('SQL_PORT', '1433')
        
        connection_string = (
            f"DRIVER={{ODBC Driver 17 for SQL Server}};"
            f"SERVER={sql_server},{sql_port};"
            f"DATABASE={sql_database};"
            f"UID={sql_username};"
            f"PWD={sql_password};"
            f"Connection Timeout=10"
        )
        
        conn = pyodbc.connect(connection_string)
        cursor = conn.cursor()
        cursor.execute("SELECT @@VERSION")
        versao = cursor.fetchone()[0]
        print(f"   ✓ Conexão estabelecida com sucesso!")
        print(f"   ✓ Servidor: {sql_server}")
        print(f"   ✓ Database: {sql_database}")
        conn.close()
        return True
    except Exception as e:
        print(f"   ✗ Erro ao conectar: {e}")
        return False

def testar_smtp():
    """Testa conexão SMTP"""
    print("\n4. Testando conexão SMTP...")
    try:
        import smtplib
        
        smtp_server = os.getenv('SMTP_SERVER')
        smtp_port = int(os.getenv('SMTP_PORT', '587'))
        smtp_username = os.getenv('SMTP_USERNAME')
        smtp_password = os.getenv('SMTP_PASSWORD')
        
        server = smtplib.SMTP(smtp_server, smtp_port, timeout=10)
        server.starttls()
        server.login(smtp_username, smtp_password)
        server.quit()
        
        print(f"   ✓ Conexão SMTP estabelecida com sucesso!")
        print(f"   ✓ Servidor: {smtp_server}:{smtp_port}")
        return True
    except Exception as e:
        print(f"   ✗ Erro ao conectar SMTP: {e}")
        return False

if __name__ == "__main__":
    print("\n🔍 Iniciando testes de configuração...\n")
    
    resultados = []
    
    resultados.append(("Variáveis de Ambiente", testar_variaveis_ambiente()))
    resultados.append(("Bibliotecas Python", testar_bibliotecas()))
    resultados.append(("Conexão SQL Server", testar_conexao_sql()))
    resultados.append(("Conexão SMTP", testar_smtp()))
    
    print("\n" + "="*60)
    print("RESUMO DOS TESTES")
    print("="*60)
    
    for nome, resultado in resultados:
        status = "✓ PASSOU" if resultado else "✗ FALHOU"
        print(f"{nome}: {status}")
    
    todos_ok = all(r[1] for r in resultados)
    
    print("\n" + "="*60)
    if todos_ok:
        print("✓ TODOS OS TESTES PASSARAM!")
        print("Você pode executar o script principal agora:")
        print("  python custo_medio_analyzer.py")
    else:
        print("✗ ALGUNS TESTES FALHARAM")
        print("Corrija os problemas acima antes de executar o script principal.")
        sys.exit(1)
    print("="*60 + "\n")
