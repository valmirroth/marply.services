
#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Inicializa a interface web do Analisador de Custo Médio
Execute este arquivo para acessar a interface em http://localhost:5000
"""

import sys
import os

# Adicionar src ao path para importar módulos
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'src'))

from web_interface import app

if __name__ == '__main__':
    print("="*60)
    print("🌐 INTERFACE WEB - Analisador de Custo Médio")
    print("="*60)
    print("\n🚀 Iniciando servidor web...")
    print("📱 Acesse: http://localhost:5000")
    print("📱 Acesse externamente: http://0.0.0.0:5000")
    print("\n💡 Para parar o servidor, pressione Ctrl+C")
    print("="*60)
    
    app.run(host='0.0.0.0', port=5000, debug=False)
