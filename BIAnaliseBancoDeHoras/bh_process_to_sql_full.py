import pandas as pd
import pyodbc
import math
from typing import List, Optional, Tuple
import logging
from decimal import Decimal, InvalidOperation, ROUND_HALF_UP
# ======================
# CONFIGURAÇÕES
# ======================

# Se True, lê da origem via SQL Server; se False, usa CSV
USE_SQL_ORIGIN = True

# Conexão da ORIGEM (onde a query é executada)
SRC_CONN = (
      "DRIVER={ODBC Driver 17 for SQL Server};"
        "SERVER=186.250.94.237;"
        "DATABASE=vetorhprod;"
        "Trusted_Connection=no;"
        "UID=marinip;PWD=ZdnAin3Lp1$n7xh;"
)

# Query de ORIGEM (a mesma que você forneceu)
SRC_QUERY = r"""
select    codcal, perref, codccu, dtApuracao, codfil, resumo.CRACHA, resumo.Colaborador, resumo.numcad, resumo.numemp, resumo.tipcol, resumo.ValHoraMes, ValHoraCalculado, resumo.dessit, resumo.codsit, SUM(qtdhor) as horas,
    (ValHoraCalculado * SUM( resumo.qtdhor ) ) as [R$ valor]
from (
SELECT r044cal.codcal, r044cal.perref, r034fun.codccu, convert(varchar,(r066sit.datapu),112) dtApuracao,
     r034fun.codfil, getdate() as DataAtual, 
     cast(r034fun.numcad  as varchar) +' - '+ r034fun.nomfun as Colaborador,
     CAST(r034fun.numcra AS VARCHAR) AS CRACHA,               
      r066apu.numemp,              
      r066apu.numcad,               
      r066apu.datapu,               
      r066apu.tipcol,              
      nomesc AS TURNO ,  
      r034fun.valsal / 
	  cast(r034fun.valsal / 
	  (hormes/60) as decimal(19,6))as ValHoraMes,

	 cast( case when sit.codsit = 331 then ( r034fun.valsal / 
	  (hormes/60))*1.6
	  when sit.codsit = 330 then ( r034fun.valsal / 
	  (hormes/60))*1.5
	  when sit.codsit = 332 then ( r034fun.valsal / 
	  (hormes/60))*2
	  when sit.codsit in (230,109) then ( r034fun.valsal / 
	  (hormes/60))	  
	   when sit.codsit = 313 then ( r034fun.valsal / 
	  (hormes/60))*1.5
	   when sit.codsit = 302 then (( r034fun.valsal / 
	  (hormes/60))*1.2) * 1.5
	   when sit.codsit = 301 then (( r034fun.valsal / 
	  (hormes/60))*1.5)  
	     when sit.codsit = 305 then (( r034fun.valsal / 
	  (hormes/60))*2) 
	       when sit.codsit = 311 then (( r034fun.valsal / 
	  (hormes/60))*1.6) 
	       when sit.codsit = 306 then ((( r034fun.valsal / 
	  (hormes/60))*1.2) * 2)
	       when sit.codsit = 314 then ((( r034fun.valsal / 
	  (hormes/60))*1.2) * 1.5)
	       when sit.codsit = 312 then ((( r034fun.valsal / 
	  (hormes/60))*1.2) * 1.6)

	  end as decimal(19,6)) as ValHoraCalculado,
	case when sit.codsit in (109,230) then 'Banco de Horas Negativo' else 
    replace(
    replace(
    replace(   
    replace(  
    replace(replace(replace(sit.dessit,'Banco de Horas','B.H.'),'H. Extra','H.E.'),'Horas Extras','H.E.') 
    ,'Intrajornada','Intr')
    ,'Noturnas','Not')
    ,'Intraj','Intr')
    ,'Apurad','Apu.')
    end as dessit,    
    case when sit.codsit in (230,109) then 230 else sit.codsit end as codsit,
    case when sit.codsit in (230,109) then (r066sit.qtdhor/60.000) * (-1) else (r066sit.qtdhor/60.000)  end qtdhor 
FROM r066sit
    left join r066apu on  r066sit.numemp = r066apu.numemp and r066sit.tipcol = r066apu.tipcol and r066sit.numcad = r066apu.numcad           
      and r066sit.datapu = r066apu.datapu
    left join r034fun (nolock) on r034fun.numcad = r066apu.numcad   
    LEFT JOIN r010sit sit on sit.codsit = r066sit.codsit
    INNER JOIN r006esc (NOLOCK) ON r066apu.codesc = r006esc.codesc
    inner join r010sit on r010sit.codsit = r034fun.sitafa
    inner join r044cal on r044cal.tipcal ='11' and  r066apu.datapu between r044cal.iniapu and r044cal.fimapu and r044cal.numemp = r066apu.numemp 
where 
    r010sit.codsit <> '7'
    --and r034fun.codccu in ('cc103','cc789','cc252')
    and convert(varchar,( eomonth( r044cal.fimapu )),112) >= convert(varchar,dateadd(day,1,eomonth(dateadd(month,-6,getdate()))),112)  
    and convert(varchar,(r066sit.datapu),112) < convert(varchar,(getdate()),112)
    and (  sit.dessit like '%banco%' )
    and codfil in (1,5)
 ) as resumo
 where CRACHA <> 6065 -- saiu do setor e não tem mais banco
   -- and codccu  = 'CC196'
  --  and codfil = 1
   	and numemp = 1
  --  and CRACHA = '3732'
 group by codcal, perref, codccu, dtApuracao, codfil, resumo.CRACHA, ValHoraCalculado, Colaborador, resumo.numcad, resumo.numemp, resumo.tipcol, resumo.ValHoraMes, resumo.dessit, resumo.codsit
"""

