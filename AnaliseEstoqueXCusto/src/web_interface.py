#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Interface Web para o Analisador de Custo Médio
Permite seleção de período e tipos de relatório
"""

from flask import Flask, render_template, request, jsonify, redirect, url_for, session
from datetime import datetime, date
import os
import threading
import pandas as pd
import requests
from functools import wraps
from custo_medio_analyzer import CustoMedioAnalyzer

# Definir caminho correto para templates
template_dir = os.path.join(os.path.dirname(os.path.dirname(__file__)), 'templates')
app = Flask(__name__, template_folder=template_dir)

# Configuração para sessões
app.secret_key = os.getenv('SECRET_KEY', 'sua-chave-secreta-aqui')

# URL da API de autenticação
AUTH_API_BASE = 'http://192.168.1.28:9091/api/auth'

# Status global para acompanhar execução
status_execucao = {
    'executando': False,
    'progresso': '',
    'resultado': '',
    'erro': ''
}

def verificar_jwt_token(token):
    """Verifica se o token JWT é válido fazendo request para API de auth"""
    try:
        headers = {'Authorization': f'Bearer {token}'}
        response = requests.get(f'{AUTH_API_BASE}/me', headers=headers, timeout=10)
        if response.status_code == 200:
            return response.json()
        return None
    except Exception as e:
        print(f"Erro ao verificar token: {e}")
        return None

def login_required(f):
    """Decorador para rotas que requerem autenticação"""
    @wraps(f)
    def decorated_function(*args, **kwargs):
        token = session.get('jwt_token')
        if not token:
            return redirect(url_for('login'))

        # Verificar se o token ainda é válido
        user_data = verificar_jwt_token(token)
        if not user_data:
            session.clear()
            return redirect(url_for('login'))

        # Disponibilizar dados do usuário para a rota
        request.user = user_data
        return f(*args, **kwargs)
    return decorated_function

@app.route('/login', methods=['GET', 'POST'])
def login():
    """Página e processo de login"""
    if request.method == 'GET':
        # Se já está logado, redirecionar
        if session.get('jwt_token') and verificar_jwt_token(session['jwt_token']):
            return redirect(url_for('index'))
        return render_template('login.html')

    # POST - processo de login
    try:
        username = request.form.get('username')
        password = request.form.get('password')

        if not username or not password:
            return jsonify({'erro': 'Username e password são obrigatórios'}), 400

        # Fazer request para API de login
        login_data = {
            'username': username,
            'password': password
        }

        response = requests.post(f'{AUTH_API_BASE}/login', json=login_data, timeout=10)

        if response.status_code == 200:
            result = response.json()
            token = result.get('token') or result.get('access_token') or result.get('jwt')

            if token:
                session['jwt_token'] = token
                return jsonify({'sucesso': True})
            else:
                return jsonify({'erro': 'Token não encontrado na resposta'}), 500
        else:
            error_msg = 'Credenciais inválidas'
            try:
                error_data = response.json()
                error_msg = error_data.get('message', error_msg)
            except:
                pass
            return jsonify({'erro': error_msg}), 401

    except requests.exceptions.RequestException as e:
        return jsonify({'erro': f'Erro de conexão com servidor de autenticação: {str(e)}'}), 500
    except Exception as e:
        return jsonify({'erro': f'Erro interno: {str(e)}'}), 500

@app.route('/logout')
def logout():
    """Logout do usuário"""
    session.clear()
    return redirect(url_for('login'))

@app.route('/')
@login_required
def index():
    """Página principal com formulário"""
    hoje = date.today()
    return render_template('index.html',
                         mes_atual=hoje.month,
                         ano_atual=hoje.year,
                         user=request.user)

@app.route('/executar', methods=['POST'])
@login_required
def executar():
    """Executa a análise com os parâmetros selecionados"""
    global status_execucao

    if status_execucao['executando']:
        return jsonify({'erro': 'Já existe uma análise em execução'}), 400

    try:
        mes = int(request.form.get('mes'))
        ano = int(request.form.get('ano'))
        empresa = request.form.get('empresa') # Captura o valor do filtro 'empresa'
        incluir_custo = request.form.get('incluir_custo') == 'on'
        incluir_saldo = request.form.get('incluir_saldo') == 'on'

        if not incluir_custo and not incluir_saldo:
            return jsonify({'erro': 'Selecione pelo menos um tipo de relatório'}), 400

        # Executar em thread separada
        thread = threading.Thread(target=executar_analise,
                                args=(mes, ano, empresa, incluir_custo, incluir_saldo)) # Passa 'empresa'
        thread.start()

        return jsonify({'sucesso': True})

    except Exception as e:
        return jsonify({'erro': str(e)}), 400

@app.route('/status')
@login_required
def obter_status():
    """Retorna status atual da execução"""
    return jsonify(status_execucao)

@app.route('/analisar')
@login_required
def analisar_resultados():
    """Página para análise de resultados"""
    hoje = date.today()
    return render_template('analisar.html',
                         mes_atual=hoje.month,
                         ano_atual=hoje.year,
                         user=request.user)

@app.route('/buscar_resultados', methods=['POST'])
@login_required
def buscar_resultados():
    """Busca resultados baseado nos filtros aplicados"""
    try:
        mes = int(request.form.get('mes'))
        ano = int(request.form.get('ano'))
        empresa = request.form.get('empresa') # Captura o valor do filtro 'empresa'
        filtro_codigo = request.form.get('filtro_codigo', '').strip()
        filtro_descricao = request.form.get('filtro_descricao', '').strip()
        incluir_custo = request.form.get('incluir_custo') == 'on'
        incluir_saldo = request.form.get('incluir_saldo') == 'on'
        apenas_diferencas = request.form.get('apenas_diferencas') == 'on'

        if not incluir_custo and not incluir_saldo:
            return jsonify({'erro': 'Selecione pelo menos um tipo de análise'}), 400

        analyzer = CustoMedioAnalyzer()
        analyzer.mes_filtro = mes
        analyzer.ano_filtro = ano
        analyzer.empresa_filtro = empresa # Define o filtro de empresa

        if not analyzer.conectar_sql_server():
            return jsonify({'erro': 'Falha ao conectar ao SQL Server'}), 500

        df_original = analyzer.executar_query()
        if df_original is None or len(df_original) == 0:
            return jsonify({'erro': 'Nenhum dado encontrado para o período selecionado'}), 400

        df_resultados = analyzer.calcular_custo_medio(df_original)

        # Aplicar filtros
        if filtro_codigo:
            df_resultados = df_resultados[df_resultados['Codigo'].str.contains(filtro_codigo, case=False, na=False)]

        if filtro_descricao:
            df_resultados = df_resultados[df_resultados['Item'].str.contains(filtro_descricao, case=False, na=False)]

        # Filtrar apenas diferenças se solicitado
        if apenas_diferencas:
            condicoes = []
            if incluir_custo:
                condicoes.append(df_resultados['TemDiferencaCusto'] == True)
            if incluir_saldo:
                condicoes.append(df_resultados['TemDiferencaSaldo'] == True)

            if condicoes:
                df_filtrado = df_resultados[condicoes[0]]
                for condicao in condicoes[1:]:
                    df_filtrado = pd.concat([df_filtrado, df_resultados[condicao]]).drop_duplicates()
                df_resultados = df_filtrado

        # Preparar dados para exibição
        colunas_exibicao = [
            'Ordem', 'Item', 'Codigo', 'TipoMovimento', 'QtdeMovimento',
            'DataMovimento', 'MovimentoEfetivo', 'SaldoOriginal', 'SaldoCalculado'
        ]

        if incluir_custo:
            colunas_exibicao.extend(['CustoMedioOriginal', 'CustoMedioCalculado', 'DiferencaCusto', 'TemDiferencaCusto'])

        if incluir_saldo:
            colunas_exibicao.extend(['DiferencaSaldo', 'TemDiferencaSaldo'])

        # Limitar a 1000 registros para performance
        df_resultado_final = df_resultados[colunas_exibicao].head(1000)

        # Converter para formato JSON
        dados = df_resultado_final.to_dict('records')

        # Formatar dados para melhor exibição
        for item in dados:
            if 'DataMovimento' in item and item['DataMovimento'] is not None:
                # Verificar se é datetime ou string
                if hasattr(item['DataMovimento'], 'strftime'):
                    item['DataMovimento'] = item['DataMovimento'].strftime('%d/%m/%Y')

            # Formatar valores numéricos
            for col in ['SaldoOriginal', 'SaldoCalculado', 'DiferencaSaldo']:
                if col in item and item[col] is not None:
                    try:
                        item[col] = f"{float(item[col]):.2f}"
                    except (ValueError, TypeError):
                        pass

            for col in ['CustoMedioOriginal', 'CustoMedioCalculado', 'DiferencaCusto']:
                if col in item and item[col] is not None:
                    try:
                        item[col] = f"{float(item[col]):.4f}"
                    except (ValueError, TypeError):
                        pass

        analyzer.conn.close()

        return jsonify({
            'sucesso': True,
            'dados': dados,
            'total_registros': len(df_resultados),
            'registros_exibidos': len(dados)
        })

    except Exception as e:
        return jsonify({'erro': str(e)}), 500

def executar_analise(mes, ano, empresa, incluir_custo, incluir_saldo):
    """Executa a análise em background"""
    global status_execucao

    status_execucao.update({
        'executando': True,
        'progresso': 'Iniciando análise...',
        'resultado': '',
        'erro': ''
    })

    try:
        analyzer = CustoMedioAnalyzer()
        analyzer.mes_filtro = mes
        analyzer.ano_filtro = ano
        analyzer.empresa_filtro = empresa # Define o filtro de empresa
        analyzer.incluir_custo = incluir_custo
        analyzer.incluir_saldo = incluir_saldo

        # Executar com callback de progresso
        resultado = analyzer.executar_com_interface()

        status_execucao.update({
            'executando': False,
            'progresso': 'Concluído',
            'resultado': resultado
        })

    except Exception as e:
        status_execucao.update({
            'executando': False,
            'progresso': 'Erro',
            'erro': str(e)
        })

if __name__ == '__main__':
    # Criar pasta templates se não existir
    if not os.path.exists('templates'):
        os.makedirs('templates')

    app.run(host='0.0.0.0', port=5000, debug=True)