#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Analisador de Custo Médio - SQL Server
Calcula custo médio de itens e compara com valores originais.
Envia email com diferenças encontradas.
"""

import pyodbc
import pandas as pd
import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from datetime import datetime
from decimal import Decimal
import os
from dotenv import load_dotenv
import sys

load_dotenv()


class CustoMedioAnalyzer:

    def __init__(self):
        self.sql_server = os.getenv('SQL_SERVER')
        self.sql_database = os.getenv('SQL_DATABASE')
        self.sql_username = os.getenv('SQL_USERNAME')
        self.sql_password = os.getenv('SQL_PASSWORD')
        self.sql_port = os.getenv('SQL_PORT', '1433')

        self.smtp_server = os.getenv('SMTP_SERVER')
        self.smtp_port = int(os.getenv('SMTP_PORT', '587'))
        self.smtp_username = os.getenv('SMTP_USERNAME')
        self.smtp_password = os.getenv('SMTP_PASSWORD')
        self.email_from = os.getenv('EMAIL_FROM')
        self.email_to = os.getenv('EMAIL_TO', 'ti01@marply.com.br')

        self.diferenca_threshold = float(
            os.getenv('DIFERENCA_THRESHOLD', '0.05'))

        # Parâmetros de filtro
        self.mes_filtro = None
        self.ano_filtro = None
        self.empresa_filtro = 1  # Default para Matriz
        self.incluir_custo = True
        self.incluir_saldo = True

        self.conn = None

    def conectar_sql_server(self):
        """Conecta ao SQL Server"""
        try:
            connection_string = (f"DRIVER={{ODBC Driver 17 for SQL Server}};"
                                 f"SERVER={self.sql_server},{self.sql_port};"
                                 f"DATABASE={self.sql_database};"
                                 f"UID={self.sql_username};"
                                 f"PWD={self.sql_password}")
            self.conn = pyodbc.connect(connection_string)
            print(
                f"✓ Conectado ao SQL Server: {self.sql_server}/{self.sql_database}"
            )
            return True
        except Exception as e:
            print(f"✗ Erro ao conectar ao SQL Server: {e}")
            return False

    def executar_query(self):
        """Executa a query principal e retorna os dados"""
        # Determinar filtro de data
        if self.mes_filtro and self.ano_filtro:
            filtro_data = f"MONTH(DTRECEB) = {self.mes_filtro} AND YEAR(DTRECEB) = {self.ano_filtro}"
        else:
            filtro_data = "EOMONTH(DTRECEB) = EOMONTH(GETDATE())"

        query = f"""
        SELECT 
            ROW_NUMBER() OVER(PARTITION BY hisreal.codigo, hisreal.empresa_recno 
                             ORDER BY hisreal.codigo, hisreal.empresa_recno) AS Ordem,
            HISREAL.CODIGO + ' - ' + HISREAL.DESCRICAO AS Item, 
            HISREAL.CODIGO AS Codigo,
            CASE WHEN FORMA = 'E' THEN 'Entrada' ELSE 'Saida' END AS TipoMovimento,
            QTRECEB AS QtdeMovimento,
            saldo_empresa AS SaldoEmpresa, 
            HISREAL.MOV_EFETIVA AS MovimentoEfetivo,
            DTRECEB AS DataMovimento,
            EMPRESA_RECNO AS Empresa,
            HISREAL.CUSTO_UNITARIO AS CustoUnitario,
            HISREAL.CUSTO_MEDIO AS CustoMedio,
            HISREAL.CUSTO_TOTAL AS CustoTotal
        FROM HISREAL 
        INNER JOIN ESTOQUE ON ESTOQUE.CODIGO = HISREAL.CODIGO 
            AND ESTOQUE.FAMILIA IN (2,3,16,4,166,255,5,15)
            and HISREAL.FORMA IN ('E','S')
            and HISREAL.EMPRESA_RECNO = {self.empresa_filtro}
        WHERE {filtro_data}
        ORDER BY hisreal.codigo, DTRECEB ASC, HISREAL.R_E_C_N_O_ ASC
        """

        try:
            df = pd.read_sql(query, self.conn)
            print(
                f"✓ Query executada com sucesso: {len(df)} registros encontrados"
            )
            return df
        except Exception as e:
            print(f"✗ Erro ao executar query: {e}")
            return None

    def calcular_custo_medio(self, df):
        """Calcula o custo médio para cada item começando pela ordem 2"""
        resultados = []

        grupos = df.groupby(['Codigo', 'Empresa'])

        for (codigo, empresa), grupo_df in grupos:
            grupo_df = grupo_df.sort_values('Ordem').reset_index(drop=True)

            if len(grupo_df) < 1:
                continue

            linha_inicial = grupo_df[grupo_df['Ordem'] == 1]
            if len(linha_inicial) == 0:
                continue

            saldo_anterior = float(linha_inicial.iloc[0]['SaldoEmpresa'])
            custo_medio_anterior = float(linha_inicial.iloc[0]['CustoMedio'])
            item_nome = linha_inicial.iloc[0]['Item']

            linhas_movimento = grupo_df[grupo_df['Ordem'] >= 2]

            for idx, row in linhas_movimento.iterrows():
                tipo_movimento = row['TipoMovimento']
                qtde_movimento = float(row['QtdeMovimento'])
                custo_unitario = float(row['CustoUnitario'])
                custo_medio_original = float(row['CustoMedio'])
                saldo_original = float(row['SaldoEmpresa'])
                movimento_efetivo = row['MovimentoEfetivo']

                # Se MovimentoEfetivo = 'N', não recalcula custo médio
                if movimento_efetivo == 'N':
                    if tipo_movimento == 'Entrada':
                        saldo_calculado = saldo_anterior + qtde_movimento
                    else:
                        saldo_calculado = saldo_anterior - qtde_movimento

                    # Mantém o custo médio anterior (não recalcula)
                    custo_medio_calculado = custo_medio_anterior
                    saldo_anterior = saldo_calculado
                else:
                    # Lógica normal quando MovimentoEfetivo != 'N'
                    if tipo_movimento == 'Entrada':
                        saldo_calculado = saldo_anterior + qtde_movimento

                        if saldo_calculado > 0:
                            custo_medio_calculado = (
                                (saldo_anterior * custo_medio_anterior) +
                                (qtde_movimento *
                                 custo_unitario)) / saldo_calculado
                        else:
                            custo_medio_calculado = custo_medio_anterior

                        saldo_anterior = saldo_calculado
                        custo_medio_anterior = custo_medio_calculado

                    else:
                        saldo_calculado = saldo_anterior - qtde_movimento
                        custo_medio_calculado = custo_medio_anterior

                        saldo_anterior = saldo_calculado

                diferenca_custo = abs(custo_medio_calculado -
                                      custo_medio_original)
                diferenca_saldo = abs(saldo_calculado - saldo_original)

                resultados.append({
                    'Ordem': row['Ordem'],
                    'Item': item_nome,
                    'Codigo': codigo,
                    'Empresa': empresa,
                    'TipoMovimento': tipo_movimento,
                    'QtdeMovimento': qtde_movimento,
                    'DataMovimento': row['DataMovimento'],
                    'MovimentoEfetivo': movimento_efetivo,
                    'SaldoOriginal': saldo_original,
                    'SaldoCalculado': saldo_calculado,
                    'CustoMedioOriginal': custo_medio_original,
                    'CustoMedioCalculado': custo_medio_calculado,
                    'DiferencaCusto': diferenca_custo,
                    'DiferencaSaldo': diferenca_saldo,
                    'TemDiferencaCusto': diferenca_custo
                    > self.diferenca_threshold,
                    'TemDiferencaSaldo':
                    diferenca_saldo > 0.01  # Threshold para saldo
                })

        return pd.DataFrame(resultados)

    def gerar_relatorio_html(self, df_resultados):
        """Gera relatório HTML com as diferenças encontradas"""
        # Filtrar apenas movimentações efetivas (S) para buscar diferenças
        df_efetivas = df_resultados[df_resultados['MovimentoEfetivo'] == 'S']

        # Aplicar filtros baseados na configuração
        df_diferencas_custo = pd.DataFrame()
        df_diferencas_saldo = pd.DataFrame()

        if self.incluir_custo:
            df_diferencas_custo = df_efetivas[df_efetivas['TemDiferencaCusto']
                                              == True]

        if self.incluir_saldo:
            df_diferencas_saldo = df_efetivas[df_efetivas['TemDiferencaSaldo']
                                              == True]

        total_diferencas = len(df_diferencas_custo) + len(df_diferencas_saldo)

        if total_diferencas == 0:
            data_execucao = datetime.now().strftime('%d/%m/%Y %H:%M:%S')
            html = f"""
            <html>
            <head>
                <style>
                    body {{ font-family: Arial, sans-serif; padding: 20px; }}
                    .success {{ color: green; font-weight: bold; }}
                </style>
            </head>
            <body>
                <h2>Análise de Custo Médio e Saldos</h2>
                <p class="success">✓ Nenhuma diferença significativa encontrada!</p>
                <p>Todos os custos médios e saldos estão corretos</p>
                <p><small>Executado em: {data_execucao}</small></p>
            </body>
            </html>
            """
            return html

        # Tabela de diferenças de custo médio
        tabela_custo_html = ""
        if len(df_diferencas_custo) > 0:
            colunas_custo = [
                'Item', 'Ordem', 'TipoMovimento', 'MovimentoEfetivo',
                'DataMovimento', 'CustoMedioOriginal', 'CustoMedioCalculado',
                'DiferencaCusto'
            ]
            tabela_custo_html = df_diferencas_custo[colunas_custo].to_html(
                index=False,
                float_format=lambda x: f'{x:.4f}',
                classes='tabela-diferencas',
                border=0)

        # Tabela de diferenças de saldo
        tabela_saldo_html = ""
        if len(df_diferencas_saldo) > 0:
            colunas_saldo = [
                'Item', 'Ordem', 'TipoMovimento', 'QtdeMovimento',
                'MovimentoEfetivo', 'DataMovimento', 'SaldoOriginal',
                'SaldoCalculado', 'DiferencaSaldo'
            ]
            tabela_saldo_html = df_diferencas_saldo[colunas_saldo].to_html(
                index=False,
                float_format=lambda x: f'{x:.2f}',
                classes='tabela-diferencas',
                border=0)

        html = f"""
        <html>
        <head>
            <style>
                body {{
                    font-family: Arial, sans-serif;
                    padding: 20px;
                    background-color: #f5f5f5;
                }}
                .header {{
                    background-color: #d32f2f;
                    color: white;
                    padding: 15px;
                    border-radius: 5px;
                    margin-bottom: 20px;
                }}
                .summary {{
                    background-color: #fff;
                    padding: 15px;
                    border-left: 4px solid #d32f2f;
                    margin-bottom: 20px;
                }}
                .secao {{
                    background-color: #fff;
                    padding: 15px;
                    margin-bottom: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }}
                .secao h3 {{
                    color: #1976d2;
                    margin-top: 0;
                }}
                .tabela-diferencas {{
                    width: 100%;
                    border-collapse: collapse;
                    background-color: white;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }}
                .tabela-diferencas th {{
                    background-color: #1976d2;
                    color: white;
                    padding: 12px;
                    text-align: left;
                    font-weight: bold;
                }}
                .tabela-diferencas td {{
                    padding: 10px;
                    border-bottom: 1px solid #ddd;
                }}
                .tabela-diferencas tr:hover {{
                    background-color: #f5f5f5;
                }}
                .footer {{
                    margin-top: 20px;
                    font-size: 12px;
                    color: #666;
                }}
            </style>
        </head>
        <body>
            <div class="header">
                <h2>⚠️ Diferenças Encontradas na Análise</h2>
            </div>
            
            <div class="summary">
                <p><strong>Diferenças de Custo Médio:</strong> {len(df_diferencas_custo)}</p>
                <p><strong>Diferenças de Saldo:</strong> {len(df_diferencas_saldo)}</p>
                <p><strong>Threshold de diferença de custo:</strong> {self.diferenca_threshold:.2f}</p>
                <p><strong>Threshold de diferença de saldo:</strong> 0.01</p>
                <p><strong>Data da análise:</strong> {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}</p>
            </div>
            
            {f'<div class="secao"><h3>🔴 Diferenças de Custo Médio ({len(df_diferencas_custo)} itens)</h3>{tabela_custo_html}</div>' if len(df_diferencas_custo) > 0 else ''}
            
            {f'<div class="secao"><h3>📊 Diferenças de Saldo ({len(df_diferencas_saldo)} itens)</h3>{tabela_saldo_html}</div>' if len(df_diferencas_saldo) > 0 else ''}
            
            <div class="footer">
                <p>Este relatório foi gerado automaticamente pelo Analisador de Custo Médio.</p>
            </div>
        </body>
        </html>
        """
        return html

    def enviar_email(self, html_content):
        """Envia email com o relatório"""
        try:
            msg = MIMEMultipart('alternative')
            msg['Subject'] = f'Análise de Custo Médio - {datetime.now().strftime("%d/%m/%Y")}'
            msg['From'] = self.email_from
            msg['To'] = self.email_to

            html_part = MIMEText(html_content, 'html', 'utf-8')
            msg.attach(html_part)

            with smtplib.SMTP(self.smtp_server, self.smtp_port) as server:
                server.starttls()
                server.login(self.smtp_username, self.smtp_password)
                server.send_message(msg)

            print(f"✓ Email enviado com sucesso para {self.email_to}")
            return True
        except Exception as e:
            print(f"✗ Erro ao enviar email: {e}")
            return False

    def executar(self):
        """Executa o processo completo"""
        print("=" * 60)
        print("ANALISADOR DE CUSTO MÉDIO")
        print("=" * 60)
        print(f"Início: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}\n")

        if not self.conectar_sql_server():
            sys.exit(1)

        print("\n[1/4] Executando query...")
        df_original = self.executar_query()
        if df_original is None or len(df_original) == 0:
            print("✗ Nenhum dado encontrado.")
            return

        print(f"\n[2/4] Calculando custos médios...")
        df_resultados = self.calcular_custo_medio(df_original)
        print(f"✓ {len(df_resultados)} linhas processadas")

        print(f"\n[3/4] Identificando diferenças...")
        # Filtrar apenas movimentações efetivas (S) para buscar diferenças
        df_efetivas = df_resultados[df_resultados['MovimentoEfetivo'] == 'S']

        df_diferencas_custo = df_efetivas[df_efetivas['TemDiferencaCusto'] ==
                                          True]
        df_diferencas_saldo = df_efetivas[df_efetivas['TemDiferencaSaldo'] ==
                                          True]
        total_diferencas = len(df_diferencas_custo) + len(df_diferencas_saldo)

        print(
            f"✓ {len(df_diferencas_custo)} diferenças de custo médio encontradas (threshold: {self.diferenca_threshold})"
        )
        print(
            f"✓ {len(df_diferencas_saldo)} diferenças de saldo encontradas (threshold: 0.01)"
        )

        if len(df_diferencas_custo) > 0:
            print("\nItens com diferença de CUSTO MÉDIO:")
            for _, row in df_diferencas_custo.iterrows():
                print(f"  - {row['Item']} | Ordem {row['Ordem']} | "
                      f"Diferença: {row['DiferencaCusto']:.4f}")

        if len(df_diferencas_saldo) > 0:
            print("\nItens com diferença de SALDO:")
            for _, row in df_diferencas_saldo.iterrows():
                print(f"  - {row['Item']} | Ordem {row['Ordem']} | "
                      f"Original: {row['SaldoOriginal']:.2f} | "
                      f"Calculado: {row['SaldoCalculado']:.2f} | "
                      f"Diferença: {row['DiferencaSaldo']:.2f}")

        print(f"\n[4/4] Gerando e enviando relatório...")
        html_relatorio = self.gerar_relatorio_html(df_resultados)
        self.enviar_email(html_relatorio)

        if self.conn:
            self.conn.close()

        print(f"\n{'=' * 60}")
        print(
            f"Processo concluído: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}"
        )
        print("=" * 60)

    def executar_com_interface(self):
        """Executa o processo para interface web e retorna resultado"""
        if not self.conectar_sql_server():
            raise Exception("Falha ao conectar ao SQL Server")

        df_original = self.executar_query()
        if df_original is None or len(df_original) == 0:
            raise Exception(
                "Nenhum dado encontrado para o período selecionado")

        df_resultados = self.calcular_custo_medio(df_original)

        # Filtrar apenas movimentações efetivas (S) para buscar diferenças
        df_efetivas = df_resultados[df_resultados['MovimentoEfetivo'] == 'S']

        df_diferencas_custo = pd.DataFrame()
        df_diferencas_saldo = pd.DataFrame()

        if self.incluir_custo:
            df_diferencas_custo = df_efetivas[df_efetivas['TemDiferencaCusto']
                                              == True]

        if self.incluir_saldo:
            df_diferencas_saldo = df_efetivas[df_efetivas['TemDiferencaSaldo']
                                              == True]

        total_diferencas = len(df_diferencas_custo) + len(df_diferencas_saldo)

        html_relatorio = self.gerar_relatorio_html(df_resultados)
        self.enviar_email(html_relatorio)

        if self.conn:
            self.conn.close()

        resultado = f"Análise concluída! {total_diferencas} diferenças encontradas e email enviado."
        return resultado


if __name__ == "__main__":
    analyzer = CustoMedioAnalyzer()
    analyzer.executar()
