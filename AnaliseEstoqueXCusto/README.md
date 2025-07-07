# Analisador de Custo Médio - SQL Server

Script Python standalone que calcula o custo médio de itens baseado em movimentações do SQL Server, compara com os valores originais e envia email com diferenças encontradas.

## 📋 Requisitos

### Sistema Operacional
- Windows, Linux ou macOS

### Software Necessário
- **Python 3.8 ou superior**
- **Driver ODBC para SQL Server**
  - Windows: [Microsoft ODBC Driver 17 for SQL Server](https://docs.microsoft.com/en-us/sql/connect/odbc/download-odbc-driver-for-sql-server)
  - Linux: `sudo apt-get install unixodbc-dev`
- **Acesso ao SQL Server** (porta 1433 ou customizada)
- **Servidor SMTP** para envio de emails

## 🚀 Instalação

### 1. Clone ou baixe os arquivos
```bash
# Baixe os arquivos para uma pasta no seu servidor
mkdir /opt/custo_medio_analyzer
cd /opt/custo_medio_analyzer
```

### 2. Instale o Python (se não estiver instalado)
```bash
# Windows: Baixe de python.org
# Linux Ubuntu/Debian:
sudo apt-get update
sudo apt-get install python3 python3-pip
```

### 3. Instale as dependências Python
```bash
pip install -r requirements.txt
```

### 4. Configure as credenciais

Copie o arquivo de exemplo e edite com suas credenciais:
```bash
cp .env.example .env
```

Edite o arquivo `.env` com suas informações:
```bash
# No Windows
notepad .env

# No Linux
nano .env
```

**Exemplo de configuração:**
```ini
# SQL Server
SQL_SERVER=192.168.1.100
SQL_DATABASE=SEU_BANCO
SQL_USERNAME=usuario_sql
SQL_PASSWORD=senha_sql
SQL_PORT=1433

# SMTP (exemplo Gmail)
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=seuemail@gmail.com
SMTP_PASSWORD=sua_senha_app
EMAIL_FROM=seuemail@gmail.com
EMAIL_TO=ti01@marply.com.br

# Threshold
DIFERENCA_THRESHOLD=0.05
```

### 5. Teste a execução
```bash
python custo_medio_analyzer.py
```

## ⏰ Agendamento Automático (5h da manhã todos os dias)

### No Windows (Task Scheduler)

1. Abra o **Agendador de Tarefas** (Task Scheduler)
2. Clique em **Criar Tarefa Básica**
3. Configure:
   - **Nome:** Analisador Custo Médio
   - **Gatilho:** Diariamente às 05:00
   - **Ação:** Iniciar um programa
   - **Programa:** `C:\Python\python.exe` (caminho do seu Python)
   - **Argumentos:** `custo_medio_analyzer.py`
   - **Iniciar em:** `C:\caminho\para\pasta\do\script`

**Ou via linha de comando:**
```cmd
schtasks /create /tn "Analisador Custo Medio" /tr "python C:\caminho\custo_medio_analyzer.py" /sc daily /st 05:00
```

### No Linux (Cron)

1. Abra o crontab:
```bash
crontab -e
```

2. Adicione a linha:
```bash
0 5 * * * /usr/bin/python3 /opt/custo_medio_analyzer/custo_medio_analyzer.py >> /var/log/custo_medio.log 2>&1
```

**Explicação:**
- `0 5 * * *` = Todos os dias às 5h da manhã
- `>> /var/log/custo_medio.log 2>&1` = Salva logs em arquivo

3. Salve e saia (Ctrl+X, depois Y, depois Enter)

4. Verifique se foi adicionado:
```bash
crontab -l
```

## 📧 Configuração de Email

### Gmail
1. Ative a **verificação em 2 etapas** na sua conta Google
2. Crie uma **senha de app**: https://myaccount.google.com/apppasswords
3. Use essa senha no campo `SMTP_PASSWORD` do arquivo `.env`

### Outlook/Office 365
```ini
SMTP_SERVER=smtp.office365.com
SMTP_PORT=587
SMTP_USERNAME=seuemail@outlook.com
SMTP_PASSWORD=sua_senha
```

### Servidor SMTP Corporativo
Consulte seu departamento de TI para obter:
- Endereço do servidor SMTP
- Porta (geralmente 587 ou 25)
- Credenciais de autenticação

## 🔍 Como Funciona

1. **Conecta ao SQL Server** e executa a query de movimentações do mês atual
2. **Para cada item:**
   - Lê a linha com Ordem = 1 (saldo e custo médio inicial)
   - A partir da Ordem = 2, calcula:
     - **Entrada:** Atualiza saldo e recalcula custo médio
     - **Saída:** Reduz saldo, mantém custo médio
3. **Compara** o custo médio calculado com o da query original
4. **Identifica** diferenças absolutas maiores que 0.05
5. **Envia email** com relatório HTML das diferenças

### Fórmula do Custo Médio
```
Custo Médio = (Saldo Anterior × Custo Médio Anterior + Qtde Entrada × Custo Unitário) / (Saldo Anterior + Qtde Entrada)
```

## 📊 Exemplo de Saída

```
============================================================
ANALISADOR DE CUSTO MÉDIO
============================================================
Início: 02/10/2025 05:00:00

✓ Conectado ao SQL Server: servidor/banco
✓ Query executada com sucesso: 1500 registros encontrados
✓ 1500 linhas processadas
✓ 3 diferenças encontradas (threshold: 0.05)

Itens com diferença:
  - 001234 - ITEM EXEMPLO A | Ordem 5 | Diferença: 0.0823
  - 001235 - ITEM EXEMPLO B | Ordem 3 | Diferença: 0.1245
  - 001236 - ITEM EXEMPLO C | Ordem 7 | Diferença: 0.0567

✓ Email enviado com sucesso para ti01@marply.com.br
============================================================
```

## 🐛 Solução de Problemas

### Erro: "Driver não encontrado"
```bash
# Windows: Instale ODBC Driver 17 for SQL Server
# Linux:
sudo apt-get install unixodbc-dev
```

### Erro de conexão SQL Server
- Verifique se o servidor aceita conexões remotas
- Confirme que o firewall permite a porta 1433
- Teste a conexão: `telnet servidor 1433`

### Email não enviado
- Verifique credenciais SMTP
- Para Gmail, use senha de app (não a senha normal)
- Confirme que a porta SMTP não está bloqueada

### Logs de execução
```bash
# Linux - ver últimos logs
tail -f /var/log/custo_medio.log

# Windows - adicionar logging ao script
python custo_medio_analyzer.py > log.txt 2>&1
```

## 📝 Logs

O script gera logs no console durante a execução. Para salvar em arquivo:

**Linux:**
```bash
python custo_medio_analyzer.py >> /var/log/custo_medio.log 2>&1
```

**Windows:**
```cmd
python custo_medio_analyzer.py >> C:\logs\custo_medio.log 2>&1
```

## 🔐 Segurança

- ⚠️ **Nunca commite o arquivo `.env`** em repositórios Git
- ✓ Use senhas fortes para SQL Server e SMTP
- ✓ Restrinja permissões do arquivo `.env`: `chmod 600 .env` (Linux)
- ✓ Para Gmail, sempre use senha de app, nunca a senha principal

## 📞 Suporte

Em caso de problemas:
1. Verifique os logs de execução
2. Teste a conexão SQL Server manualmente
3. Teste o envio de email separadamente
4. Verifique se todas as bibliotecas Python estão instaladas

---

**Versão:** 1.0  
**Criado:** Outubro 2025  
**Email de relatórios:** ti01@marply.com.br
