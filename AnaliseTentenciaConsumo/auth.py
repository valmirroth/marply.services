# ==========================================================
# ARQUIVO 3: auth.py
# Autenticação e gerenciamento de usuários
# ==========================================================

import os
from functools import wraps
from flask import session, redirect, url_for
import sys

def carregar_usuarios():
    """Carrega usuários do arquivo .inv e grava o caminho em um .txt"""
#sys.executable
# EM AMBIENTE DEBUG USAR __file__
    base_dir = os.path.dirname(os.path.abspath( sys.executable  ))
    arquivo_usuarios = os.path.join(base_dir, '.inv')
    arquivo_log = os.path.join(base_dir, 'caminho_inv.txt')

    usuarios = {}

    # Grava o caminho do .inv no txt
    def gravar_caminho():
        try:
            with open(arquivo_log, 'w', encoding='utf-8') as f:
                f.write("Caminho do arquivo .inv:\n")
                f.write(arquivo_usuarios)
        except Exception as e:
            print(f"Erro ao gravar caminho_inv.txt: {e}")

    if os.path.exists(arquivo_usuarios):
        print(f"Existe o arquivo .inv: {arquivo_usuarios}")
        gravar_caminho()

        try:
            with open(arquivo_usuarios, 'r', encoding='utf-8') as f:
                for linha in f:
                    linha = linha.strip()
                    if linha and not linha.startswith('#'):
                        partes = linha.split(':', 1)
                        if len(partes) == 2:
                            usuario, senha = partes
                            usuarios[usuario.strip()] = senha.strip()
        except Exception as e:
            print(f"Erro ao carregar arquivo .inv: {e}")
    else:
        try:
            with open(arquivo_usuarios, 'w', encoding='utf-8') as f:
                f.write("# Arquivo de usuários - Formato: usuario:senha\n")
                f.write("# Adicione um usuário por linha\n")
                f.write("admin:admin123\n")
                f.write("usuario:senha123\n")

            print(f"Arquivo .inv criado em: {arquivo_usuarios}")
            usuarios = {'admin': 'admin123', 'usuario': 'senha123'}

            gravar_caminho()

        except Exception as e:
            print(f"Erro ao criar arquivo .inv: {e}")

    return usuarios

def login_required(f):
    """Decorador para proteger rotas"""
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'logged_in' not in session:
            return redirect(url_for('login'))
        return f(*args, **kwargs)
    return decorated_function
