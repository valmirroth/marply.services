# BI Análise Banco de Horas

Serviço em Go que processa o banco de horas dos colaboradores (ERP Senior/Vetorh), grava o resultado em SQL Server e oferece uma interface web protegida por senha para consulta e análise.

## O que faz

1. **ETL (`POST /run`)**: lê as apurações da origem (SQL Server do Vetorh ou CSV), abate o banco de horas negativo (`codsit 230`) das horas positivas, da mais antiga para a mais recente, e grava:
   - o **detalhado**, linha a linha, em `dbo.BH_DetalhadoPosAbatimento`;
   - o **resumo mensal** por colaborador, em `dbo.BH_ResumoMensalSaldo` (via `MERGE`, idempotente por `numemp, numcad, perref`).
2. **Interface (`/ui`)**: consulta o resumo com KPIs, filtros, ordenação, drill-down no detalhado e exportação CSV.
3. **Agendador**: dispara o `/run` automaticamente todo dia às **03:00** (fuso BRT, UTC-3).

## Estrutura

```
main.go                      rotas HTTP, login, UI embutida, agendador (cron)
internal/config/             leitura das variáveis de ambiente (.env)
internal/model/              tipos RawRow / ProcessedRow
internal/repository/         qry.go (queries de origem), source.go (leitura), dest.go (gravação/limpeza)
internal/service/            processor.go (regra de abatimento), summary.go (resumo mensal)
internal/util/               leitura de CSV
Dockerfile, docker-compose.yml
```

## Rotas

| Rota | Método | Autenticação | Descrição |
|---|---|---|---|
| `/health` | GET | não | healthcheck (`ok`) |
| `/login`, `/logout` | GET/POST, GET | não | tela de senha e encerramento de sessão |
| `/ui` | GET | sim | interface web |
| `/run` | POST | sim | executa o ETL completo |
| `/summary` | GET | sim | resumo mensal em JSON |
| `/periods` | GET | sim | períodos (`YYYYMM`) existentes no resumo |
| `/detail` | GET | sim | linhas do detalhado de um colaborador/período |

Filtros do `/summary` (todos opcionais e combináveis):

- `perref=YYYYMM` ou `YYYY-MM`, **repetível** (`?perref=202608&perref=202609`); sem nenhum, traz todos os períodos
- `numemp` (inteiro): ver nota abaixo
- `codccu`: prefixo do centro de custo
- `colab`: parte do nome do colaborador

`/detail` exige `numemp`, `numcad` e `perref`.

## Autenticação

Defina `UI_PASSWORD` no `.env`. Quem acessa `/ui` é redirecionado para `/login`; após informar a senha, recebe um cookie de sessão (`HttpOnly`) válido por **12 horas**.

- Reiniciar o serviço **derruba todas as sessões** (o token é gerado a cada start).
- Para scripts e para o agendador, envie a senha no header `X-BH-Password`:
  ```bash
  curl -X POST -H "X-BH-Password: <senha>" http://localhost:8093/run
  ```
- Se `UI_PASSWORD` ficar vazio, o serviço sobe **sem proteção** e registra um aviso no log.
- O repositório é público: nunca coloque a senha real no código nem no `.env.example`.

## Configuração (`.env`)

Copie `.env.example` para `.env` e preencha.

| Variável | Padrão | Descrição |
|---|---|---|
| `DEST_SQL_CONNECTION` | (obrigatória) | SQL Server de destino |
| `SRC_SQL_CONNECTION` | | SQL Server de origem (obrigatória se `USE_SQL_ORIGIN=true`) |
| `USE_SQL_ORIGIN` | `true` | `false` lê do CSV em `INPUT_CSV` |
| `INPUT_CSV`, `CSV_SEP` | `./resultado_query.csv`, `;` | usados só com origem CSV |
| `TBL_DETALHADO` | `dbo.BH_DetalhadoPosAbatimento` | tabela do detalhado |
| `TBL_RESUMO` | `dbo.BH_ResumoMensalSaldo` | tabela do resumo |
| `CLEAR_DEST` | `true` | apaga o destino (nas faixas fixas de `dest.go`) antes de gravar |
| `UPSERT_RESUMO` | `true` | grava o resumo via `MERGE` |
| `HTTP_PORT` | `:8093` | porta HTTP |
| `UI_PASSWORD` | (vazio) | senha da interface e da API |

Formato da string de conexão: `sqlserver://usuario:senha@host:1433?database=Banco&encrypt=disable`.

**Atenção com `$` na senha:** o Docker Compose interpreta `$` como variável. Escreva `$$` no `.env` (ex.: `Abc$xyz` vira `Abc$$xyz`), senão a senha é truncada sem aviso.

## Como executar

Direto na máquina:

```bash
go run .
```

Com Docker:

```bash
docker compose up -d --build
```

A interface fica em `http://localhost:8093/ui`.

### Docker Desktop no macOS e bancos na rede local

O Docker Desktop no Mac não alcança hosts da mesma sub-rede física do Mac (por exemplo `192.168.1.x`), embora a internet funcione. Em um servidor Linux isso não acontece. No Mac, a solução é fazer o próprio Mac servir de ponte, com `socat`, e apontar o `.env` para `host.docker.internal`:

```bash
brew install socat
socat TCP-LISTEN:14330,bind=127.0.0.1,fork,reuseaddr TCP:192.168.1.28:1433 &
socat TCP-LISTEN:14331,bind=127.0.0.1,fork,reuseaddr TCP:186.250.94.237:1433 &
```

```
DEST_SQL_CONNECTION=sqlserver://usuario:senha@host.docker.internal:14330?database=...
SRC_SQL_CONNECTION=sqlserver://usuario:senha@host.docker.internal:14331?database=...
```

Os `socat` não sobrevivem a reinício do Mac nem ao fechamento do terminal.

## Atenção ao alterar o ETL

- As **janelas de data** das queries ficam fixas em `internal/repository/qry.go` (`>= '20260926'` e `< '20270226'`, por exemplo) e precisam ser ajustadas a cada novo período.
- A limpeza do destino em `dest.go` (`DeleteAll`, `DeleteAllBhMonth`) também usa datas fixas.
- O agendador chama `http://192.168.1.28:8093/run` (endereço fixo em `triggerRunOnce`, em `main.go`). Se o serviço rodar em outro endereço, ajuste esse ponto.
- `InsertDetalhado` sempre insere; não há modo de simulação. Com `CLEAR_DEST=false`, rodar `/run` de novo duplica linhas no detalhado.

## Notas sobre os dados

- Pela ordem das colunas na query de origem e no `Scan` de `source.go`, o campo chamado `numemp` no código guarda o **crachá** (número de cadastro) e `numcad` guarda a **empresa**. A interface já exibe as colunas corretamente ("Empresa" e "Crachá/Núm. Cad."), mas no banco e no filtro `numemp` o significado é o do crachá.
- O detalhado só guarda os últimos meses, enquanto o resumo tem histórico maior. Por isso colaboradores antigos, sem linhas no detalhado, aparecem sem nome e sem CCU na interface.

## Composição do resumo mensal

- `horas_positivas_original`: soma das horas originais com `codsit != 230` no mês
- `banco_230_consumido_no_mes`: soma do banco usado nas linhas do mês
- `horas_saldo_mes` e `valor_saldo_mes`: soma do saldo e do valor após o abatimento
- `banco_total_aplicado_no_grupo`: banco total do colaborador (valor máximo do grupo, constante entre os meses; não some entre períodos)
