# ==========================================================
# TEMPLATE DE LOGIN
# ==========================================================
LOGIN_HTML = """
<!doctype html>
<html>
<head>
    <title>Login - Análise de Risco</title>
    <meta charset="UTF-8">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        
        .login-wrapper {
            flex: 1;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        
        .login-brand {
            text-align: center;
            margin-bottom: 40px;
            color: white;
        }
        
        .login-brand h1 {
            font-size: 42px;
            font-weight: 700;
            margin-bottom: 10px;
            text-shadow: 0 2px 10px rgba(0,0,0,0.3);
        }
        
        .login-brand p {
            font-size: 18px;
            opacity: 0.95;
            font-weight: 300;
        }
        
        .login-container {
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.4);
            width: 100%;
            max-width: 450px;
            overflow: hidden;
        }
        
        .login-header {
            background: rgba(255,255,255,0.1);
            backdrop-filter: blur(10px);
            border-bottom: 1px solid rgba(139, 35, 50, 0.2);
            color: #2d3748;
            padding: 30px;
            text-align: center;
        }
        
        .login-header h2 {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 8px;
            color: #8B2332;
        }
        
        .login-header p {
            font-size: 14px;
            color: #718096;
        }
        
        .login-form {
            padding: 40px 35px;
        }
        
        .form-group {
            margin-bottom: 25px;
        }
        
        .form-group label {
            display: block;
            font-weight: 600;
            color: #4a5568;
            font-size: 14px;
            margin-bottom: 8px;
        }
        
        .form-group input {
            width: 100%;
            padding: 14px 16px;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-size: 15px;
            transition: all 0.2s;
        }
        
        .form-group select {
            width: 100%;
            padding: 14px 16px;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-size: 15px;
            transition: all 0.2s;
            background: white;
            cursor: pointer;
        }
        
        .form-group input:focus,
        .form-group select:focus {
            outline: none;
            border-color: #8B2332;
            box-shadow: 0 0 0 4px rgba(139, 35, 50, 0.1);
        }
        
        .btn-login {
            width: 100%;
            padding: 14px;
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
            box-shadow: 0 4px 15px rgba(139, 35, 50, 0.3);
            margin-top: 10px;
        }
        
        .btn-login:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(139, 35, 50, 0.4);
        }
        
        .btn-login:active {
            transform: translateY(0);
        }
        
        .error-message {
            background: #fff5f5;
            color: #c53030;
            padding: 14px;
            border-radius: 8px;
            border-left: 4px solid #f56565;
            margin-bottom: 25px;
            font-size: 14px;
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .login-footer {
            text-align: center;
            padding: 25px;
            background: #f7fafc;
            color: #718096;
            font-size: 13px;
            border-top: 1px solid #e2e8f0;
        }
        
        .login-decoration {
            position: absolute;
            width: 300px;
            height: 300px;
            border-radius: 50%;
            background: rgba(255,255,255,0.05);
            z-index: 0;
        }
        
        .decoration-1 {
            top: -100px;
            left: -100px;
        }
        
        .decoration-2 {
            bottom: -150px;
            right: -150px;
            width: 400px;
            height: 400px;
        }
        
        @media (max-width: 768px) {
            .login-brand h1 {
                font-size: 32px;
            }
            
            .login-brand p {
                font-size: 16px;
            }
            
            .login-container {
                max-width: 100%;
            }
            
            .login-form {
                padding: 30px 25px;
            }
            
            .login-decoration {
                display: none;
            }
        }
    </style>
</head>
<body>
    <div class="login-decoration decoration-1"></div>
    <div class="login-decoration decoration-2"></div>
    
    <div class="login-wrapper">
        <div class="login-brand">
            <h1> Marply</h1>
            <p>Análise de Tendências de Consumo</p>
        </div>
        
        <div class="login-container">
            <div class="login-header">
                <h2>Bem-vindo(a)</h2>
                <p>Faça login para acessar o sistema</p>
            </div>
            
            <form class="login-form" method="POST">
                {% if error %}
                <div class="error-message">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" y1="8" x2="12" y2="12"></line>
                        <line x1="12" y1="16" x2="12.01" y2="16"></line>
                    </svg>
                    <span>{{ error }}</span>
                </div>
                {% endif %}
                
                <div class="form-group">
                    <label for="username">Usuário</label>
                    <input type="text" id="username" name="username" required autofocus placeholder="Digite seu usuário">
                </div>
                
                <div class="form-group">
                    <label for="password">Senha</label>
                    <input type="password" id="password" name="password" required placeholder="Digite sua senha">
                </div>
                
                <div class="form-group">
                    <label for="empresa">Empresa</label>
                    <select id="empresa" name="empresa" required>
                        <option value="">Selecione a empresa...</option>
                        <option value="matriz">🏢 Matriz</option>
                        <option value="filial">🏭 Filial</option>
                        <option value="consolidado">🌐 Consolidado (Todas)</option>
                    </select>
                </div>
                
                <div class="form-group">
                    <label for="mes_atual">Considerar Mês Atual</label>
                    <select id="mes_atual" name="mes_atual" required>
                        <option value="">Selecione uma opção...</option>
                        <option value="sim">✅ Sim - Incluir mês atual</option>
                        <option value="nao">❌ Não - Excluir mês atual</option>
                    </select>
                </div>
                
                <button type="submit" class="btn-login">Entrar no Sistema</button>
            </form>
            
            <div class="login-footer">
                Sistema de Análise de Tendência de Consumo © 2025
            </div>
        </div>
    </div>
</body>
</html>
"""

