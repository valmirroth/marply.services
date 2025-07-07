# API de Contagem de Inventário - Go

API REST desenvolvida em Go usando o padrão MVC para gerenciar contagens de inventário no SQL Server.

## Estrutura do Projeto

```
backend/
├── controllers/         # Controllers (lógica de negócio)
│   └── contagem_controller.go
├── models/             # Models (estruturas de dados e acesso ao banco)
│   ├── contagem.go
│   └── database.go
├── routes/             # Rotas da API
│   └── routes.go
├── database/           # Scripts SQL
│   └── create_table.sql
├── main.go            # Arquivo principal
├── go.mod             # Dependências Go
└── .env              # Configurações (não versionar em produção)
```

## Configuração

### Opção 1: Configuração Manual

#### 1. Instalar Go

Se Go não estiver instalado, execute:
```bash
./setup.sh
```

#### 2. Banco de Dados SQL Server

Execute o script `database/create_table.sql` no seu SQL Server para criar a tabela necessária.

#### 3. Variáveis de Ambiente

Configure o arquivo `.env` com suas credenciais do SQL Server:

```env
DB_SERVER=localhost
DB_PORT=1433
DB_NAME=InventoryDB
DB_USER=sa
DB_PASSWORD=YourPassword123
API_PORT=8080
```

#### 4. Instalar Dependências e Executar

```bash
cd backend
go mod tidy
go run main.go
```

### Opção 2: Usando Docker

Para uma configuração mais simples com SQL Server incluído:

```bash
cd backend
docker-compose up -d
```

Isso iniciará:
- API Go na porta 8080
- SQL Server na porta 1433
- Criará automaticamente a tabela necessária

## Endpoints da API

### Contagens de Inventário

- `GET /api/contagens` - Listar todas as contagens
- `POST /api/contagens` - Criar nova contagem
- `GET /api/contagens/{id}` - Buscar contagem por ID
- `PUT /api/contagens/{id}` - Atualizar contagem
- `DELETE /api/contagens/{id}` - Excluir contagem

### Health Check

- `GET /health` - Verificar status da API

## Exemplo de Uso

### Criar Contagem (POST /api/contagens)

```json
{
  "itemCode": "A001234",
  "itemDescription": "Parafuso Sextavado M12x50 Aço Inox",
  "location": "EST-A01",
  "quantity": "10.50000",
  "volumeCount": 5,
  "usuario_contagem": "Admin"
}
```

### Resposta

```json
{
  "success": true,
  "message": "Contagem criada com sucesso",
  "data": {
    "id": "uuid-gerado",
    "codigo_item": "A001234",
    "descricao_item": "Parafuso Sextavado M12x50 Aço Inox",
    "local": "EST-A01",
    "quantidade": 10.50000,
    "volumes": 5,
    "data_contagem": "2024-01-01T10:00:00Z",
    "usuario_contagem": "Admin",
    "status": "ATIVO"
  }
}
```

## Tecnologias Utilizadas

- **Go 1.22** - Linguagem de programação
- **Gorilla Mux** - Router HTTP
- **go-mssqldb** - Driver SQL Server
- **godotenv** - Gerenciamento de variáveis de ambiente

## Padrão MVC

- **Models**: Gerenciam dados e acesso ao banco de dados
- **Controllers**: Contêm a lógica de negócio e processam requisições HTTP
- **Routes**: Definem as rotas da API e conectam URLs aos controllers