# CSV de fallback (se USE_SQL_ORIGIN=False)
INPUT_CSV = r"./resultado_query.csv"
CSV_SEP = ";"  # mude se necessário

# Grupo para aplicar o banco
GROUP_COLS: List[str] = ["numemp", "numcad"]

# Conexão do DESTINO (onde serão gravadas as tabelas)
DEST_CONN = (
    "DRIVER={ODBC Driver 17 for SQL Server};"
    "SERVER=192.168.1.28,1433;"
    "DATABASE=Marini_PRD;"
    "UID=Octopus;PWD=A45182008069199;"
    "TrustServerCertificate=no;"
)

# Tabelas no destino
TBL_DETALHADO = "dbo.BH_DetalhadoPosAbatimento"
TBL_RESUMO    = "dbo.BH_ResumoMensalSaldo"

# Se True, usa MERGE (upsert) no resumo
UPSERT_RESUMO = True


# ======================
# LEITURA ORIGEM
# ======================

def read_from_sql(conn_str: str, query: str) -> pd.DataFrame:
    logging.info("Iniciando conexão com servidor Senior")
    with pyodbc.connect(conn_str) as con:
        df = pd.read_sql(query, con)
    return df

def read_from_csv(path: str, sep: str = ";") -> pd.DataFrame:
    return pd.read_csv(path, sep=sep, dtype=str)

def normalize_columns(df: pd.DataFrame) -> pd.DataFrame:
    # forçar nomes minúsculos e sem espaços
    logging.info("forçar nomes minúsculos e sem espaços")
    df = df.copy()
    df.columns = [c.strip().lower() for c in df.columns]

    # datas
    logging.info("forçar formato de datas")
    if "dtapuracao" in df.columns:
        # aceitar yyyyMMdd, yyyy-MM-dd, etc.
        df["dtapuracao"] = pd.to_datetime(df["dtapuracao"], errors="coerce", format=None)

    # perref: se não existir, derive de dtapuracao
    if "perref" not in df.columns:
        if "dtapuracao" in df.columns:
            df["perref"] = df["dtapuracao"].dt.strftime("%Y%m")
        else:
            df["perref"] = pd.NA

    logging.info("forçar formato de númericos")
    # numericos
    for col in ["horas","valhoracalculado","valhorames","codsit","numemp","numcad","codfil","codcal","tipcol"]:
        if col in df.columns:
            # codsit/num* podem ter string; convertemos com coerce
            if col in ["codsit","numemp","numcad","codfil","codcal","tipcol"]:
                df[col] = pd.to_numeric(df[col], errors="coerce")
            else:
                df[col] = pd.to_numeric(df[col], errors="coerce").fillna(0.0)

    # valor original (se existir com nome "r$ valor")
    if "r$ valor" in df.columns and "valor_reais" not in df.columns:
        df = df.rename(columns={"r$ valor":"valor_reais"})

    return df

