import logging
import threading

import pyodbc
from tenacity import (
    retry,
    retry_if_exception_type,
    stop_after_attempt,
    wait_exponential,
)

from config import DatabaseConfig

log = logging.getLogger(__name__)


class WeightRepository:
    """Thread-safe writer for weight readings. Reconnects on transient failure."""

    def __init__(self, config: DatabaseConfig):
        self._config = config
        self._lock = threading.Lock()
        self._conn: pyodbc.Connection | None = None

    def _ensure_connection(self) -> pyodbc.Connection:
        if self._conn is None:
            log.info("Opening SQL connection to %s/%s", self._config.server, self._config.database)
            self._conn = pyodbc.connect(self._config.connection_string, autocommit=False)
        return self._conn

    def _reset_connection(self) -> None:
        if self._conn is not None:
            try:
                self._conn.close()
            except pyodbc.Error:
                pass
            self._conn = None

    @retry(
        retry=retry_if_exception_type(pyodbc.Error),
        stop=stop_after_attempt(5),
        wait=wait_exponential(multiplier=1, min=1, max=30),
        reraise=True,
    )
    def insert_weight(self, id_balanca: str, weight: float) -> None:
        """Cancel any PENDING weighings for this scale and insert the new reading.

        Runs both statements in a single transaction so we never end up with
        the cancel applied but the insert lost (or vice versa).
        """
        cancel_sql = (
            f"UPDATE {self._config.table} "
            "SET STATUS_PESAGEM = 'CANCELADO', DTH_CANCELA_PESO = GETDATE() "
            "WHERE STATUS_PESAGEM = 'PENDENTE' AND ID_BALANCA = ?"
        )
        insert_sql = f"INSERT INTO {self._config.table} (ID_BALANCA, PESO) VALUES (?, ?)"
        with self._lock:
            try:
                conn = self._ensure_connection()
                cursor = conn.cursor()
                cursor.execute(cancel_sql, (id_balanca,))
                cancelled = cursor.rowcount
                cursor.execute(insert_sql, (id_balanca, weight))
                conn.commit()
                cursor.close()
                if cancelled > 0:
                    log.info(
                        "[%s] cancelled %d pending weighing(s) before insert",
                        id_balanca, cancelled,
                    )
            except pyodbc.Error:
                log.exception("SQL write failed; will reset connection and retry")
                self._reset_connection()
                raise

    @retry(
        retry=retry_if_exception_type(pyodbc.Error),
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=1, max=10),
        reraise=True,
    )
    def count_by_status_today_per_scale(self) -> dict[str, dict[str, int]]:
        """Count today's rows grouped by ID_BALANCA and STATUS_PESAGEM.

        Returns: { "Balanca_1": {"COLETADO": 5, "PENDENTE": 1}, ... }
        Status keys are upper-cased and trimmed so the UI can match reliably
        regardless of how rows were originally written.
        """
        sql = (
            "SELECT ID_BALANCA, UPPER(LTRIM(RTRIM(STATUS_PESAGEM))) AS STATUS, COUNT(*) AS QTD "
            f"FROM {self._config.table} "
            f"WHERE CAST({self._config.date_column} AS DATE) = CAST(GETDATE() AS DATE) "
            "AND STATUS_PESAGEM IS NOT NULL "
            "AND ID_BALANCA IS NOT NULL "
            "GROUP BY ID_BALANCA, UPPER(LTRIM(RTRIM(STATUS_PESAGEM)))"
        )
        with self._lock:
            try:
                conn = self._ensure_connection()
                cursor = conn.cursor()
                cursor.execute(sql)
                result: dict[str, dict[str, int]] = {}
                for row in cursor.fetchall():
                    id_balanca = str(row[0]).strip() if row[0] is not None else ""
                    status = row[1]
                    qtd = int(row[2])
                    if not id_balanca:
                        continue
                    result.setdefault(id_balanca, {})[status] = qtd
                cursor.close()
                return result
            except pyodbc.Error:
                log.exception("SQL count_by_status_today_per_scale failed; resetting connection")
                self._reset_connection()
                raise

    def close(self) -> None:
        with self._lock:
            self._reset_connection()
