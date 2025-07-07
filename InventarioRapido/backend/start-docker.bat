@echo off
title Docker - Sistema de Inventário
echo 🐳 Iniciando Sistema com Docker
echo ================================

REM Verificar se Docker está instalado
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker não encontrado!
    echo 📥 Instale Docker Desktop em: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM Verificar se Docker está rodando
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker não está rodando!
    echo 🚀 Inicie o Docker Desktop e tente novamente
    pause
    exit /b 1
)

echo ✅ Docker encontrado e rodando

REM Navegar para o diretório do script
cd /d "%~dp0"

echo.
echo 🏗️  Construindo e iniciando containers...
echo 📊 SQL Server será criado automaticamente
echo 📡 API será compilada e iniciada
echo.

REM Usar arquivo específico do Windows
docker-compose -f docker-compose-windows.yml up -d

if %errorlevel% neq 0 (
    echo ❌ Erro ao iniciar containers
    pause
    exit /b 1
)

echo.
echo 🎉 Containers iniciados com sucesso!
echo.
echo ⏳ Aguardando SQL Server inicializar (pode demorar alguns minutos)...
timeout /t 10 /nobreak >nul

echo.
echo 🔗 Serviços disponíveis:
echo    📡 API: http://localhost:8080
echo    📊 Health Check: http://localhost:8080/health
echo    🗄️  SQL Server: localhost:1433
echo.
echo 📋 Credenciais do SQL Server:
echo    Usuário: sa
echo    Senha: YourStrongPassword123!
echo    Banco: InventoryDB
echo.
echo 🛠️  Comandos úteis:
echo    Ver logs: docker-compose -f docker-compose-windows.yml logs
echo    Parar: docker-compose -f docker-compose-windows.yml down
echo    Reiniciar: docker-compose -f docker-compose-windows.yml restart

pause