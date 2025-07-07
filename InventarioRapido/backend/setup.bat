@echo off
echo 🚀 Configurando API de Contagem de Inventário em Go - Windows
echo ===========================================================

REM Verificar se Go está instalado
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go não encontrado. 
    echo 📥 Baixe e instale Go em: https://golang.org/dl/
    echo 📝 Após instalar, reinicie o terminal e execute este script novamente.
    pause
    exit /b 1
) else (
    echo ✅ Go encontrado: 
    go version
)

REM Navegar para o diretório do script
cd /d "%~dp0"

REM Instalar dependências
echo.
echo 📦 Instalando dependências Go...
go mod tidy

if %errorlevel% neq 0 (
    echo ❌ Erro ao instalar dependências
    pause
    exit /b 1
)

REM Compilar a aplicação
echo.
echo 🔨 Compilando a aplicação...
go build -o inventory-api.exe main.go

if %errorlevel% neq 0 (
    echo ❌ Erro ao compilar
    pause
    exit /b 1
)

echo.
echo 🎉 Configuração concluída com sucesso!
echo.
echo 📋 Próximos passos:
echo 1. Configure o SQL Server e execute o script database\create_table.sql
echo 2. Edite o arquivo .env com suas credenciais do banco
echo 3. Execute a API com: start-api.bat
echo.
echo 🔗 A API estará disponível em: http://localhost:8080
echo 📖 Documentação: http://localhost:8080/health

pause