# ==========================================================
# ARQUIVO 4: data_processor.py
# Processamento de dados e cálculos
# ==========================================================

import pandas as pd
from database import QUERY_TEMPLATE, get_conn

def classify_risk(media: float, ultimo: float) -> str:
    """Classifica risco baseado em média e último valor"""
    if media <= 0:
        return "BAIXO"
    if ultimo >= media * 1.5 and media >= 1000:
        return "ALTO"
    if ultimo >= media * 1.2:
        return "MÉDIO"
    return "BAIXO"

def parse_periodo_to_datetime(series: pd.Series) -> pd.Series:
    """Converte período para datetime"""
    s = series.astype(str).str.strip()
    is_ymd8 = s.str.fullmatch(r"\d{8}", na=False)
    dt = pd.to_datetime(pd.Series([pd.NaT] * len(s)))

    if is_ymd8.any():
        dt.loc[is_ymd8] = pd.to_datetime(s[is_ymd8], format="%Y%m%d", errors="coerce")

    if (~is_ymd8).any():
        dt.loc[~is_ymd8] = pd.to_datetime(s[~is_ymd8], errors="coerce", dayfirst=False)

    return dt

def carregar_dados_detalhado(mes_atual, empresa_filter, num_meses=10):
    """Carrega dados detalhados (por item)"""
    query = QUERY_TEMPLATE.format(empresa_filter=empresa_filter , mes_atual=mes_atual)
    
    with get_conn() as conn:
        df = pd.read_sql(query, conn)

    df["Valor"] = pd.to_numeric(df["Valor"], errors="coerce").fillna(0.0)
    df["ContaContabil"] = df.get("ContaContabil", "").astype(str)

    dt = parse_periodo_to_datetime(df["Periodo"])
    df = df.assign(PeriodoDT=dt).dropna(subset=["PeriodoDT"]).copy()

    df["MesKey"] = df["PeriodoDT"].dt.to_period("M").astype(str)
    df["MesLabel"] = df["PeriodoDT"].dt.strftime("%m/%Y")

    mensal = (
        df.groupby(["ContaContabil", "Observação", "MesKey", "MesLabel"], as_index=False)["Valor"]
          .sum()
    )

    meses_key = sorted(mensal["MesKey"].unique())[-num_meses:]
    mensal_4 = mensal[mensal["MesKey"].isin(meses_key)].copy()

    key_to_label = (
        mensal_4.sort_values("MesKey")
                .drop_duplicates("MesKey")[["MesKey", "MesLabel"]]
    )
    meses_label_ordenados = key_to_label["MesLabel"].tolist()

    pivot = mensal_4.pivot_table(
        index=["ContaContabil", "Observação"],
        columns="MesLabel",
        values="Valor",
        aggfunc="sum",
        fill_value=0.0
    ).reset_index()

    for m in meses_label_ordenados:
        if m not in pivot.columns:
            pivot[m] = 0.0

    pivot["MediaMensal"] = pivot[meses_label_ordenados].mean(axis=1)
    pivot["ValorUltimoMes"] = pivot[meses_label_ordenados[-1]]

    LIMITE_MEDIA = 1  # ajuste conforme seu negócio
 
    pivot["DesvioPercentual"] = (
        (pivot["ValorUltimoMes"] - pivot["MediaMensal"])
        / pivot["MediaMensal"].where(pivot["MediaMensal"] >= LIMITE_MEDIA)
    ) * 100

    # Regra de negócio: média = 0 e houve consumo → 100%
    pivot.loc[
        (pivot["MediaMensal"] < LIMITE_MEDIA) & (pivot["ValorUltimoMes"] > 0),
        "DesvioPercentual"
    ] = 100.0

    # Segurança final
    pivot["DesvioPercentual"] = pivot["DesvioPercentual"].fillna(0.0)
