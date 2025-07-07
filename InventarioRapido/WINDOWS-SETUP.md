# 🪟 Sistema de Contagem de Inventário - Windows

## 🚀 Inicio Rápido

### Opção 1: Sistema Completo (Mais Fácil)
```cmd
start-full-system.bat
```
Inicia automaticamente frontend e backend juntos.

### Opção 2: Apenas Frontend
```cmd
start-frontend.bat
```
Inicia apenas o aplicativo React (porta 5000).

### Opção 3: Com Docker (Recomendado)
```cmd
cd backend
start-docker.bat
```
Inicia tudo com Docker: SQL Server + API automática.

## 📋 Configuração Detalhada

### 1. Frontend (Sistema Web)
- **Arquivo:** `start-frontend.bat`
- **Função:** Interface web para tablets/celulares
- **Porta:** 5000
- **URL:** http://localhost:5000

### 2. Backend (API Go + SQL Server)
- **Setup:** `backend\setup.bat`
- **Banco:** `backend\create-database.bat`
- **Iniciar:** `backend\start-api.bat`
- **Porta:** 8080

### 3. Docker (Tudo Automatizado)
- **Arquivo:** `backend\start-docker.bat`
- **Inclui:** SQL Server + API + Banco configurado
- **Portas:** 8080 (API) + 1433 (SQL Server)

## 🗄️ Configuração do Banco

### Automática (Docker)
```cmd
cd backend
start-docker.bat
```

### Manual (SQL Server Local)
1. Execute: `backend\create-database.bat`
2. Edite: `backend\.env` com suas credenciais
3. Execute: `backend\start-api.bat`

## 📱 Como Usar

1. **Inicie o sistema:** `start-full-system.bat`
2. **Abra no navegador:** http://localhost:5000
3. **Digite código do item**
4. **Preencha local, quantidade e volumes**
5. **Clique "Registrar Contagem"**
6. **Use "Encerrar Contagem" quando terminar**

## 🔧 Solução de Problemas

### Node.js não encontrado
- Baixe em: https://nodejs.org/

### Go não encontrado  
- Baixe em: https://golang.org/dl/

### SQL Server não conecta
- Verifique se está rodando
- Use `backend\create-database.bat`
- Edite credenciais em `backend\.env`

### Erro de porta ocupada
- Frontend: Feche outros servidores na porta 5000
- Backend: Feche outros serviços na porta 8080

## 📁 Estrutura de Arquivos

```
📂 Projeto/
├── 🚀 start-full-system.bat    # Inicia tudo
├── 📱 start-frontend.bat       # Só frontend  
├── 📂 backend/
│   ├── 🔧 setup.bat           # Instala Go
│   ├── 🗄️ create-database.bat # Cria banco
│   ├── 🚀 start-api.bat       # Inicia API
│   ├── 🐳 start-docker.bat    # Docker
│   └── 📄 .env                # Configurações
└── 📂 client/                 # Aplicativo web
```

## ✅ Verificação

Acesse para testar:
- 📱 **Frontend:** http://localhost:5000
- 📡 **API:** http://localhost:8080/health
- 🗄️ **Banco:** localhost:1433 (via SQL Management Studio)