# ==========================================================
# TEMPLATE PRINCIPAL
# ==========================================================
MAIN_HTML = """
<!doctype html>
<html>
<head>
    <title>Análise de Tendências de Consumo</title>
    <meta charset="UTF-8">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            padding: 30px;
            min-height: 100vh;
        }
        
        .container {
            max-width: 100%;
            margin: 0 auto;
            background: white;
            border-radius: 15px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            color: white;
            padding: 30px;
            text-align: center;
            position: relative;
        }
        
        .header h2 {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 8px;
        }
        
        .header p {
            font-size: 14px;
            opacity: 0.9;
        }
        
        .header-buttons {
            position: absolute;
            top: 20px;
            right: 20px;
            display: flex;
            gap: 10px;
            align-items: center;
        }
        
        .logout-btn {
            padding: 8px 16px;
            background: rgba(255,255,255,0.2);
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            border-radius: 6px;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
            text-decoration: none;
            transition: all 0.2s;
            white-space: nowrap;
        }
        
        .logout-btn:hover {
            background: rgba(255,255,255,0.3);
        }
        
        .user-info {
            position: absolute;
            top: 20px;
            left: 20px;
            color: white;
            font-size: 13px;
            opacity: 0.9;
        }
        
        .view-toggle {
            background: white;
            padding: 15px 30px;
            border-bottom: 2px solid #e2e8f0;
            display: flex;
            justify-content: center;
            gap: 15px;
        }
        
        .toggle-btn {
            padding: 10px 24px;
            border: 2px solid #8B2332;
            background: white;
            color: #8B2332;
            border-radius: 6px;
            font-size: 14px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
            text-decoration: none;
            display: inline-block;
        }
        
        .toggle-btn:hover {
            background: #FFF5F7;
        }
        
        .toggle-btn.active {
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            color: white;
            border-color: #8B2332;
        }
        
        .filters {
            background: white;
            padding: 15px 30px;
            border-bottom: 3px solid #e2e8f0;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        
        .filter-group {
            display: flex;
            align-items: center;
            gap: 25px;
            flex-wrap: wrap;
        }
        
        .filter-item {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        
        .filter-group label {
            font-weight: 600;
            color: #4a5568;
            font-size: 13px;
            white-space: nowrap;
        }
        
        .input-group {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        
        .input-group span {
            color: #718096;
            font-size: 13px;
            font-weight: 500;
        }
        
        .input-group input[type="number"] {
            width: 120px;
        }
        
        .filter-group input,
        .filter-group select {
            padding: 8px 12px;
            border: 2px solid #e2e8f0;
            border-radius: 6px;
            font-size: 13px;
            transition: all 0.2s;
            background: white;
        }
        
        .filter-group select {
            min-width: 180px;
        }
        
        .filter-group input:hover,
        .filter-group select:hover {
            border-color: #cbd5e0;
        }
        
        .filter-group input:focus,
        .filter-group select:focus {
            outline: none;
            border-color: #8B2332;
            box-shadow: 0 0 0 3px rgba(139, 35, 50, 0.1);
        }
        
        .btn-clear {
            padding: 8px 16px;
            background: white;
            color: #e53e3e;
            border: 2px solid #e53e3e;
            border-radius: 6px;
            font-size: 13px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
            white-space: nowrap;
        }
        
        .btn-clear:hover {
            background: #e53e3e;
            color: white;
        }
        
        .stats {
            padding: 8px 16px;
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            color: white;
            border-radius: 6px;
            font-size: 13px;
            font-weight: 600;
            white-space: nowrap;
            margin-left: auto;
        }
        
        .stats span {
            font-size: 16px;
            font-weight: 700;
        }
        
        .table-wrapper {
            overflow-x: auto;
            padding: 20px;
        }
        
        table { 
            border-collapse: collapse;
            width: 100%;
            font-size: 13px;
            background: white;
        }
        
        th {
            background: #2d3748;
            color: white;
            padding: 14px 10px;
            text-align: center;
            font-weight: 600;
            text-transform: uppercase;
            font-size: 11px;
            letter-spacing: 0.5px;
            position: sticky;
            top: 0;
            z-index: 10;
        }
        
        td {
            padding: 12px 10px;
            text-align: right;
            border-bottom: 1px solid #e2e8f0;
            transition: background-color 0.2s;
        }
        
        td:first-child, td:nth-child(2) {
            text-align: left;
            font-weight: 500;
        }
        
        td:first-child {
            color: #4a5568;
            font-weight: 600;
        }
        
        td:nth-child(2) {
            color: #2d3748;
            max-width: 300px;
        }
        
        tbody tr:hover {
            background-color: #f7fafc;
        }
        
        .ultimo-mes {
            background-color: #e6f3ff !important;
            font-weight: 600;
            border-left: 2px solid #4299e1;
        }
        
        tbody tr:hover .ultimo-mes {
            background-color: #cce7ff !important;
        }
        
        .media-col {
            background-color: #fff5e6 !important;
            font-weight: 700;
            color: #2d3748;
            border-left: 2px solid #ed8936;
        }
        
        tbody tr:hover .media-col {
            background-color: #ffe8cc !important;
        }
        
        .ALTO {
            background-color: #fff5f5;
            border-left: 4px solid #f56565;
        }
        
        .MÉDIO {
            background-color: #fffbeb;
            border-left: 4px solid #ecc94b;
        }
        
        .BAIXO {
            background-color: #f0fff4;
            border-left: 4px solid #48bb78;
        }
        
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .badge-alto {
            background-color: #feb2b2;
            color: #742a2a;
        }
        
        .badge-medio {
            background-color: #fbd38d;
            color: #7c2d12;
        }
        
        .badge-baixo {
            background-color: #9ae6b4;
            color: #22543d;
        }
        
        .desvio-positivo {
            color: #e53e3e;
            font-weight: 600;
        }
        
        .desvio-negativo {
            color: #38a169;
            font-weight: 600;
        }
        
        @media (max-width: 768px) {
            body {
                padding: 15px;
            }
            
            .header {
                padding: 20px 15px;
            }
            
            .user-info {
                position: static;
                display: block;
                margin: 10px auto;
                text-align: center;
            }
            
            .header-buttons {
                position: static;
                justify-content: center;
                margin: 15px auto 0;
            }
            
            .logout-btn {
                font-size: 12px;
                padding: 6px 12px;
            }
            
            .view-toggle {
                padding: 10px 15px;
            }
            
            .toggle-btn {
                padding: 8px 16px;
                font-size: 12px;
            }
            
            .filters {
                padding: 15px;
            }
            
            .filter-group {
                flex-direction: column;
                align-items: stretch;
                gap: 15px;
            }
            
            .filter-item {
                flex-direction: column;
                align-items: flex-start;
                width: 100%;
            }
            
            .filter-group select,
            .filter-group input {
                width: 100%;
            }
            
            .input-group {
                width: 100%;
            }
            
            .input-group input[type="number"] {
                flex: 1;
            }
            
            .btn-clear {
                width: 100%;
            }
            
            .stats {
                width: 100%;
                text-align: center;
                margin-left: 0;
            }
            
            table {
                font-size: 11px;
            }
            
            th, td {
                padding: 8px 6px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="user-info">👤 {{ session.username }} | 🏢 {{ session.empresa_nome }}</div>
            <h2>📊 Análise de Tendências de Consumo</h2>
            <p>Monitoramento de Tendências de Consumo</p>
            <div class="header-buttons">
                <a href="{{ url_for('manual') }}" class="logout-btn" target="_blank" style="background: rgba(255,255,255,0.15);">📖 Manual</a>
                <a href="{{ url_for('logout') }}" class="logout-btn">Sair</a>
            </div>
        </div>
        
        <div class="view-toggle">
            <a href="/?view=detalhado&num_meses={{ num_meses }}" class="toggle-btn {% if view_mode == 'detalhado' %}active{% endif %}">
                📋 Visão Detalhada (Por Item)
            </a>
            <a href="/?view=totalizado&num_meses={{ num_meses }}" class="toggle-btn {% if view_mode == 'totalizado' %}active{% endif %}">
                📊 Visão Totalizada (Por Família)
            </a>
        </div>
        
        <div class="filters">
            <div class="filter-group">                
                <!-- 1. PERÍODO -->
                <div class="filter-item">
                    <label for="numMeses">📅 Período para Média:</label>
                    <select id="numMeses" onchange="atualizarMeses()">
                        <option value="3" {% if num_meses == 3 %}selected{% endif %}>3 meses</option>
                        <option value="6" {% if num_meses == 6 %}selected{% endif %}>6 meses</option>
                        <option value="10" {% if num_meses == 10 %}selected{% endif %}>10 meses</option>
                        <option value="12" {% if num_meses == 12 %}selected{% endif %}>12 meses</option>
                     
                    </select>
                </div>
                
                <!-- 2. FAMÍLIA -->
                <div class="filter-item">
                    <label for="filterConta">🏢 Família:</label>
                    <select id="filterConta" onchange="filtrarTabela()">
                        <option value="">Todas as famílias</option>
                        {% for conta in contas_unicas %}
                        <option value="{{ conta }}">{{ conta }}</option>
                        {% endfor %}
                    </select>
                </div>
                
                <!-- 3. BUSCAR ITEM (só na visão detalhada) -->
                {% if view_mode == 'detalhado' %}
                <div class="filter-item">
                    <label for="searchItem">🔍 Buscar Item:</label>
                    <input type="text" id="searchItem" placeholder="Digite para buscar..." onkeyup="filtrarTabela()" style="min-width: 200px;">
                </div>
                {% endif %}
                
                <!-- 4. FAIXA DE VALOR -->
                <div class="filter-item">
                    <label for="valorMin">💰 Valor Médio:</label>
                    <div class="input-group">
                        <input type="number" id="valorMin" placeholder="Mínimo" step="0.01" onchange="filtrarTabela()">
                        <span>até</span>
                        <input type="number" id="valorMax" placeholder="Máximo" step="0.01" onchange="filtrarTabela()">
                    </div>
                </div>
                
                <!-- BOTÃO LIMPAR -->
                <button class="btn-clear" onclick="limparFiltros()">✖ Limpar</button>
                
                <!-- CONTADOR (margem automática à direita) -->
                <div class="stats">
                    <span id="totalLinhas">{{ dados|length }}</span> registros
                </div>
            </div>
        </div>
        
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Família</th>
                        {% if view_mode == 'detalhado' %}
                        <th>Item / Observação</th>
                        {% endif %}
                        {% for m in meses[:-1] %}
                            <th>{{ m }}</th>
                        {% endfor %}
                        <th style="background: #2c5282;">{{ meses[-1] }}</th>
                        <th style="background: #c05621;">Média</th>
                        <th>Desvio %</th>
                        <th>Risco</th>
                    </tr>
                </thead>
                <tbody>
                    {% for _, r in dados.iterrows() %}
                    <tr class="{{ r.Risco }}">
                        <td>{{ r.ContaContabil }}</td>
                        {% if view_mode == 'detalhado' %}
                        <td>{{ r['Observação'] }}</td>
                        {% endif %}

                        {% for m in meses[:-1] %}
                            <td>{{ "{:,.2f}".format(r[m]) }}</td>
                        {% endfor %}
                        
                        <td class="ultimo-mes">{{ "{:,.2f}".format(r[meses[-1]]) }}</td>
                        <td class="media-col">{{ "{:,.2f}".format(r.MediaMensal) }}</td>
                        
                        <td class="{% if r.DesvioPercentual > 0 %}desvio-positivo{% else %}desvio-negativo{% endif %}">
                            {{ "{:+.1f}".format(r.DesvioPercentual) }}%
                        </td>
                        
                        <td style="text-align: center;">
                            <span class="badge badge-{{ r.Risco|lower }}">{{ r.Risco }}</span>
                        </td>
                    </tr>
                    {% endfor %}
                </tbody>
            </table>
        </div>
    </div>
    
    <script>
        const viewMode = '{{ view_mode }}';
        
        function atualizarMeses() {
            const numMeses = document.getElementById('numMeses').value;
            const urlParams = new URLSearchParams(window.location.search);
            urlParams.set('num_meses', numMeses);
            window.location.search = urlParams.toString();
        }
        
        function filtrarTabela() {
            const contaSelecionada = document.getElementById('filterConta').value.toLowerCase();
            const valorMin = parseFloat(document.getElementById('valorMin').value) || null;
            const valorMax = parseFloat(document.getElementById('valorMax').value) || null;
            
            let buscaItem = '';
            if (viewMode === 'detalhado') {
                buscaItem = document.getElementById('searchItem').value.toLowerCase();
            }
            
            const linhas = document.querySelectorAll('tbody tr');
            let contador = 0;
            
            linhas.forEach(linha => {
                const conta = linha.cells[0].textContent.toLowerCase();
                
                let item = '';
                if (viewMode === 'detalhado') {
                    item = linha.cells[1].textContent.toLowerCase();
                }
                
                const numColunas = linha.cells.length;
                const celulaMedia = linha.cells[numColunas - 3];
                const mediaTexto = celulaMedia.textContent.replace(/,/g, '');
                const mediaValor = parseFloat(mediaTexto) || 0;
                
                const matchConta = !contaSelecionada || conta.includes(contaSelecionada);
                const matchItem = !buscaItem || item.includes(buscaItem);
                const matchMin = valorMin === null || mediaValor >= valorMin;
                const matchMax = valorMax === null || mediaValor <= valorMax;
                
                if (matchConta && matchItem && matchMin && matchMax) {
                    linha.style.display = '';
                    contador++;
                } else {
                    linha.style.display = 'none';
                }
            });
            
            document.getElementById('totalLinhas').textContent = contador;
        }
        
        function limparFiltros() {
            document.getElementById('filterConta').value = '';
            if (viewMode === 'detalhado') {
                document.getElementById('searchItem').value = '';
            }
            document.getElementById('valorMin').value = '';
            document.getElementById('valorMax').value = '';
            filtrarTabela();
        }
    </script>
</body>
</html>
"""

