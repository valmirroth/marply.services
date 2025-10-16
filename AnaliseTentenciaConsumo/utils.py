# ==========================================================
# ARQUIVO 5: utils.py
# Utilitários gerais
# ==========================================================

import os

def carregar_manual():
    """Carrega o HTML do manual se existir"""
    manual_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'manual_marini.html')
    if os.path.exists(manual_path):
        try:
            with open(manual_path, 'r', encoding='utf-8') as f:
                return f.read()
        except Exception as e:
            print(f"Erro ao carregar manual: {e}")
    
    return """
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Manual do Usuário</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            padding: 40px;
            max-width: 800px;
            margin: 0 auto;
            background: #f5f5f5;
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 { color: #8B2332; }
        .warning {
            background: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📖 Manual do Usuário</h1>
        <div class="warning">
            <strong>⚠️ Arquivo não encontrado</strong>
            <p>O arquivo <code>manual_marini.html</code> não foi encontrado.</p>
            <p>Coloque o arquivo na mesma pasta do aplicativo.</p>
        </div>
        <p><a href="/" style="color: #8B2332;">← Voltar para o sistema</a></p>
    </div>
</body>
</html>
    """

# ==========================================================