def load_data() -> pd.DataFrame:
    if USE_SQL_ORIGIN:
        df = read_from_sql(SRC_CONN, SRC_QUERY)
    else:
        df = read_from_csv(INPUT_CSV, sep=CSV_SEP)
    df = normalize_columns(df)
    return df


# ==============================================
# APLICAÇÃO DO BANCO 230 (mais antigo → atual)
# ==============================================

def apply_bank_offset_for_group(gdf: pd.DataFrame) -> pd.DataFrame:
    gdf = gdf.copy()
    gdf["horas_original"] = gdf["horas"]
    logging.info("Realizando calculo do saldos do banco de horas situação 230 e 109....")

    # total de banco (soma de 230, que vem negativo na query). Transformamos em positivo para abater.
    banco_total_horas = -gdf.loc[gdf["codsit"] == 230, "horas"].sum()
    banco_restante = float(max(banco_total_horas, 0.0))

    # separa positivas (não-230) e ordena da mais antiga p/ mais recente
    pos_mask = gdf["codsit"] != 230
    pos = gdf.loc[pos_mask].sort_values(["dtapuracao", "perref"], na_position="last").copy()
    pos["banco_usado_na_linha"] = 0.0
    pos["horas_saldo"] = 0.0

    for idx, row in pos.iterrows():
        horas_linha = float(row.get("horas", 0.0) or 0.0)
        if horas_linha <= 0:
            # nada para abater desta linha
            pos.at[idx, "horas_saldo"] = 0.0
            continue

        if banco_restante <= 0:
            # não há mais banco para abater -> saldo é a própria hora positiva
            pos.at[idx, "horas_saldo"] = horas_linha
            continue

        abatimento = min(horas_linha, banco_restante)
        saldo = horas_linha - abatimento
        banco_restante -= abatimento

        pos.at[idx, "horas_saldo"] = saldo
        pos.at[idx, "banco_usado_na_linha"] = abatimento

    # linhas de banco (230)
    bank = gdf.loc[~pos_mask].copy()
    if not bank.empty:
        bank["banco_usado_na_linha"] = 0.0
        bank["horas_saldo"] = 0.0

        # Se sobrou banco após abater tudo (banco_restante > 0),
        # expõe o saldo negativo em horas_saldo na linha 230 mais antiga.
        if banco_restante > 0 and not bank.empty:
            sel_idx = bank.sort_values(
                ["dtapuracao", "perref"], ascending=[False, False], na_position="last"
            ).index[0]
            bank.loc[sel_idx, "horas_saldo"] = -banco_restante  # saldo negativo de horas
    # Consolida e calcula valor_saldo
    out = pd.concat([pos, bank], ignore_index=True)

    if "valhoracalculado" in out.columns:
        out["valor_saldo"] = out["horas_saldo"] * out["valhoracalculado"].fillna(0.0)
    else:
        out["valor_saldo"] = 0.0

    out["banco_total_aplicado_no_grupo"] = banco_total_horas
    return out


def apply_bank_offset(df: pd.DataFrame, group_cols: List[str]) -> pd.DataFrame:
    needed = set(group_cols + ["codsit","horas","dtapuracao","perref"])
    missing = needed - set(df.columns)
    if missing:
        raise ValueError(f"Faltam colunas para processamento: {missing}")
    processed = (
        df.groupby(group_cols, dropna=False, group_keys=True, sort=False)
          .apply(apply_bank_offset_for_group)
          .reset_index(drop=True)
    )
    return processed


# =========================================
# RESUMO MENSAL
# =========================================

def build_monthly_summary(processed: pd.DataFrame, group_cols: List[str]) -> pd.DataFrame:
    df = processed.copy()
    effective = df[df["codsit"] != 230].copy()
    agg_map = {
        "horas_original":"sum",
        "banco_usado_na_linha":"sum",
        "horas_saldo":"sum",
        "valor_saldo":"sum"
    }
    summary = (
        effective.groupby(group_cols+["perref"], dropna=False, as_index=False)
        .agg(agg_map)
        .rename(columns={
            "horas_original":"horas_positivas_original",
            "banco_usado_na_linha":"banco_230_consumido_no_mes",
            "horas_saldo":"horas_saldo_mes",
            "valor_saldo":"valor_saldo_mes"
        })
    )
    banco_total = (
        df.groupby(group_cols, dropna=False, as_index=False)["banco_total_aplicado_no_grupo"].max()
    )
    summary = summary.merge(banco_total, on=group_cols, how="left")
    return summary


