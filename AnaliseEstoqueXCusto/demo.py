#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ARQUIVO DE DEMONSTRAÇÃO
Este é apenas um exemplo para demonstrar a estrutura do código.
O script real está em src/custo_medio_analyzer.py e deve ser executado
no seu ambiente local com as credenciais corretas.
"""

import pandas as pd
from datetime import datetime

print("="*70)
print(" DEMONSTRAÇÃO - Analisador de Custo Médio")
print("="*70)
print("\n⚠️  IMPORTANTE:")
print("   Este é um SCRIPT STANDALONE para execução LOCAL")
print("   NÃO deve ser executado no Replit\n")
print("="*70)
print("\n📦 Arquivos do projeto:\n")
print("   src/custo_medio_analyzer.py  - Script principal")
print("   src/test_conexao.py         - Script de teste de configuração")
print("   .env.example                 - Template de configuração")
print("   requirements.txt             - Dependências Python")
print("   README.md                    - Documentação completa")

print("\n" + "="*70)
print(" COMO USAR")
print("="*70)
print("""
1. Baixe todos os arquivos deste Replit para seu servidor local

2. Instale as dependências:
   pip install -r requirements.txt

3. Configure suas credenciais:
   cp .env.example .env
   # Edite o arquivo .env com suas credenciais

4. Teste a configuração:
   python src/test_conexao.py

5. Execute o script principal:
   python src/custo_medio_analyzer.py

6. Agende execução diária às 5h:
   - Windows: Use Task Scheduler
   - Linux: Use cron (0 5 * * *)

Consulte o README.md para instruções detalhadas.
""")

print("="*70)
print("\n📊 Exemplo da lógica de cálculo:\n")

dados_exemplo = {
    'Ordem': [1, 2, 3, 4],
    'Item': ['001 - ITEM A', '001 - ITEM A', '001 - ITEM A', '001 - ITEM A'],
    'TipoMovimento': ['Inicial', 'Entrada', 'Entrada', 'Saida'],
    'QtdeMovimento': [0, 100, 50, 30],
    'SaldoEmpresa': [200, 300, 350, 320],
    'CustoMedio': [10.00, 11.00, 11.50, 11.50],
}

df = pd.DataFrame(dados_exemplo)
print(df.to_string(index=False))

print("\n💡 Lógica:")
print("   - Ordem 1: Saldo e custo inicial (200 unidades @ R$ 10,00)")
print("   - Ordem 2: Entrada de 100 @ R$ 15,00")
print("     Novo CM = (200×10 + 100×15) / 300 = R$ 11,67")
print("   - Ordem 3: Entrada de 50 @ R$ 13,00")  
print("     Novo CM = (300×11,67 + 50×13) / 350 = R$ 11,86")
print("   - Ordem 4: Saída de 30")
print("     CM mantém: R$ 11,86")

print("\n" + "="*70)
print(f"Data: {datetime.now().strftime('%d/%m/%Y %H:%M:%S')}")
print("="*70)
