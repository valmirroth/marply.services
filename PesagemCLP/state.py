import threading
from dataclasses import dataclass, field
from datetime import datetime


@dataclass
class ScaleReading:
    id_balanca: str
    plc_ip: str
    connected: bool = False
    last_weight: float | None = None
    last_timestamp: datetime | None = None
    last_error: str | None = None
    last_error_at: datetime | None = None
    stats: dict[str, int] = field(default_factory=dict)


@dataclass
class StatsMeta:
    """Tracks whether the latest stats query succeeded."""
    updated_at: datetime | None = None
    error: str | None = None


class ScaleState:
    """Thread-safe in-memory state shared between PLC workers and the web UI."""

    def __init__(self):
        self._lock = threading.Lock()
        self._scales: dict[str, ScaleReading] = {}
        self._stats_meta = StatsMeta()

    def register(self, id_balanca: str, plc_ip: str) -> None:
        with self._lock:
            self._scales[id_balanca] = ScaleReading(id_balanca=id_balanca, plc_ip=plc_ip)

    def update_weight(self, id_balanca: str, weight: float, moment: datetime) -> None:
        with self._lock:
            reading = self._scales.get(id_balanca)
            if reading is None:
                return
            reading.last_weight = weight
            reading.last_timestamp = moment

    def set_connected(self, id_balanca: str, connected: bool) -> None:
        with self._lock:
            reading = self._scales.get(id_balanca)
            if reading is None:
                return
            reading.connected = connected

    def set_error(self, id_balanca: str, error: str) -> None:
        with self._lock:
            reading = self._scales.get(id_balanca)
            if reading is None:
                return
            reading.last_error = error
            reading.last_error_at = datetime.now()

    def clear_error(self, id_balanca: str) -> None:
        with self._lock:
            reading = self._scales.get(id_balanca)
            if reading is None:
                return
            reading.last_error = None
            reading.last_error_at = None

    def update_stats_per_scale(self, per_scale: dict[str, dict[str, int]]) -> None:
        """Update each known scale's stats from a {scale: {status: count}} mapping.

        Scales present in the state but not in `per_scale` are reset to zero
        (the DB has no rows for them today). Scales in `per_scale` that aren't
        registered in the state are ignored (they may belong to other systems).
        """
        with self._lock:
            for reading in self._scales.values():
                reading.stats = dict(per_scale.get(reading.id_balanca, {}))
            self._stats_meta.updated_at = datetime.now()
            self._stats_meta.error = None

    def set_stats_error(self, error: str) -> None:
        with self._lock:
            self._stats_meta.error = error
            self._stats_meta.updated_at = datetime.now()

    def snapshot(self) -> dict:
        with self._lock:
            scales = [
                {
                    "id_balanca": r.id_balanca,
                    "plc_ip": r.plc_ip,
                    "connected": r.connected,
                    "last_weight": r.last_weight,
                    "last_timestamp": r.last_timestamp.isoformat(timespec="seconds")
                    if r.last_timestamp
                    else None,
                    "last_error": r.last_error,
                    "last_error_at": r.last_error_at.isoformat(timespec="seconds")
                    if r.last_error_at
                    else None,
                    "has_error": r.last_error is not None or not r.connected,
                    "stats": dict(r.stats),
                }
                for r in self._scales.values()
            ]
            stats_meta = {
                "updated_at": self._stats_meta.updated_at.isoformat(timespec="seconds")
                if self._stats_meta.updated_at
                else None,
                "error": self._stats_meta.error,
            }
            return {"scales": scales, "stats_meta": stats_meta}
