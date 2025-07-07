@echo off
echo 🗄️  Criando Banco de Dados - SQL Server
echo =====================================

set /p SERVER="Digite o servidor SQL Server (localhost): "
if "%SERVER%"=="" set SERVER=localhost

set /p DATABASE="Digite o nome do banco (InventoryDB): "
if "%DATABASE%"=="" set DATABASE=InventoryDB

set /p USERNAME="Digite o usuário (sa): "
if "%USERNAME%"=="" set USERNAME=sa

set /p PASSWORD="Digite a senha: "

echo.
echo 🔧 Conectando em %SERVER% com usuário %USERNAME%...

REM Criar o banco se não existir
sqlcmd -S %SERVER% -U %USERNAME% -P %PASSWORD% -Q "IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = '%DATABASE%') CREATE DATABASE [%DATABASE%]"

if %errorlevel% neq 0 (
    echo ❌ Erro ao conectar no SQL Server
    echo 🔍 Verifique se:
    echo    - SQL Server está rodando
    echo    - Credenciais estão corretas
    echo    - Usuário tem permissões
    pause
    exit /b 1
)

echo ✅ Banco %DATABASE% verificado/criado com sucesso!

REM Executar script de criação da tabela
echo.
echo 📊 Criando tabela CST_CONTAGEM_LOCAL...
sqlcmd -S %SERVER% -U %USERNAME% -P %PASSWORD% -d %DATABASE% -i "database\create_table.sql"

if %errorlevel% neq 0 (
    echo ❌ Erro ao criar tabela
    pause
    exit /b 1
)

echo.
echo 🎉 Banco de dados configurado com sucesso!
echo.
echo 📝 Atualize o arquivo .env com estas configurações:
echo DB_SERVER=%SERVER%
echo DB_NAME=%DATABASE%
echo DB_USER=%USERNAME%
echo DB_PASSWORD=%PASSWORD%

pause