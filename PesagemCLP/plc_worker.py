import logging
import threading
from datetime import datetime

import snap7
from snap7.util import get_bool, get_real

from config import ScaleConfig
from database import WeightRepository
from state import ScaleState

log = logging.getLogger(__name__)

# python-snap7's check_error raises RuntimeError for any PLC/transport failure.
PlcError = RuntimeError

_RECONNECT_BACKOFF_MIN_S = 2.0
_RECONNECT_BACKOFF_MAX_S = 30.0


def _format_error(exc: BaseException) -> str:
    """Render an exception as a clean one-line message for the UI."""
    message = str(exc).strip()
    if not message:
        message = exc.__class__.__name__
    if message.startswith("b'") and message.endswith("'"):
        message = message[2:-1]
    if message.startswith("b\"") and message.endswith("\""):
        message = message[2:-1]
    return message.strip()


class PlcGroupWorker(threading.Thread):
    """Polls one PLC for all scales attached to it.

    Several scales (each with their own pulse-bit and weight DB offset) share
    a single Snap7 connection. This avoids opening N connections to the same
    Siemens CPU — its PG/OP connection slots are limited (typically 3).
    """

    def __init__(
        self,
        plc_ip: str,
        rack: int,
        slot: int,
        scales: list[ScaleConfig],
        repository: WeightRepository,
        state: ScaleState,
        stop_event: threading.Event,
    ):
        if not scales:
            raise ValueError("PlcGroupWorker requires at least one scale")
        super().__init__(name=f"plc-{plc_ip}", daemon=True)
        self._plc_ip = plc_ip
        self._rack = rack
        self._slot = slot
        self._scales = list(scales)
        self._repository = repository
        self._state = state
        self._stop_event = stop_event
        self._client: snap7.client.Client | None = None
        self._last_pulse: dict[str, bool] = {s.id_balanca: False for s in scales}
        self._backoff_s = _RECONNECT_BACKOFF_MIN_S
        # Use the smallest poll interval among scales in the group.
        self._poll_interval_s = min(s.poll_interval_s for s in scales)

    # ------------------------------------------------------------------ loop
    def run(self) -> None:
        names = [s.id_balanca for s in self._scales]
        log.info("[%s] worker starting with %d scale(s): %s", self._plc_ip, len(names), names)
        try:
            while not self._stop_event.is_set():
                if not self._ensure_connected():
                    continue
                try:
                    self._poll_once()
                except PlcError as exc:
                    msg = _format_error(exc)
                    log.warning("[%s] PLC error during poll, will reconnect: %s", self._plc_ip, msg)
                    self._mark_group_offline(f"PLC: {msg}")
                    self._disconnect()
                    self._sleep(_RECONNECT_BACKOFF_MIN_S)
                    continue
                except Exception as exc:
                    msg = _format_error(exc)
                    log.exception("[%s] unexpected error in poll loop", self._plc_ip)
                    self._mark_group_offline(msg)
                    self._sleep(_RECONNECT_BACKOFF_MIN_S)
                    continue

                self._sleep(self._poll_interval_s)
        finally:
            self._disconnect()
            self._mark_group_offline("worker parado")
            log.info("[%s] worker stopped", self._plc_ip)

    def _sleep(self, seconds: float) -> None:
        self._stop_event.wait(timeout=seconds)

    # ------------------------------------------------------------ connection
    def _ensure_connected(self) -> bool:
        if self._client is not None and self._client.get_connected():
            return True
        try:
            client = snap7.client.Client()
            client.connect(self._plc_ip, self._rack, self._slot)
        except PlcError as exc:
            msg = _format_error(exc)
            log.warning(
                "[%s] connect failed (%s); retrying in %.1fs",
                self._plc_ip, msg, self._backoff_s,
            )
            self._mark_group_offline(f"CLP inacessível: {msg}")
            self._sleep(self._backoff_s)
            self._backoff_s = min(self._backoff_s * 2, _RECONNECT_BACKOFF_MAX_S)
            return False

        log.info("[%s] connected to PLC (%d scales)", self._plc_ip, len(self._scales))
        self._client = client
        for sid in self._last_pulse:
            self._last_pulse[sid] = False
        self._backoff_s = _RECONNECT_BACKOFF_MIN_S
        for scale in self._scales:
            self._state.set_connected(scale.id_balanca, True)
            self._state.clear_error(scale.id_balanca)
        return True

    def _disconnect(self) -> None:
        if self._client is None:
            return
        try:
            self._client.disconnect()
        except PlcError:
            pass
        self._client = None

    def _mark_group_offline(self, error_msg: str) -> None:
        for scale in self._scales:
            self._state.set_connected(scale.id_balanca, False)
            self._state.set_error(scale.id_balanca, error_msg)

    # ------------------------------------------------------------ polling
    def _poll_once(self) -> None:
        """Read each scale's pulse bit + weight in turn. One connection, many scales."""
        assert self._client is not None
        for scale in self._scales:
            condition_bytes = self._client.mb_read(scale.pulse_byte, 1)
            current_pulse = get_bool(condition_bytes, 0, scale.pulse_bit)
            last = self._last_pulse[scale.id_balanca]

            if current_pulse and not last:
                weight_bytes = self._client.db_read(
                    scale.weight_db, scale.weight_offset, 4
                )
                weight = get_real(weight_bytes, 0)
                moment = datetime.now()
                log.info("[%s] pulse detected, weight=%.3f", scale.id_balanca, weight)
                self._state.update_weight(scale.id_balanca, weight, moment)
                try:
                    self._repository.insert_weight(scale.id_balanca, weight)
                    self._state.clear_error(scale.id_balanca)
                except Exception as exc:
                    msg = _format_error(exc)
                    log.exception("[%s] failed to persist weight=%.3f", scale.id_balanca, weight)
                    self._state.set_error(scale.id_balanca, f"Banco: {msg}")

            self._last_pulse[scale.id_balanca] = current_pulse


# Backwards-compatible alias — older code/import paths can still reference PlcWorker.
PlcWorker = PlcGroupWorker
