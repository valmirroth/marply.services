# ==========================================================
# ARQUIVO 7: app.py (PRINCIPAL)
# Aplicação Flask e rotas
# ==========================================================

#=========================================================
# comando gerar exe produção
# C:\AnaliseRiscoConsumo\pj>py -3.11 -m PyInstaller --onefile --clean --noconfirm --noconsole --name AnaliseTendenciaConsumo --collect-all pandas --collect-all numpy --collect-all pyodbc --collect-all flask --hidden-import jinja2 --hidden-import markupsafe --hidden-import werkzeug app.py
#=========================================================
from flask import Flask, render_template_string, request, session, redirect, url_for
from config import SECRET_KEY, PORT
from auth import carregar_usuarios, login_required
from database import get_empresa_filter
from data_processor import carregar_dados_detalhado, carregar_dados_totalizado
from GravaLog import inserir_log_acesso
from utils import carregar_manual
from templates import LOGIN_HTML, MAIN_HTML, MANUAL_HTML

app = Flask(__name__)
app.secret_key = SECRET_KEY

@app.route("/login", methods=["GET", "POST"])
def login():
    error = None
    if request.method == "POST":
        username = request.form.get("username", "").strip()
        password = request.form.get("password", "").strip()
        empresa = request.form.get("empresa", "").strip()
        mes_atual = request.form.get("mes_atual", "").strip() 
        if not empresa:
            error = "Por favor, selecione uma empresa!"
        elif not username or not password:
            error = "Usuário e senha são obrigatórios!"
        else:
            usuarios = carregar_usuarios()
            
            if username in usuarios and usuarios[username] == password:
                session.clear()
                session['logged_in'] = True
                session['username'] = username
                session['empresa'] = empresa
                session['mes_atual'] = mes_atual
                empresa_nomes = {
                    'matriz': 'Matriz',
                    'filial': 'Filial',
                    'consolidado': 'Consolidado'
                }
                session['empresa_nome'] = empresa_nomes.get(empresa, 'Consolidado')
                inserir_log_acesso(mes_atual, empresa, num_meses=10, username=username)
                return redirect(url_for('index'))
            else:
                error = "Usuário ou senha inválidos!"
    
    return render_template_string(LOGIN_HTML, error=error)

@app.route("/logout")
def logout():
    session.clear()
    return redirect(url_for('login'))

@app.route("/manual")
def manual():
    return render_template_string(MANUAL_HTML, error="Erro ao carregar manual!")

@app.route("/")
@login_required
def index():
    try:
        if 'empresa' not in session:
            session.clear()
            return redirect(url_for('login'))
        
        view_mode = request.args.get('view', 'detalhado')
        num_meses = int(request.args.get('num_meses', 10))
        num_meses = max(3, min(24, num_meses))
        
        empresa_tipo = session.get('empresa', 'consolidado')
        mes_atual = session.get('mes_atual', 'não')
        empresa_filter = get_empresa_filter(empresa_tipo)
        
        if view_mode == 'totalizado':
            dados, meses = carregar_dados_totalizado(mes_atual, empresa_filter, num_meses)
        else:
            dados, meses = carregar_dados_detalhado(mes_atual, empresa_filter, num_meses)
        
        contas_unicas = sorted(dados['ContaContabil'].unique())
        
        return render_template_string(
            MAIN_HTML, 
            dados=dados, 
            meses=meses, 
            contas_unicas=contas_unicas,
            view_mode=view_mode,
            session=session,
            num_meses=num_meses
        )
    except Exception as e:
        import traceback
        return f"<pre>Erro ao gerar relatório:\n{traceback.format_exc()}</pre>", 500

if __name__ == "__main__":
    print("=" * 60)
    print("Sistema de Análise de Tendências de Consumo")
    print("=" * 60)
    
    usuarios = carregar_usuarios()
#    print(f"✓ Usuários: {len(usuarios)}")
#    print(f"✓ Servidor: http://0.0.0.0:{PORT}")
#    print(f"✓ Manual: http://0.0.0.0:{PORT}/manual")
#    print("=" * 60)
    
    app.run(host="0.0.0.0", port=PORT)