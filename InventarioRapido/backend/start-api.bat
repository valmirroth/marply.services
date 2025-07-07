@echo off
title API de Contagem de Inventário
echo 🚀 Iniciando API de Contagem de Inventário...
echo ============================================

REM Navegar para o diretório do script
cd /d "%~dp0"

REM Verificar se o arquivo .env existe
if not exist ".env" (
    echo ❌ Arquivo .env não encontrado!
    echo 📝 Criando arquivo .env modelo...
    echo.
    echo # SQL Server Configuration > .env
    echo DB_SERVER=localhost >> .env
    echo DB_PORT=1433 >> .env
    echo DB_NAME=InventoryDB >> .env
    echo DB_USER=sa >> .env
    echo DB_PASSWORD=YourPassword123 >> .env
    echo. >> .env
    echo # API Configuration >> .env
    echo API_PORT=8080 >> .env
    echo.
    echo ✅ Arquivo .env criado! Edite-o com suas credenciais do SQL Server.
    echo.
    pause
)

REM Verificar se Go está instalado
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go não encontrado. Execute setup.bat primeiro.
    pause
    exit /b 1
)

echo 🔧 Configuração:
echo 📊 Banco: %DB_SERVER%:%DB_PORT%/%DB_NAME%
echo 🌐 API: http://localhost:%API_PORT%
echo.

echo 🏃 Iniciando servidor...
echo 🔗 Endpoints disponíveis:
echo    GET    /api/contagens - Listar contagens
echo    POST   /api/contagens - Criar contagem  
echo    GET    /api/contagens/{id} - Buscar contagem por ID
echo    PUT    /api/contagens/{id} - Atualizar contagem
echo    DELETE /api/contagens/{id} - Excluir contagem
echo    GET    /health - Status da API
echo.
echo 💡 Pressione Ctrl+C para parar o servidor
echo.

go run main.go

if %errorlevel% neq 0 (
    echo.
    echo ❌ Erro ao iniciar API
    echo 🔍 Verifique:
    echo    - Credenciais do SQL Server no arquivo .env
    echo    - Se o SQL Server está rodando
    echo    - Se a tabela CST_CONTAGEM_LOCAL foi criada
    pause
)