#    pivot["DesvioPercentual"] = (
#        (pivot["ValorUltimoMes"] - pivot["MediaMensal"])
#        / pivot["MediaMensal"].where(pivot["MediaMensal"] != 0)
#    ) * 100
    pivot["DesvioPercentual"] = pivot["DesvioPercentual"].fillna(0.0)

    pivot["Risco"] = pivot.apply(
        lambda r: classify_risk(float(r["MediaMensal"]), float(r["ValorUltimoMes"])),
        axis=1
    )

    ordem = {"ALTO": 1, "MÉDIO": 2, "BAIXO": 3}
    pivot["Ordem"] = pivot["Risco"].map(ordem).fillna(9).astype(int)

    pivot = pivot.sort_values(
        ["Ordem", "DesvioPercentual", "ValorUltimoMes"],
        ascending=[True, False, False]
    ).reset_index(drop=True)

    for col in meses_label_ordenados + ["MediaMensal", "ValorUltimoMes", "DesvioPercentual"]:
        pivot[col] = pd.to_numeric(pivot[col], errors="coerce").fillna(0.0)

    return pivot, meses_label_ordenados

def carregar_dados_totalizado(mes_atual,empresa_filter, num_meses=10):
    """Carrega dados totalizados (por família)"""
    query = QUERY_TEMPLATE.format(empresa_filter=empresa_filter, mes_atual=mes_atual)
    
    with get_conn() as conn:
        df = pd.read_sql(query, conn)

    df["Valor"] = pd.to_numeric(df["Valor"], errors="coerce").fillna(0.0)
    df["ContaContabil"] = df.get("ContaContabil", "").astype(str)

    dt = parse_periodo_to_datetime(df["Periodo"])
    df = df.assign(PeriodoDT=dt).dropna(subset=["PeriodoDT"]).copy()

    df["MesKey"] = df["PeriodoDT"].dt.to_period("M").astype(str)
    df["MesLabel"] = df["PeriodoDT"].dt.strftime("%m/%Y")

    mensal = (
        df.groupby(["ContaContabil", "MesKey", "MesLabel"], as_index=False)["Valor"]
          .sum()
    )

    meses_key = sorted(mensal["MesKey"].unique())[-num_meses:]
    mensal_4 = mensal[mensal["MesKey"].isin(meses_key)].copy()

    key_to_label = (
        mensal_4.sort_values("MesKey")
                .drop_duplicates("MesKey")[["MesKey", "MesLabel"]]
    )
    meses_label_ordenados = key_to_label["MesLabel"].tolist()

    pivot = mensal_4.pivot_table(
        index=["ContaContabil"],
        columns="MesLabel",
        values="Valor",
        aggfunc="sum",
        fill_value=0.0
    ).reset_index()

    for m in meses_label_ordenados:
        if m not in pivot.columns:
            pivot[m] = 0.0

    pivot["MediaMensal"] = pivot[meses_label_ordenados].mean(axis=1)
    pivot["ValorUltimoMes"] = pivot[meses_label_ordenados[-1]]

    LIMITE_MEDIA = 1  # ajuste conforme seu negócio



    pivot["DesvioPercentual"] = (
        (pivot["ValorUltimoMes"] - pivot["MediaMensal"])
        / pivot["MediaMensal"].where(pivot["MediaMensal"] != 0)
    ) * 100
    pivot["DesvioPercentual"] = pivot["DesvioPercentual"].fillna(0.0)

    pivot["Risco"] = pivot.apply(
        lambda r: classify_risk(float(r["MediaMensal"]), float(r["ValorUltimoMes"])),
        axis=1
    )

    ordem = {"ALTO": 1, "MÉDIO": 2, "BAIXO": 3}
    pivot["Ordem"] = pivot["Risco"].map(ordem).fillna(9).astype(int)

    pivot = pivot.sort_values(
        ["Ordem", "DesvioPercentual", "ValorUltimoMes"],
        ascending=[True, False, False]
    ).reset_index(drop=True)

    for col in meses_label_ordenados + ["MediaMensal", "ValorUltimoMes", "DesvioPercentual"]:
        pivot[col] = pd.to_numeric(pivot[col], errors="coerce").fillna(0.0)

    return pivot, meses_label_ordenados