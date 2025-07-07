# BH Processor (Go MVC) — v2 (Resumo + UI)

Agora com:
- **Resumo mensal** calculado e gravado em `dbo.BH_ResumoMensalSaldo` (com **UPSERT** opcional).
- **API** para consulta do resumo: `GET /summary?perref=YYYYMM&numemp=1&codccu=CC123&colab=ana`
- **UI** mínima (HTML/JS) em `GET /ui` com filtros (mês/ano, empresa, centro de custo, nome do colaborador).

## Rotas

- `GET /health` — simples healthcheck.
- `POST /run` — executa o ETL completo (limpa destino conforme `.env`, lê, processa, insere detalhado e resumo).
- `GET /summary` — consulta o resumo mensal (JSON). Filtros:
  - `perref=YYYYMM` (opcional; aceita também `YYYY-MM`)
  - `numemp` (int, opcional)
  - `codccu` (string, opcional; aceita prefixo exato)
  - `colab` (string, opcional; faz `LIKE %%colab%%` no `Colaborador` do detalhado agregado)
- `GET /ui` — página HTML com formulário e tabela dinâmica.

## Executando

```bash
go mod tidy
go run ./cmd/app    # inicia o servidor http em :8080
```

Para processar os dados, faça um `POST /run` (ou execute a main sem servidor, se preferir adaptar). A UI usa `/summary` apenas para leitura.

## Observações

- O **Resumo** é gerado a partir do `processed` (mesmo cálculo do Python comentado), agregando:
  - `horas_positivas_original` (sum das horas originais onde `codsit != 230` no mês)
  - `banco_230_consumido_no_mes` (sum dos `banco_usado_na_linha` no mês)
  - `horas_saldo_mes` (sum dos `horas_saldo` no mês)
  - `valor_saldo_mes` (sum dos `valor_saldo` no mês)
  - `banco_total_aplicado_no_grupo` (max dentro do grupo)
- O UPSERT usa `MERGE` por linha (idempotente para `(numemp,numcad,perref)`).
