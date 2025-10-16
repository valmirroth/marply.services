import logging
import os
import signal
import sys
import threading
from collections import defaultdict

from config import load_database_config, load_dotenv, load_scales, stats_refresh_seconds
from database import WeightRepository
from plc_worker import PlcGroupWorker
from state import ScaleState
from stats_worker import StatsWorker
from web import WebServerThread


def _configure_logging() -> None:
    level = os.environ.get("LOG_LEVEL", "INFO").upper()
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s [%(threadName)s] %(name)s: %(message)s",
        stream=sys.stdout,
    )
    logging.getLogger("snap7.client").setLevel(logging.WARNING)
    logging.getLogger("snap7.common").setLevel(logging.CRITICAL)


def main() -> int:
    _configure_logging()
    log = logging.getLogger("main")

    loaded_env = load_dotenv()
    if loaded_env:
        log.info("Loaded environment from %s", loaded_env)
    else:
        log.info("No .env file found; using OS environment only")

    try:
        db_config = load_database_config()
        scales = load_scales()
    except Exception as exc:
        log.error("Failed to load configuration: %s", exc)
        return 2

    log.info("Loaded %d scales: %s", len(scales), [s.id_balanca for s in scales])

    state = ScaleState()
    for scale in scales:
        state.register(scale.id_balanca, scale.plc_ip)

    repository = WeightRepository(db_config)
    stop_event = threading.Event()

    def _shutdown(signum, _frame):
        log.info("Received signal %s, stopping workers", signum)
        stop_event.set()

    signal.signal(signal.SIGINT, _shutdown)
    signal.signal(signal.SIGTERM, _shutdown)

    web_host = os.environ.get("WEB_HOST", "0.0.0.0")
    web_port = int(os.environ.get("WEB_PORT", "8080"))
    web_thread = WebServerThread(state, web_host, web_port)
    web_thread.start()

    # Group scales by PLC endpoint so each CPU gets ONE shared Snap7 connection
    # instead of one-per-scale (Siemens PG/OP slots are limited).
    groups: dict[tuple, list] = defaultdict(list)
    for scale in scales:
        groups[(scale.plc_ip, scale.rack, scale.slot)].append(scale)

    log.info("Grouped into %d PLC connection(s):", len(groups))
    for (ip, rack, slot), group_scales in groups.items():
        log.info(
            "  PLC %s (rack=%d slot=%d) -> %s",
            ip, rack, slot, [s.id_balanca for s in group_scales],
        )

    workers = [
        PlcGroupWorker(ip, rack, slot, group_scales, repository, state, stop_event)
        for (ip, rack, slot), group_scales in groups.items()
    ]
    for worker in workers:
        worker.start()

    stats_worker = StatsWorker(repository, state, stop_event, stats_refresh_seconds())
    stats_worker.start()

    try:
        while not stop_event.is_set():
            stop_event.wait(timeout=1.0)
            for worker in workers:
                if not worker.is_alive():
                    log.error("Worker %s died unexpectedly; shutting down", worker.name)
                    stop_event.set()
                    break
    finally:
        for worker in workers:
            worker.join(timeout=10)
        stats_worker.join(timeout=5)
        repository.close()
        log.info("Shutdown complete")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