# ======================
# EXPORTAÇÃO PARA SQL
# ======================

def ensure_types_for_sql(df: pd.DataFrame) -> pd.DataFrame:
    out = df.copy()
    if "perref" in out.columns:
        out["perref"] = out["perref"].astype(str).str[:10]
    if "dtapuracao" in out.columns and pd.api.types.is_datetime64_any_dtype(out["dtapuracao"]):
        out["dtapuracao"] = out["dtapuracao"].dt.strftime("%Y-%m-%d")
    return out
import re
# --- util: normalizar números em string para float/decimal ---
_num_pt_re = re.compile(r"^\s*-?(\d{1,3}(\.\d{3})+|\d+)(,\d+)?\s*$")   # "1.234,56" ou "1234,56"
_num_us_re = re.compile(r"^\s*-?\d+(\.\d+)?\s*$")                      # "1234.56"
def parse_number_str_to_float(s: str):
    if s is None:
        return None
    s = str(s).strip()
    if s == "":
        return None
    # 1) pt-BR "1.234,56"
    if _num_pt_re.match(s):
        s = s.replace('.', '').replace(',', '.')
        return float(s)
    # 2) US "1234.56" (ou "2.5e-3")
    if _num_us_re.match(s) or 'e' in s.lower():
        return float(s)
    # 3) "2,9744820918" (sem milhar, só vírgula decimal)
    if s.count(',') == 1 and s.count('.') == 0 and all(ch.isdigit() or ch in ',.- ' for ch in s):
        return float(s.replace(',', '.'))
    # 4) números com milhar US "1,234.56" (menos comum em dados pt-BR)
    if s.count(',') > 0 and s.count('.') == 1 and s.replace(',', '').replace('.', '').replace('-', '').isdigit():
        return float(s.replace(',', ''))
    # último recurso: tenta Decimal -> float
    try:
        return float(Decimal(s.replace(',', '.')))
    except Exception:
        raise ValueError(f"Texto numérico inválido para float: '{s}'")

def _read_table_schema(cn, table):
    if '.' in table:
        schema, name = table.split('.', 1)
    else:
        schema, name = 'dbo', table
    rows = cn.cursor().execute("""
        SELECT COLUMN_NAME, DATA_TYPE,
               CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, IS_NULLABLE
        FROM INFORMATION_SCHEMA.COLUMNS
        WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
    """, schema, name).fetchall()
    return {
        r.COLUMN_NAME.lower(): {
            "type": r.DATA_TYPE.lower(),
            "maxlen": r.CHARACTER_MAXIMUM_LENGTH,
            "prec": r.NUMERIC_PRECISION,
            "scale": r.NUMERIC_SCALE,
            "nullable": (r.IS_NULLABLE == 'YES')
        } for r in rows
    }

def _to_decimal_like(val):
    # aceita float, int, str com vírgula decimal etc.
    if val is None:
        return None
    if isinstance(val, Decimal):
        return val
    s = str(val).strip()
    if s == "":
        return None
    # normaliza pt-BR "1.234,56" ou "0,016666"
    if any(c in s for c in ",.") and "," in s and "." in s:
        s = s.replace('.', '').replace(',', '.')
    elif "," in s and "." not in s:
        s = s.replace(',', '.')
    try:
        return Decimal(s)
    except Exception:
        # último recurso: via float (pode introduzir binário, mas passa)
        return Decimal(str(float(val)))

def _quantize_to_scale(d: Decimal, scale: int) -> Decimal:
    # quantiza para a escala do schema, com arredondamento HALF_UP
    if scale and scale > 0:
        return d.quantize(Decimal(1).scaleb(-scale), rounding=ROUND_HALF_UP)
    return d.quantize(Decimal(1), rounding=ROUND_HALF_UP)

