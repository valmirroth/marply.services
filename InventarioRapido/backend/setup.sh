#!/bin/bash

echo "🚀 Configurando API de Contagem de Inventário em Go"
echo "=================================================="

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go não encontrado. Instalando Go..."
    
    # Baixar e instalar Go (Linux)
    GO_VERSION="1.22.0"
    wget -c https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    
    # Adicionar Go ao PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    
    # Limpar arquivo de download
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    echo "✅ Go instalado com sucesso!"
else
    echo "✅ Go já está instalado: $(go version)"
fi

# Navegar para o diretório backend
cd "$(dirname "$0")"

# Instalar dependências
echo "📦 Instalando dependências Go..."
go mod tidy

# Compilar a aplicação
echo "🔨 Compilando a aplicação..."
go build -o inventory-api main.go

echo ""
echo "🎉 Configuração concluída com sucesso!"
echo ""
echo "📋 Próximos passos:"
echo "1. Configure o SQL Server e execute o script database/create_table.sql"
echo "2. Edite o arquivo .env com suas credenciais do banco"
echo "3. Execute a API com: ./inventory-api"
echo ""
echo "🔗 A API estará disponível em: http://localhost:8080"
echo "📖 Documentação: http://localhost:8080/health"