MANUAL_HTML = """<!DOCTYPE html>
<html lang="pt-BR">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Manual do Usuário - Sistema de Análise de Tendências</title>
    <style>
        @page {
            size: A4;
            margin: 2cm;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            background: white;
        }

        .container {
            max-width: 210mm;
            margin: 0 auto;
            padding: 20px;
        }

        /* CAPA */
        .cover {
            page-break-after: always;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            text-align: center;
            background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
            color: white;
            padding: 40px;
        }

        .cover h1 {
            font-size: 48px;
            margin-bottom: 20px;
            font-weight: 700;
            color: #e6f3ff;
        }

        .cover h2 {
            font-size: 32px;
            margin-bottom: 40px;
            font-weight: 300;
        }

        .cover .logo {
            font-size: 72px;
            margin-bottom: 30px;
        }

        .cover .info {
            margin-top: 60px;
            font-size: 18px;
            opacity: 0.9;
        }

        /* CONTEÚDO */
        h1 {
            color: #8B2332;
            font-size: 32px;
            margin: 40px 0 20px 0;
            padding-bottom: 10px;
            border-bottom: 3px solid #8B2332;
        }

        h2 {
            color: #9e8487;
            font-size: 24px;
            margin: 30px 0 15px 0;
            padding-left: 10px;
            border-left: 4px solid #8B2332;
        }

        h3 {
            color: #333;
            font-size: 18px;
            margin: 20px 0 10px 0;
        }

        p {
            margin: 10px 0;
            text-align: justify;
        }

        ul,
        ol {
            margin: 15px 0 15px 30px;
        }

        li {
            margin: 8px 0;
        }

        .section {
            page-break-inside: avoid;
            margin-bottom: 30px;
        }

        .highlight {
            background: #fff5e6;
            padding: 15px;
            border-left: 4px solid #ed8936;
            margin: 20px 0;
        }

        .info-box {
            background: #e6f3ff;
            padding: 15px;
            border-left: 4px solid #4299e1;
            margin: 20px 0;
        }

        .warning-box {
            background: #fff5f5;
            padding: 15px;
            border-left: 4px solid #f56565;
            margin: 20px 0;
        }

        table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }

        th,
        td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: left;
        }

        th {
            background: #8B2332;
            color: white;
            font-weight: 600;
        }

        tr:nth-child(even) {
            background: #f9f9f9;
        }

        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: 700;
            text-transform: uppercase;
        }

        .badge-alto {
            background: #feb2b2;
            color: #742a2a;
        }

        .badge-medio {
            background: #fbd38d;
            color: #7c2d12;
        }

        .badge-baixo {
            background: #9ae6b4;
            color: #22543d;
        }

        .footer {
            position: fixed;
            bottom: 0;
            left: 0;
            right: 0;
            text-align: center;
            padding: 10px;
            font-size: 12px;
            color: #666;
            border-top: 1px solid #ddd;
        }

        .page-break {
            page-break-after: always;
        }

        code {
            background: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
        }

        .diagram {
            background: #f9f9f9;
            border: 2px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            font-family: 'Courier New', monospace;
            font-size: 14px;
            line-height: 1.8;
        }

        @media print {
            .cover {
                min-height: 297mm;
            }

            .no-print {
                display: none;
            }
        }
    </style>
</head>

<body>

    <!-- CAPA -->
    <div class="cover">
        <div class="logo">Marply</div>
        <h1>Sistema de Análise de</h1>
        <h1>Tendências de Consumo</h1>
        <h2>Manual do Usuário</h2>
        <div class="info">
            <p><strong>Marini Indústria de Compensados</strong></p>
            <p>Versão 1.0 | Janeiro 2026</p>
        </div>
    </div>

    <div class="container">

        <!-- ÍNDICE -->
        <div class="section">
            <h1>📑 Índice</h1>
            <ol>
                <li><strong>Descrição do Sistema</strong></li>
                <li><strong>Objetivos</strong></li>
                <li><strong>Funcionalidades Principais</strong></li>
                <li><strong>Arquitetura Técnica</strong></li>
                <li><strong>Regras de Negócio</strong></li>
                <li><strong>Manual do Usuário</strong></li>
                <li><strong>Casos de Uso</strong></li>
            </ol>
        </div>

        <div class="page-break"></div>

        <!-- 1. DESCRIÇÃO -->
        <div class="section">
            <h1>1. Descrição do Sistema</h1>

            <p>O <strong>Sistema de Análise de Tendências de Consumo</strong> é uma ferramenta web desenvolvida
                especialmente para a <strong>Marini Indústria de Compensados</strong>, com o objetivo de monitorar,
                analisar e prever padrões de consumo de materiais do almoxarifado.</p>

            <p>O sistema processa dados históricos de movimentação de estoque, calcula médias móveis configuráveis e
                identifica automaticamente desvios e tendências, permitindo uma gestão proativa de custos e recursos.
            </p>

            <div class="info-box">
                <strong>💡 Principais Benefícios:</strong>
                <ul>
                    <li>Identificação automática de anomalias no consumo</li>
                    <li>Previsibilidade para planejamento de compras</li>
                    <li>Redução de custos através de análise preditiva</li>
                    <li>Comparação entre unidades (Matriz e Filial)</li>
                    <li>Tomada de decisão baseada em dados concretos</li>
                </ul>
            </div>
        </div>

        <!-- 2. OBJETIVOS -->
        <div class="section">
            <h1>2. Objetivos</h1>

            <h2>2.1 Objetivo Geral</h2>
            <p>Fornecer uma plataforma integrada de <strong>Business Intelligence</strong> para análise de consumo de
                materiais, permitindo identificação rápida de desvios e facilitando a tomada
                de decisões estratégicas.</p>

            <h2>2.2 Objetivos Específicos</h2>
            <ul>
                <li><strong>Monitoramento Contínuo:</strong> Acompanhar o consumo mensal de todos os itens do
                    almoxarifado</li>
                <li><strong>Detecção de Anomalias:</strong> Identificar automaticamente itens com consumo acima do
                    esperado</li>
                <li><strong>Classificação de Riscos:</strong> Categorizar itens em níveis de risco (ALTO, MÉDIO, BAIXO)
                </li>
                <li><strong>Análise Comparativa:</strong> Permitir comparação entre diferentes períodos e unidades</li>
                <li><strong>Suporte à Decisão:</strong> Fornecer dados confiáveis para negociação com fornecedores</li>
                <li><strong>Controle de Custos:</strong> Auxiliar na redução e previsibilidade de gastos operacionais
                </li>
            </ul>
        </div>

        <div class="page-break"></div>

        <!-- 3. FUNCIONALIDADES -->
        <div class="section">
            <h1>3. Funcionalidades Principais</h1>

            <h2>3.1 Autenticação e Controle de Acesso</h2>
            <p>O sistema possui controle de acesso seguro através de login com usuário e senha. Cada usuário tem
                credenciais individuais e pode selecionar a empresa que deseja analisar no momento do login.</p>

            <div class="highlight">
                <strong>📋 Opções de Empresa:</strong>
                <ul>
                    <li><strong>🏢 Matriz:</strong> Visualiza apenas dados da empresa principal (RECNO = 1)</li>
                    <li><strong>🏭 Filial:</strong> Visualiza apenas dados da unidade filial (RECNO = 5)</li>
                    <li><strong>🌐 Consolidado:</strong> Visualiza dados de ambas as empresas simultaneamente</li>
                </ul>
            </div>

            <h2>3.2 Análise de Tendências Dinâmica</h2>

            <h3>Período Configurável</h3>
            <p>O usuário pode escolher o período de análise para cálculo da média:</p>
            <ul>
                <li><strong>3 meses:</strong> Análise de curto prazo, mais sensível a variações recentes</li>
                <li><strong>6 meses:</strong> Análise de médio prazo, equilibrada</li>
                <li><strong>10 meses:</strong> Período recomendado (padrão do sistema)</li>
                <li><strong>12 meses:</strong> Análise anual completa</li>

            </ul>

            <h3>Modos de Visualização</h3>
            <p>O sistema oferece duas formas de visualizar os dados:</p>

            <table>
                <tr>
                    <th style="width: 30%;">Modo</th>
                    <th>Descrição</th>
                    <th style="width: 25%;">Quando Usar</th>
                </tr>
                <tr>
                    <td><strong>📋 Visão Detalhada</strong></td>
                    <td>Exibe cada item do almoxarifado individualmente com seu histórico completo de consumo</td>
                    <td>Auditoria específica, investigação de itens</td>
                </tr>
                <tr>
                    <td><strong>📊 Visão Totalizada</strong></td>
                    <td>Agrupa itens por família de produtos, mostrando consumo consolidado por categoria</td>
                    <td>Visão macro, análise gerencial, tendências por categoria</td>
                </tr>
            </table>

            <h2>3.3 Classificação de Risco</h2>
            <p>O sistema classifica automaticamente cada item em três níveis de risco:</p>

            <table>
                <tr>
                    <th style="width: 20%;">Nível</th>
                    <th>Critério</th>
                    <th>Ação Recomendada</th>
                </tr>
                <tr>
                    <td><span class="badge badge-alto">🔴 ALTO</span></td>
                    <td>Consumo ≥ 150% da média <strong>E</strong> média ≥ R$ 1.000</td>
                    <td>Investigação imediata, ação corretiva urgente</td>
                </tr>
                <tr>
                    <td><span class="badge badge-medio">🟡 MÉDIO</span></td>
                    <td>Consumo ≥ 120% da média</td>
                    <td>Monitoramento próximo, análise de causa</td>
                </tr>
                <tr>
                    <td><span class="badge badge-baixo">🟢 BAIXO</span></td>
                    <td>Consumo dentro do esperado</td>
                    <td>Manutenção do monitoramento padrão</td>
                </tr>
            </table>

            <h2>3.4 Filtros Avançados</h2>
            <p>O sistema permite refinar a visualização através de múltiplos filtros:</p>

            <ul>
                <li><strong>📅 Período para Média:</strong> Define quantos meses usar no cálculo</li>
                <li><strong>🏢 Família:</strong> Filtra por categoria específica de produtos</li>
                <li><strong>🔍 Buscar Item:</strong> Pesquisa textual por nome do item (visão detalhada)</li>
                <li><strong>💰 Faixa de Valor:</strong> Define valores mínimo e máximo da média</li>
            </ul>

            <div class="info-box">
                <strong>💡 Dica:</strong> Os filtros podem ser combinados! Por exemplo: "Mostrar apenas família
                'Madeiras' com média entre R$ 5.000 e R$ 20.000 nos últimos 12 meses"
            </div>
        </div>

        <div class="page-break"></div>

        <!-- 4. ARQUITETURA -->
        <div class="section">
            <h1>4. Arquitetura Técnica</h1>

            <h2>4.1 Visão Geral</h2>
            <p>O sistema utiliza uma arquitetura web clássica de três camadas:</p>
            <pre class="diagram">
┌──────────────────────────────────────┐
│ CAMADA DE APRESENTAÇÃO               │
│                                      │
│ - Interface Web (HTML/CSS/JS)        │
│ - Responsiva e Intuitiva             │
│ - Filtros em Tempo Real              │
└──────────────────────────────────────┘
                ↓
┌──────────────────────────────────────┐
│ CAMADA DE APLICAÇÃO                  │
│                                      │
│ - Flask (Web Framework Python)       │
│ - Pandas (Processamento de Dados)    │
│ - Lógica de Negócio                  │
└──────────────────────────────────────┘
                ↓
┌──────────────────────────────────────────┐
│ CAMADA DE DADOS                          │
│                                          │
│ - SQL Server (Banco de Dados)            │
│ - Tabelas: CST_BI_ANALISE_CUSTEIO_GERAL, │
| ESTOQUE, GRUPOE                          │
└──────────────────────────────────────────┘
</pre>


            <h2>4.2 Fluxo de Processamento</h2>
            <ol>
                <li><strong>Autenticação:</strong> Usuário faz login e seleciona empresa</li>
                <li><strong>Consulta SQL:</strong> Sistema busca dados filtrados no banco (WHERE EMPRESA_RECNO IN ...)
                </li>
                <li><strong>Agregação:</strong> Pandas agrupa dados por mês e item/família</li>
                <li><strong>Cálculo de Métricas:</strong> Sistema calcula médias, desvios e classifica riscos</li>
                <li><strong>Ordenação:</strong> Dados ordenados por risco e desvio percentual</li>
                <li><strong>Renderização:</strong> Interface HTML exibe tabela interativa</li>
                <li><strong>Filtros Cliente:</strong> JavaScript aplica filtros sem recarregar página</li>
            </ol>

            <h2>4.3 Tecnologias Utilizadas</h2>
            <table>
                <tr>
                    <th>Componente</th>
                    <th>Tecnologia</th>
                    <th>Versão</th>
                </tr>
                <tr>
                    <td>Backend</td>
                    <td>Python + Flask</td>
                    <td>3.8+ / 3.0+</td>
                </tr>
                <tr>
                    <td>Análise de Dados</td>
                    <td>Pandas</td>
                    <td>2.0+</td>
                </tr>
                <tr>
                    <td>Banco de Dados</td>
                    <td>SQL Server</td>
                    <td>2016+</td>
                </tr>
                <tr>
                    <td>Conexão BD</td>
                    <td>pyODBC</td>
                    <td>4.0+</td>
                </tr>
                <tr>
                    <td>Frontend</td>
                    <td>HTML5 + CSS3 + JavaScript</td>
                    <td>-</td>
                </tr>
            </table>
        </div>

        <div class="page-break"></div>

        <!-- 5. REGRAS DE NEGÓCIO -->
        <div class="section">
            <h1>5. Regras de Negócio</h1>

            <h2>5.1 Cálculo da Média Mensal</h2>
            <p>A média mensal é calculada considerando os últimos N meses selecionados pelo usuário:</p>

            <div class="highlight">
                <strong>Fórmula:</strong><br>
                <code>Média Mensal = SOMA(Consumo dos últimos N meses) ÷ N</code>
                <br><br>
                <strong>Exemplo:</strong><br>
                Período: 6 meses<br>
                Consumos: R$ 5.000, R$ 4.800, R$ 5.200, R$ 4.900, R$ 5.100, R$ 7.500<br>
                Média = (5.000 + 4.800 + 5.200 + 4.900 + 5.100 + 7.500) ÷ 6 = <strong>R$ 5.416,67</strong>
            </div>

            <h2>5.2 Cálculo do Desvio Percentual</h2>
            <p>O desvio indica o quanto o último mês variou em relação à média:</p>

            <div class="highlight">
                <strong>Fórmula:</strong><br>
                <code>Desvio % = ((Valor Último Mês - Média Mensal) ÷ Média Mensal) × 100</code>
                <br><br>
                <strong>Exemplo:</strong><br>
                Média Mensal: R$ 5.416,67<br>
                Último Mês: R$ 7.500,00<br>
                Desvio = ((7.500 - 5.416,67) ÷ 5.416,67) × 100 = <strong>+38,5%</strong>
            </div>

            <div class="info-box">
                <strong>Interpretação:</strong>
                <ul>
                    <li><strong>Desvio Positivo (+):</strong> Consumo acima da média (vermelho)</li>
                    <li><strong>Desvio Negativo (-):</strong> Consumo abaixo da média (verde)</li>
                </ul>
            </div>

            <h2>5.3 Classificação de Risco</h2>
            <p>O algoritmo de classificação segue esta lógica:</p>

            <div class="diagram">
                SE Média Mensal ≤ 0:
                Risco = BAIXO

                SENÃO SE (Último Mês ≥ Média × 1,5) E (Média ≥ R$ 1.000):
                Risco = ALTO

                SENÃO SE Último Mês ≥ Média × 1,2:
                Risco = MÉDIO

                SENÃO:
                Risco = BAIXO
            </div>

            <h2>5.4 Exemplos Práticos</h2>

            <table>
                <tr>
                    <th>Média Mensal</th>
                    <th>Último Mês</th>
                    <th>Desvio %</th>
                    <th>Classificação</th>
                    <th>Motivo</th>
                </tr>
                <tr>
                    <td>R$ 5.000</td>
                    <td>R$ 8.000</td>
                    <td>+60%</td>
                    <td><span class="badge badge-alto">ALTO</span></td>
                    <td>8.000 ≥ 5.000×1,5 E 5.000 ≥ 1.000</td>
                </tr>
                <tr>
                    <td>R$ 800</td>
                    <td>R$ 1.500</td>
                    <td>+87,5%</td>
                    <td><span class="badge badge-medio">MÉDIO</span></td>
                    <td>Média < R$ 1.000 (não atinge ALTO)</td>
                </tr>
                <tr>
                    <td>R$ 10.000</td>
                    <td>R$ 12.500</td>
                    <td>+25%</td>
                    <td><span class="badge badge-medio">MÉDIO</span></td>
                    <td>12.500 ≥ 10.000×1,2 mas < 10.000×1,5</td>
                </tr>
                <tr>
                    <td>R$ 3.000</td>
                    <td>R$ 3.100</td>
                    <td>+3,3%</td>
                    <td><span class="badge badge-baixo">BAIXO</span></td>
                    <td>Variação dentro do esperado</td>
                </tr>
            </table>

            <div class="warning-box">
                <strong>⚠️ Importante:</strong> A classificação de risco considera tanto o percentual de aumento quanto
                o valor absoluto da média. Itens de baixo valor não são classificados como ALTO risco mesmo com grandes
                variações percentuais.
            </div>
        </div>

        <div class="page-break"></div>

        <!-- 6. MANUAL DO USUÁRIO -->
        <div class="section">
            <h1>6. Manual do Usuário</h1>

            <h2>6.1 Acessando o Sistema</h2>
            <ol>
                <li>Abra o navegador (Chrome, Firefox ou Edge)</li>
                <li>Digite o endereço: <code>https://CustoCerto.local</code></li>
                <li>A tela de login será exibida</li>
            </ol>

            <h2>6.2 Fazendo Login</h2>
            <ol>
                <li><strong>Usuário:</strong> Digite seu nome de usuário</li>
                <li><strong>Senha:</strong> Digite sua senha</li>
                <li><strong>Empresa:</strong> Selecione a empresa:
                    <ul>
                        <li><strong>🏢 Matriz:</strong> Para análise apenas da matriz</li>
                        <li><strong>🏭 Filial:</strong> Para análise apenas da filial</li>
                        <li><strong>🌐 Consolidado:</strong> Para análise de ambas</li>
                    </ul>
                </li>
                <li>Clique em <strong>"Entrar no Sistema"</strong></li>
            </ol>

            <h2>6.3 Navegando pela Interface</h2>

            <h3>Cabeçalho do Sistema</h3>
            <p>No topo da tela você verá:</p>
            <ul>
                <li><strong>Esquerda:</strong> Seu nome de usuário e empresa selecionada</li>
                <li><strong>Centro:</strong> Título "Análise de Tendências de Consumo"</li>
                <li><strong>Direita:</strong> Botão "Sair" para fazer logout</li>
            </ul>

            <h3>Botões de Visualização</h3>
            <p>Logo abaixo do cabeçalho, você encontra dois botões:</p>
            <ul>
                <li><strong>📋 Visão Detalhada (Por Item):</strong> Mostra cada item separadamente</li>
                <li><strong>📊 Visão Totalizada (Por Família):</strong> Mostra totais por categoria</li>
            </ul>

            <h3>Barra de Filtros</h3>
            <p>A barra de filtros permite refinar os dados exibidos:</p>

            <table>
                <tr>
                    <th style="width: 30%;">Filtro</th>
                    <th>Descrição</th>
                    <th style="width: 25%;">Dica de Uso</th>
                </tr>
                <tr>
                    <td><strong>📅 Período para Média</strong></td>
                    <td>Define quantos meses usar para calcular a média</td>
                    <td>Use 10 ou 12 meses para análise equilibrada</td>
                </tr>
                <tr>
                    <td><strong>🏢 Família</strong></td>
                    <td>Filtra por categoria de produtos</td>
                    <td>Útil para análise de categorias específicas</td>
                </tr>
                <tr>
                    <td><strong>🔍 Buscar Item</strong></td>
                    <td>Pesquisa por nome do item (apenas visão detalhada)</td>
                    <td>Digite parte do nome para encontrar rapidamente</td>
                </tr>
                <tr>
                    <td><strong>💰 Valor Médio</strong></td>
                    <td>Define faixa de valores (mínimo e máximo)</td>
                    <td>Ex: 1000 até 10000 para itens de médio valor</td>
                </tr>
            </table>

            <h2>6.4 Interpretando a Tabela</h2>

            <h3>Colunas da Tabela</h3>
            <table>
                <tr>
                    <th>Coluna</th>
                    <th>Significado</th>
                </tr>
                <tr>
                    <td><strong>Família</strong></td>
                    <td>Categoria do produto (ex: Abrasivos, Ferragens, etc.)</td>
                </tr>
                <tr>
                    <td><strong>Item/Observação</strong></td>
                    <td>Nome completo do item (só aparece na visão detalhada)</td>
                </tr>
                <tr>
                    <td><strong>Meses (colunas)</strong></td>
                    <td>Valor consumido em cada mês histórico</td>
                </tr>
                <tr>
                    <td><strong>Último Mês</strong></td>
                    <td>Consumo mais recente (destacado em azul)</td>
                </tr>
                <tr>
                    <td><strong>Média</strong></td>
                    <td>Média calculada do período selecionado (destacado em laranja)</td>
                </tr>
                <tr>
                    <td><strong>Desvio %</strong></td>
                    <td>Variação percentual do último mês em relação à média</td>
                </tr>
                <tr>
                    <td><strong>Risco</strong></td>
                    <td>Classificação: ALTO, MÉDIO ou BAIXO</td>
                </tr>
            </table>

            <h3>Cores e Destaques</h3>
            <ul>
                <li><strong>Fundo Azul Claro:</strong> Coluna do último mês</li>
                <li><strong>Fundo Laranja Claro:</strong> Coluna da média</li>
                <li><strong>Borda Vermelha:</strong> Linha de risco ALTO</li>
                <li><strong>Borda Amarela:</strong> Linha de risco MÉDIO</li>
                <li><strong>Borda Verde:</strong> Linha de risco BAIXO</li>
                <li><strong>Número Vermelho:</strong> Desvio positivo (aumento)</li>
                <li><strong>Número Verde:</strong> Desvio negativo (redução)</li>
            </ul>

            <h2>6.5 Ações Comuns</h2>

            <h3>Alternar Entre Visões</h3>
            <ol>
                <li>Clique em <strong>"📋 Visão Detalhada"</strong> para ver itens individuais</li>
                <li>Clique em <strong>"📊 Visão Totalizada"</strong> para ver totais por família</li>
                <li>Os filtros são mantidos ao trocar de visão</li>
            </ol>

            <h3>Alterar Período de Análise</h3>
            <ol>
                <li>No filtro <strong>"Período para Média"</strong>, selecione o período desejado</li>
                <li>A página será recarregada automaticamente</li>
                <li>Todas as médias e classificações serão recalculadas</li>
            </ol>

            <h3>Filtrar por Família</h3>
            <ol>
                <li>No filtro <strong>"Família"</strong>, selecione a categoria desejada</li>
                <li>A tabela mostrará apenas itens daquela família</li>
                <li>O contador de registros será atualizado</li>
            </ol>

            <h3>Buscar Item Específico (Visão Detalhada)</h3>
            <ol>
                <li>No campo <strong>"Buscar Item"</strong>, digite parte do nome</li>
                <li>A busca é feita em tempo real enquanto você digita</li>
                <li>Não é necessário digitar o nome completo</li>
            </ol>

            <h3>Filtrar por Valor</h3>
            <ol>
                <li>Digite o valor mínimo (opcional)</li>
                <li>Digite o valor máximo (opcional)</li>
                <li>Pressione Enter ou clique fora do campo</li>
                <li>Exemplos:
                    <ul>
                        <li>Apenas mínimo (1000): Mostra médias ≥ R$ 1.000</li>
                        <li>Apenas máximo (5000): Mostra médias ≤ R$ 5.000</li>
                        <li>Ambos (1000 até 5000): Mostra entre R$ 1.000 e R$ 5.000</li>
                    </ul>
                </li>
            </ol>

            <h3>Limpar Todos os Filtros</h3>
            <ol>
                <li>Clique no botão <strong>"✖ Limpar"</strong></li>
                <li>Todos os filtros voltarão ao padrão</li>
                <li>A tabela mostrará todos os registros novamente</li>
            </ol>

            <div class="info-box">
                <strong>💡 Dica Profissional:</strong> Combine filtros para análises específicas! Exemplo: "Família =
                Abrasivos" + "Valor entre 500 e 10000" + "Período = 12 meses" para uma análise anual detalhada de
                abrasivos de médio/alto valor.
            </div>
        </div>

        <div class="page-break"></div>

        <!-- 7. CASOS DE USO -->
        <div class="section">
            <h1>7. Casos de Uso Práticos</h1>

            <h2>7.1 Caso de Uso: Análise Mensal de Custos</h2>

            <h3>Objetivo</h3>
            <p>Revisar os custos do mês anterior e identificar itens que precisam de atenção imediata.</p>

            <h3>Passo a Passo</h3>
            <ol>
                <li>Faça login no sistema</li>
                <li>Selecione <strong>"Consolidado"</strong> para visão completa</li>
                <li>Mantenha o período em <strong>"10 meses"</strong></li>
                <li>Observe a tabela que já está ordenada por risco</li>
                <li>Analise primeiro os itens de <span class="badge badge-alto">ALTO RISCO</span></li>
                <li>Anote ou exporte os dados relevantes</li>
                <li>Agende reunião com gestores para discussão</li>
            </ol>

            <h3>Pontos de Atenção</h3>
            <ul>
                <li>Verifique se há projetos especiais que justifiquem o aumento</li>
                <li>Compare com o mesmo período do ano anterior</li>
                <li>Investigue causas de desvios acima de 50%</li>
            </ul>

            <h2>7.2 Caso de Uso: Comparação Matriz x Filial</h2>

            <h3>Objetivo</h3>
            <p>Comparar o consumo entre as duas unidades e identificar diferenças significativas.</p>

            <h3>Passo a Passo</h3>
            <ol>
                <li><strong>Análise da Matriz:</strong>
                    <ul>
                        <li>Faça login selecionando <strong>"Matriz"</strong></li>
                        <li>Alterne para <strong>"Visão Totalizada"</strong></li>
                        <li>Anote os valores das principais famílias</li>
                        <li>Identifique itens de alto risco</li>
                    </ul>
                </li>
                <li><strong>Análise da Filial:</strong>
                    <ul>
                        <li>Faça logout e login novamente</li>
                        <li>Selecione <strong>"Filial"</strong></li>
                        <li>Repita o processo de análise</li>
                        <li>Anote os mesmos dados para comparação</li>
                    </ul>
                </li>
                <li><strong>Comparação:</strong>
                    <ul>
                        <li>Compare família por família</li>
                        <li>Identifique discrepâncias significativas</li>
                        <li>Investigue causas das diferenças</li>
                        <li>Proponha padronização de processos</li>
                    </ul>
                </li>
            </ol>

            <h3>Métricas Importantes</h3>
            <table>
                <tr>
                    <th>Métrica</th>
                    <th>O que analisar</th>
                </tr>
                <tr>
                    <td>Consumo per capita</td>
                    <td>Valor total ÷ número de funcionários</td>
                </tr>
                <tr>
                    <td>Eficiência de uso</td>
                    <td>Produção ÷ consumo de materiais</td>
                </tr>
                <tr>
                    <td>Desperdício</td>
                    <td>Itens com alto consumo e baixa produção</td>
                </tr>
            </table>

            <h2>7.3 Caso de Uso: Análise de Sazonalidade</h2>

            <h3>Objetivo</h3>
            <p>Identificar padrões sazonais no consumo para melhor planejamento de compras.</p>

            <h3>Passo a Passo</h3>
            <ol>
                <li>Configure o período para <strong>"12 meses"</strong></li>
                <li>Selecione <strong>"Visão Totalizada"</strong></li>
                <li>Analise família por família observando as colunas mensais</li>
                <li>Identifique meses de pico de consumo</li>
                <li>Identifique meses de baixo consumo</li>
                <li>Documente os padrões encontrados</li>
            </ol>

            <h3>Benefícios</h3>
            <ul>
                <li><strong>Negociação:</strong> Comprar antecipado em períodos de baixa demanda</li>
                <li><strong>Estoque:</strong> Ajustar níveis de estoque conforme sazonalidade</li>
                <li><strong>Fornecedores:</strong> Negociar contratos com base em volume anual</li>
                <li><strong>Fluxo de Caixa:</strong> Planejar desembolsos nos meses corretos</li>
            </ul>

            <div class="highlight">
                <strong>Exemplo Prático:</strong><br>
                Se você identificar que o consumo de "PECAS RODANTES" sempre aumenta em Outubro, Novembro e Dezembro,
                pode:
                <ul>
                    <li>Negociar com desconto</li>
                    <li>Comprar estoque estratégico em Setembro</li>
                </ul>
            </div>

            <h2>7.4 Caso de Uso: Auditoria de Família Específica</h2>

            <h3>Objetivo</h3>
            <p>Investigar em detalhes uma categoria específica que apresentou problemas.</p>

            <h3>Passo a Passo</h3>
            <ol>
                <li>No filtro <strong>"Família"</strong>, selecione a categoria a auditar</li>
                <li>Alterne para <strong>"Visão Detalhada"</strong></li>
                <li>Configure período para <strong>"12 meses"</strong></li>
                <li>Observe cada item individualmente</li>
                <li>Identifique os itens com maior desvio</li>
                <li>Use o filtro de valor para focar em itens relevantes</li>
                <li>Documente achados e recomendações</li>
            </ol>

            <h3>Checklist de Investigação</h3>
            <div class="warning-box">
                <strong>🔍 Para cada item de ALTO risco, verifique:</strong>
                <ul>
                    <li>☐ Houve projetos especiais no período?</li>
                    <li>☐ O preço do fornecedor aumentou?</li>
                    <li>☐ Houve mudança no processo produtivo?</li>
                    <li>☐ Existe desperdício ou uso inadequado?</li>
                    <li>☐ Há alternativas mais econômicas?</li>
                    <li>☐ O consumo está dentro do esperado para a produção?</li>
                </ul>
            </div>

            <h2>7.5 Caso de Uso: Preparação para Orçamento Anual</h2>

            <h3>Objetivo</h3>
            <p>Coletar dados históricos para planejar o orçamento do próximo ano.</p>

            <h3>Passo a Passo</h3>
            <ol>
                <li>Configure período para <strong>"12 meses"</strong></li>
                <li>Selecione <strong>"Consolidado"</strong></li>
                <li>Use <strong>"Visão Totalizada"</strong></li>
                <li>Para cada família, calcule:
                    <ul>
                        <li>Total anual do ano corrente</li>
                        <li>Crescimento percentual</li>
                        <li>Projeção para próximo ano</li>
                    </ul>
                </li>
                <li>Considere itens de alto risco para ajustes</li>
                <li>Adicione margem de segurança (5-10%)</li>
            </ol>

            <h3>Fórmula de Projeção</h3>
            <div class="highlight">
                <strong>Projeção Simples:</strong><br>
                <code>Orçamento Ano N+1 = (Média dos últimos 12 meses × 12) × 1,05</code>
                <br><br>
                <strong>Projeção com Crescimento:</strong><br>
                <code>Orçamento = Total Ano N × (1 + Taxa de Crescimento) × 1,05</code>
            </div>
        </div>

        <div class="page-break"></div>

        <!-- APÊNDICES -->
        <div class="section">
            <h1>Apêndice A: Glossário de Termos</h1>

            <table>
                <tr>
                    <th style="width: 30%;">Termo</th>
                    <th>Definição</th>
                </tr>
                <tr>
                    <td><strong>Média Mensal</strong></td>
                    <td>Valor médio de consumo calculado com base nos últimos N meses selecionados</td>
                </tr>
                <tr>
                    <td><strong>Desvio Percentual</strong></td>
                    <td>Variação percentual do consumo do último mês em relação à média calculada</td>
                </tr>
                <tr>
                    <td><strong>Risco</strong></td>
                    <td>Classificação automática do item em ALTO, MÉDIO ou BAIXO com base em regras pré-definidas</td>
                </tr>
                <tr>
                    <td><strong>Família</strong></td>
                    <td>Categoria ou grupo de produtos similares (ex: Abrasivos, Ferragens, Químicos)</td>
                </tr>
                <tr>
                    <td><strong>Visão Detalhada</strong></td>
                    <td>Modo de visualização que exibe cada item individualmente com seu histórico completo</td>
                </tr>
                <tr>
                    <td><strong>Visão Totalizada</strong></td>
                    <td>Modo de visualização que agrupa itens por família, mostrando totais consolidados</td>
                </tr>
                <tr>
                    <td><strong>Consolidado</strong></td>
                    <td>Visualização que inclui dados de Matriz e Filial simultaneamente</td>
                </tr>
                <tr>
                    <td><strong>Sazonalidade</strong></td>
                    <td>Padrões de variação de consumo que se repetem em determinados períodos do ano</td>
                </tr>
                <tr>
                    <td><strong>Período para Média</strong></td>
                    <td>Quantidade de meses históricos usados para calcular a média de referência</td>
                </tr>
            </table>
        </div>

        <div class="section">
            <h1>Apêndice B: Perguntas Frequentes (FAQ)</h1>

            <h3>1. Por que meu item está classificado como MÉDIO risco mesmo com desvio alto?</h3>
            <p><strong>R:</strong> Para ser classificado como ALTO risco, o item precisa atender dois critérios
                simultaneamente: desvio ≥ 150% E média mensal ≥ R$ 1.000. Itens de baixo valor não atingem ALTO risco
                para evitar alarmes falsos em itens de pouco impacto financeiro.</p>

            <h3>2. Posso exportar os dados para Excel?</h3>
            <p><strong>R:</strong> Atualmente o sistema não possui exportação automática. Você pode copiar e colar os
                dados da tabela diretamente no Excel ou usar a função de impressão do navegador para gerar um PDF.</p>

            <h3>3. Qual período devo usar para a análise?</h3>
            <p><strong>R:</strong> Depende do objetivo:</p>
            <ul>
                <li><strong>3-6 meses:</strong> Análise de curto prazo, mais sensível a mudanças recentes</li>
                <li><strong>10-12 meses:</strong> Análise equilibrada (recomendado para uso geral)</li>
            </ul>

            <h3>4. Por que alguns meses aparecem com valor zero?</h3>
            <p><strong>R:</strong> Pode indicar que não houve consumo do item naquele mês, ou que não há dados
                registrados no sistema para aquele período.</p>

            <h3>5. Como sei se um aumento é justificado?</h3>
            <p><strong>R:</strong> Verifique:</p>
            <ul>
                <li>Se houve aumento de produção no período</li>
                <li>Se ocorreram projetos especiais</li>
                <li>Se houve aumento de preço pelo fornecedor</li>
                <li>Se há registro de desperdício ou problema operacional</li>
            </ul>

            <h3>6. Posso ver dados de anos anteriores?</h3>
            <p><strong>R:</strong> O sistema exibe periodos após a data de 01/01/2025.</p>

            <h3>7. Como funciona o filtro de busca?</h3>
            <p><strong>R:</strong> O filtro busca em tempo real enquanto você digita. Não é necessário digitar o nome
                completo - digite qualquer parte do nome e o sistema encontrará todos os itens que contenham aquele
                texto.</p>

            <h3>8. Por que a página recarrega ao mudar o período?</h3>
            <p><strong>R:</strong> A mudança de período requer recalcular todas as médias e classificações com base nos
                novos dados, por isso é necessário buscar novamente as informações do banco de dados.</p>

            <h3>9. Posso ter múltiplos usuários logados ao mesmo tempo?</h3>
            <p><strong>R:</strong> Sim, cada usuário tem sua própria sessão independente. Um usuário pode estar
                visualizando a Matriz enquanto outro visualiza a Filial.</p>

            <h3>10. O que fazer quando encontro um item de ALTO risco?</h3>
            <p><strong>R:</strong> Siga este checklist:</p>
            <ol>
                <li>Documente o desvio encontrado</li>
                <li>Investigue a causa do aumento</li>
                <li>Consulte o setor responsável</li>
                <li>Avalie se é um problema ou situação esperada</li>
                <li>Tome ação corretiva se necessário</li>
                <li>Continue monitorando nos próximos meses</li>
            </ol>
        </div>

        <div class="page-break"></div>

        <div class="section">
            <h1>Apêndice C: Melhores Práticas</h1>

            <h2>Para Gestores</h2>
            <ul>
                <li>✅ Revise o sistema <strong>mensalmente</strong> para acompanhar tendências a cada fechamento de
                    período</li>
                <li>✅ Mantenha registro de itens que foram investigados e suas resoluções</li>
                <li>✅ Use a <strong>Visão Totalizada</strong> para reuniões executivas</li>
                <li>✅ Estabeleça meta de redução de itens de ALTO risco mês a mês</li>
            </ul>

            <h2>Para Compradores</h2>
            <ul>
                <li>✅ Use períodos longos para identificar sazonalidade</li>
                <li>✅ Monitore itens de alto valor (use filtro de valor mínimo)</li>
                <li>✅ Documente negociações com fornecedores baseadas nos dados</li>
                <li>✅ Planeje compras antecipadas em períodos de baixa</li>
            </ul>

            <h2>Para Auditoria</h2>
            <ul>
                <li>✅ Foque em itens de <span class="badge badge-alto">ALTO RISCO</span> com desvio > 100%</li>
                <li>✅ Use a <strong>Visão Detalhada</strong> para análises específicas</li>
                <li>✅ Combine filtros para análises cirúrgicas</li>
                <li>✅ Compare dados entre Matriz e Filial sistematicamente</li>
                <li>✅ Mantenha histórico de achados para referência futura</li>
            </ul>

            <h2>Frequência Recomendada de Análise</h2>
            <table>
                <tr>
                    <th>Tipo de Análise</th>
                    <th>Frequência</th>
                </tr>
                <tr>
                    <td>Monitoramento Geral</td>
                    <td>Mensal</td>
                </tr>
                <tr>
                    <td>Análise Detalhada</td>
                    <td>Mensal</td>
                </tr>
                <tr>
                    <td>Comparação Matriz x Filial</td>
                    <td>Mensal</td>
                </tr>
                <tr>
                    <td>Planejamento Orçamentário</td>
                    <td>Anual</td>

                </tr>
                <tr>
                    <td>Auditoria Completa</td>
                    <td>Trimestral</td>

                </tr>
            </table>
        </div>

        <div class="page-break"></div>

        <!-- NOTAS FINAIS -->
        <div class="section" style="text-align: center; margin-top: 100px;">
            <h1>📊</h1>
            <h2>Sistema de Análise de Tendências de Consumo</h2>
            <p style="margin-top: 40px; font-size: 18px;">
                <strong>Marini Indústria de Compensados</strong><br>
                Versão 1.0 | Janeiro 2026
            </p>

            <div style="margin-top: 60px; padding: 30px; background: #f9f9f9; border-radius: 10px;">
                <p style="font-size: 16px; line-height: 2;">
                    <strong>Desenvolvido para otimizar a gestão de custos<br>
                        e transformar dados em decisões estratégicas</strong>
                </p>
            </div>

            <div style="margin-top: 60px; font-size: 14px; color: #666;">
                <p>© 2026 Marini Indústria de Compensados Ltda.</p>
                <p>Todos os direitos reservados</p>
            </div>
        </div>

    </div>

    <!-- BOTÃO DE IMPRESSÃO (não aparece no PDF) -->
    <div class="no-print" style="position: fixed; bottom: 20px; right: 20px;">
        <button onclick="window.print()" style="
        padding: 15px 30px;
        background: linear-gradient(135deg, #8B2332 0%, #5A1820 100%);
        color: white;
        border: none;
        border-radius: 8px;
        font-size: 16px;
        font-weight: 600;
        cursor: pointer;
        box-shadow: 0 4px 12px rgba(139, 35, 50, 0.3);
    ">
            🖨️ Imprimir / Salvar PDF
        </button>
    </div>

</body>

</html>
"""