def _fit_decimal_or_raise(d: Decimal, prec: int, scale: int, colname: str) -> Decimal:
    # 1) quantiza para a escala
    q = _quantize_to_scale(d, scale)
    # 2) representação sem expoente
    txt = format(q, 'f')   # ex.: "-12345.678900" ou "0.016666000000000000"
    # 3) separa partes
    if '.' in txt:
        int_part, frac_part = txt.split('.', 1)
    else:
        int_part, frac_part = txt, ""
    # remove sinal e zeros à esquerda
    ip = int_part.lstrip('-').lstrip('0')
    int_digits = len(ip)  # valores entre -1 e 1 resultam em 0
    if int_digits > (prec - scale):
        raise ValueError(
            f"{colname}: excede DECIMAL({prec},{scale}) "
            f"(parte inteira tem {int_digits} dígitos; máximo é {prec - scale})."
        )
    # frac_part já está com exatamente 'scale' casas devido ao quantize
    return q

def _coerce_and_validate_row(row, cols, schema):
    out = []
    for col, val in zip(cols, row):
        meta = schema.get(col.lower(), {})
        t = meta.get("type")

        # normaliza nulos
        if pd.isna(val):
            val = None

        if t in ("nvarchar","varchar","nchar","char","text","ntext"):
            if val is not None:
                s = str(val)
                mx = meta.get("maxlen")
                if isinstance(mx, int) and mx > 0 and len(s) > mx:
                    raise ValueError(f"{col}: excede tamanho ({len(s)}>{mx})")
                val = s

        elif t in ("decimal","numeric","money","smallmoney"):
            if val is not None:
                d = _to_decimal_like(val)
                if d is None:
                    val = None
                else:
                    prec = meta.get("prec") or 38
                    scale = meta.get("scale") or 0
                    try:
                        q = _fit_decimal_or_raise(d, prec, scale, col)
                    except ValueError as ve:
                        raise ValueError(f"{col}: {ve} (valor: {val!r})")
                    val = q
        elif t in ("float","real"):
            if val is not None:
                if isinstance(val, str):
                    val = parse_number_str_to_float(val)
                else:
                    try:
                        val = float(val)
                    except Exception:
                        raise ValueError(f"{col}: valor float inválido: {val}")
                if math.isnan(val) or math.isinf(val):
                    raise ValueError(f"{col}: NaN/Inf não permitido")

        elif t in ("date","datetime","smalldatetime","datetime2","time"):
            if val is not None:
                if isinstance(val, str):
                    try:
                        val = pd.to_datetime(val, dayfirst=False, errors='raise')
                    except Exception as e:
                        # tenta pt-BR dd/mm/yyyy
                        try:
                            val = pd.to_datetime(val, dayfirst=True, errors='raise')
                        except Exception:
                            raise ValueError(f"{col}: data inválida '{val}': {e}")
                if t == "date" and hasattr(val, "date"):
                    val = val.date()

        else:
            # inteiros
            if t in ("int","bigint","smallint","tinyint"):
                if val is not None:
                    if isinstance(val, str):
                        s = val.strip()
                        if s == "":
                            val = None
                        else:
                            try:
                                s = s.replace('.', '').replace(',', '')  # remove milhar
                                val = int(s)
                            except Exception:
                                raise ValueError(f"{col}: inteiro inválido: {val}")
                    else:
                        try:
                            val = int(val)
                        except Exception:
                            raise ValueError(f"{col}: inteiro inválido: {val}")

        if val is None and not meta.get("nullable", True):
            raise ValueError(f"{col}: NULL não permitido")
        out.append(val)
    return tuple(out)

