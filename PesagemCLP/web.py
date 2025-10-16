import logging
import sys
import threading
from pathlib import Path

from flask import Flask, jsonify, render_template
from waitress import serve

from state import ScaleState

log = logging.getLogger(__name__)


def _resource_path(rel_path: str) -> str:
    """Resolve a resource path that works in dev and inside a PyInstaller bundle."""
    base = getattr(sys, "_MEIPASS", None)
    if base:
        return str(Path(base) / rel_path)
    return str(Path(__file__).parent / rel_path)


def create_app(state: ScaleState) -> Flask:
    app = Flask(
        __name__,
        template_folder=_resource_path("templates"),
        static_folder=_resource_path("static"),
    )

    @app.route("/")
    def index():
        return render_template("index.html")

    @app.route("/api/scales")
    def api_scales():
        return jsonify(state.snapshot())

    @app.route("/healthz")
    def healthz():
        return {"status": "ok"}

    return app


class WebServerThread(threading.Thread):
    """Runs the Flask app on waitress in its own thread."""

    def __init__(self, state: ScaleState, host: str, port: int):
        super().__init__(name="web-server", daemon=True)
        self._app = create_app(state)
        self._host = host
        self._port = port

    def run(self) -> None:
        log.info("Web UI listening on http://%s:%d", self._host, self._port)
        serve(self._app, host=self._host, port=self._port, _quiet=True)
