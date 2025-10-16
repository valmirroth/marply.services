import json
import os
import sys
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class DatabaseConfig:
    server: str
    database: str
    username: str
    password: str
    driver: str
    table: str
    date_column: str

    @property
    def connection_string(self) -> str:
        return (
            f"DRIVER={{{self.driver}}};"
            f"SERVER={self.server};"
            f"DATABASE={self.database};"
            f"UID={self.username};"
            f"PWD={self.password};"
            f"TrustServerCertificate=yes;"
        )


@dataclass(frozen=True)
class ScaleConfig:
    id_balanca: str
    plc_ip: str
    rack: int
    slot: int
    pulse_byte: int
    pulse_bit: int
    weight_db: int
    weight_offset: int
    poll_interval_s: float


def _dotenv_candidates() -> list[Path]:
    """Possible locations for a .env file, in priority order."""
    paths: list[Path] = [Path.cwd() / ".env"]
    if getattr(sys, "frozen", False):
        # Running from a PyInstaller bundle — also look next to the exe.
        paths.append(Path(sys.executable).resolve().parent / ".env")
    else:
        # Running from source — also look next to this module.
        paths.append(Path(__file__).resolve().parent / ".env")
    # Deduplicate preserving order.
    seen: set[Path] = set()
    unique: list[Path] = []
    for p in paths:
        if p not in seen:
            seen.add(p)
            unique.append(p)
    return unique


def load_dotenv() -> Path | None:
    """Load the first .env file found into os.environ.

    Existing environment variables are NOT overwritten — values already in
    os.environ win, which lets NSSM/systemd/docker-compose override the file.
    Returns the path that was loaded, or None if no .env was found.
    """
    for path in _dotenv_candidates():
        if not path.is_file():
            continue
        for raw_line in path.read_text(encoding="utf-8").splitlines():
            line = raw_line.strip()
            if not line or line.startswith("#") or "=" not in line:
                continue
            key, _, value = line.partition("=")
            key = key.strip()
            value = value.strip()
            # Strip surrounding quotes if present.
            if len(value) >= 2 and value[0] == value[-1] and value[0] in ('"', "'"):
                value = value[1:-1]
            if key and key not in os.environ:
                os.environ[key] = value
        return path
    return None


def _require_env(name: str) -> str:
    value = os.environ.get(name)
    if not value:
        raise RuntimeError(f"Required environment variable '{name}' is not set")
    return value


def load_database_config() -> DatabaseConfig:
    return DatabaseConfig(
        server=_require_env("SQL_SERVER"),
        database=_require_env("SQL_DATABASE"),
        username=_require_env("SQL_USERNAME"),
        password=_require_env("SQL_PASSWORD"),
        driver=os.environ.get("SQL_DRIVER", "ODBC Driver 18 for SQL Server"),
        table=os.environ.get("SQL_TABLE", "PESAGEM_ACABAMENTO"),
        date_column=os.environ.get("SQL_DATE_COLUMN", "DATAHORA"),
    )


def stats_refresh_seconds() -> float:
    raw = os.environ.get("STATS_REFRESH_S", "10")
    try:
        return max(1.0, float(raw))
    except ValueError:
        return 10.0


def _scales_path_candidates() -> list[Path]:
    """Possible locations for scales.json, in priority order."""
    env_path = os.environ.get("SCALES_CONFIG_PATH")
    paths: list[Path] = []
    if env_path:
        paths.append(Path(env_path))
    paths.append(Path.cwd() / "scales.json")
    if getattr(sys, "frozen", False):
        paths.append(Path(sys.executable).resolve().parent / "scales.json")
    else:
        paths.append(Path(__file__).resolve().parent / "scales.json")
    seen: set[Path] = set()
    unique: list[Path] = []
    for p in paths:
        if p not in seen:
            seen.add(p)
            unique.append(p)
    return unique


def load_scales(path: str | Path | None = None) -> list[ScaleConfig]:
    if path is not None:
        candidates = [Path(path)]
    else:
        candidates = _scales_path_candidates()

    resolved: Path | None = None
    for candidate in candidates:
        if candidate.is_file():
            resolved = candidate
            break
    if resolved is None:
        tried = ", ".join(str(p) for p in candidates)
        raise FileNotFoundError(f"Scales config file not found. Tried: {tried}")

    raw = json.loads(resolved.read_text(encoding="utf-8"))
    if not isinstance(raw, list) or not raw:
        raise ValueError(f"{resolved} must contain a non-empty JSON array of scales")

    scales: list[ScaleConfig] = []
    for i, entry in enumerate(raw):
        try:
            scales.append(
                ScaleConfig(
                    id_balanca=entry["id_balanca"],
                    plc_ip=entry["plc_ip"],
                    rack=int(entry.get("rack", 0)),
                    slot=int(entry.get("slot", 1)),
                    pulse_byte=int(entry["pulse_byte"]),
                    pulse_bit=int(entry["pulse_bit"]),
                    weight_db=int(entry["weight_db"]),
                    weight_offset=int(entry.get("weight_offset", 0)),
                    poll_interval_s=float(entry.get("poll_interval_s", 0.5)),
                )
            )
        except (KeyError, ValueError, TypeError) as exc:
            raise ValueError(f"Invalid scale entry at index {i}: {exc}") from exc

    ids = [s.id_balanca for s in scales]
    if len(set(ids)) != len(ids):
        raise ValueError(f"Duplicate id_balanca values in {resolved}: {ids}")

    return scales