def bulk_insert_pyodbc(df: pd.DataFrame, table: str, conn_str: str, cols_in_order: list, batch_size: int = 500):
    if df.empty:
        print(f"[{table}] DataFrame vazio — nada a inserir.")
        return

    df = df[cols_in_order].copy()
    df = df.where(pd.notnull(df), None)

    placeholders = ",".join(["?"]*len(cols_in_order))
    sql = f"INSERT INTO {table} ({','.join(cols_in_order)}) VALUES ({placeholders})"

    with pyodbc.connect(conn_str) as cn:
        cn.autocommit = False
        schema = _read_table_schema(cn, table)
        cur = cn.cursor()
        cur.fast_executemany = True

        # Pré-valida e normaliza cada linha
        rows = []
        for i, row in enumerate(df.itertuples(index=False, name=None), start=1):
            try:
                rows.append(_coerce_and_validate_row(row, cols_in_order, schema))
            except Exception as e:
                raise RuntimeError(f"[VALIDAÇÃO] Linha {i}: {e}")

        # Inserção em lotes com fallback linha a linha para localizar a falha exata
        try:
            for start in range(0, len(rows), batch_size):
                chunk = rows[start:start+batch_size]
                try:
                    cur.executemany(sql, chunk)
                except pyodbc.ProgrammingError as e:
                    # Captura "Parâmetro X" e tenta localizar
                    cn.rollback()
                    msg = str(e)
                    # fallback para identificar a linha e a coluna
                    for j, one in enumerate(chunk, start=start+1):
                        try:
                            cur.execute(sql, one)
                        except pyodbc.ProgrammingError as e2:
                            msg2 = str(e2)
                            # tenta extrair "Parâmetro NN"
                            m = re.search(r"Parâmetro\s+(\d+)", msg2)
                            if m:
                                p = int(m.group(1))
                                # mapeia para coluna (1-based -> 0-based)
                                col_idx = p-1 if 1 <= p <= len(cols_in_order) else None
                                col_name = cols_in_order[col_idx] if col_idx is not None else "desconhecida"
                                val = one[col_idx] if col_idx is not None else None
                                raise RuntimeError(
                                    f"[INSERÇÃO] Falha na linha {j}, coluna {col_name} (parâmetro {p}). "
                                    f"Valor problemático: {val!r}. Erro do driver: {msg2}"
                                ) from e2
                            else:
                                raise RuntimeError(
                                    f"[INSERÇÃO] Falha na linha {j}. Valores: {dict(zip(cols_in_order, one))}. "
                                    f"Erro do driver: {msg2}"
                                ) from e2
                    # se chegou aqui, relança o erro original
                    raise
            cn.commit()
        except Exception:
            cn.rollback()
            raise

    print(f"[{table}] Inseridos {len(df)} registros.")



def delete_bh_data_all(conn_str: str) -> Tuple[int, int]:
    """
    Deleta TODOS os registros de:
      - dbo.BH_DetalhadoPosAbatimento
      - dbo.BH_ResumoMensalSaldo
    Retorna (qtde_detalhado, qtde_resumo).
    """
    detail_tbl = "BH_DetalhadoPosAbatimento"
    resumo_tbl = "BH_ResumoMensalSaldo"

    with pyodbc.connect(conn_str) as cn:
        cn.autocommit = False
        cur = cn.cursor()
        try:
            logging.info("Deletando TODOS os registros de %s ...", detail_tbl)
            cur.execute(f"DELETE FROM {detail_tbl};")
            deleted_detail = cur.rowcount if cur.rowcount is not None else -1

            logging.info("Deletando TODOS os registros de %s ...", resumo_tbl)
            cur.execute(f"DELETE FROM {resumo_tbl};")
            deleted_resumo = cur.rowcount if cur.rowcount is not None else -1

            cn.commit()
            logging.info("Delete concluído. Detalhado=%s, Resumo=%s", deleted_detail, deleted_resumo)
            return deleted_detail, deleted_resumo
        except Exception as e:
            cn.rollback()
            logging.exception("Falha ao deletar dados (ALL).")
            raise
        
# ======================
# MAIN
# ======================

def JobBH():
    delete_bh_data_all(DEST_CONN)
    # 1) carrega origem
    df = load_data()

    # 2) processa abatimento
    processed = apply_bank_offset(df, GROUP_COLS)

    # 3) resumo
    resumo_mensal = build_monthly_summary(processed, GROUP_COLS)

    # 4) grava no destino
    det_cols = [
        "numemp","numcad","codccu","codfil","codcal","perref","dtapuracao","cracha","colaborador",
        "tipcol","dessit","codsit","valhorames","valhoracalculado","horas_original",
        "banco_usado_na_linha","horas_saldo","valor_saldo","valor_reais","banco_total_aplicado_no_grupo"
    ]
    res_cols = [
        "numemp","numcad","perref","horas_positivas_original","banco_230_consumido_no_mes",
        "horas_saldo_mes","valor_saldo_mes","banco_total_aplicado_no_grupo"
    ]

    # Detalhado geralmente é append
    bulk_insert_pyodbc(processed, TBL_DETALHADO, DEST_CONN, det_cols)

    # Resumo: upsert opcional
   
 