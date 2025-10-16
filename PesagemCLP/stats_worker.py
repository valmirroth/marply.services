import logging
import threading

from database import WeightRepository
from state import ScaleState

log = logging.getLogger(__name__)


def _format_error(exc: BaseException) -> str:
    message = str(exc).strip()
    if not message:
        message = exc.__class__.__name__
    return message


class StatsWorker(threading.Thread):
    """Periodically queries the DB for today's per-scale status counts and updates state."""

    def __init__(
        self,
        repository: WeightRepository,
        state: ScaleState,
        stop_event: threading.Event,
        interval_s: float,
    ):
        super().__init__(name="stats-worker", daemon=True)
        self._repo = repository
        self._state = state
        self._stop_event = stop_event
        self._interval_s = max(1.0, float(interval_s))

    def run(self) -> None:
        log.info("Stats worker starting (interval=%.1fs)", self._interval_s)
        self._fetch_once()
        while not self._stop_event.wait(self._interval_s):
            self._fetch_once()
        log.info("Stats worker stopped")

    def _fetch_once(self) -> None:
        try:
            per_scale = self._repo.count_by_status_today_per_scale()
            self._state.update_stats_per_scale(per_scale)
            log.debug("Stats updated: %s", per_scale)
        except Exception as exc:
            msg = _format_error(exc)
            log.warning("Stats query failed: %s", msg)
            self._state.set_stats_error(msg)
