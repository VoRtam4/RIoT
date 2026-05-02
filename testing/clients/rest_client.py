from __future__ import annotations

import time
from pathlib import Path
from urllib.parse import urlparse, urlunparse


class RESTClient:
    def __init__(self, http_client, rest_url: str):
        self.http = http_client
        self.rest_url = rest_url.rstrip("/")

    def start_time_series_export(self, payload: dict) -> str:
        response = self.http.post(self._to_path("/time-series/export"), json=payload)
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict) or "url" not in response.body:
            raise RuntimeError("Unexpected export response shape")
        return response.body["url"]

    def start_aggregate_kpi_export(self, payload: dict) -> str:
        response = self.http.post(self._to_path("/time-series/export/aggregate-kpi"), json=payload)
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict) or "url" not in response.body:
            raise RuntimeError("Unexpected aggregate export response shape")
        return response.body["url"]

    def distinct_tag_values(self, payload: dict) -> list[str]:
        response = self.http.post(self._to_path("/time-series/distinct-tag-values"), json=payload)
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict):
            raise RuntimeError("Unexpected distinct tag values response shape")
        values = response.body.get("values", [])
        if isinstance(values, dict):
            values = list(values.keys())
        if not isinstance(values, list):
            raise RuntimeError(f"Unexpected distinct tag values payload: {values!r}")
        return [str(value) for value in values if value not in (None, "")]

    def download_file(self, absolute_url: str, target_path: Path) -> Path:
        absolute_url = self._normalize_download_url(absolute_url)
        response = self.http.session.get(
            absolute_url,
            timeout=(self.http.timeout_seconds, None),
            stream=True,
        )
        if response.status_code >= 400:
            body = response.text
            raise RuntimeError(f"Export download failed: {response.status_code} {body}")
        target_path.parent.mkdir(parents=True, exist_ok=True)
        with target_path.open("wb") as handle:
            for chunk in response.iter_content(chunk_size=65536):
                if chunk:
                    handle.write(chunk)
        return target_path

    def timed_export_download(self, payload: dict, target_path: Path, aggregate: bool = False) -> dict:
        started = time.perf_counter()
        if aggregate:
            download_url = self.start_aggregate_kpi_export(payload)
        else:
            download_url = self.start_time_series_export(payload)
        export_ready = time.perf_counter()
        self.download_file(download_url, target_path)
        finished = time.perf_counter()
        return {
            "download_url": download_url,
            "export_ready_time_s": export_ready - started,
            "download_time_s": finished - export_ready,
            "total_time_s": finished - started,
        }

    def _to_path(self, suffix: str) -> str:
        if self.rest_url.startswith(self.http.base_url):
            base_path = self.rest_url[len(self.http.base_url) :]
        else:
            base_path = self.rest_url
        return f"{base_path}{suffix}"

    def _normalize_download_url(self, absolute_url: str) -> str:
        parsed = urlparse(absolute_url)
        if not parsed.scheme or not parsed.netloc:
            if absolute_url.startswith("/"):
                return f"{self.http.base_url}{absolute_url}"
            return f"{self.rest_url}/{absolute_url.lstrip('/')}"

        rest_parsed = urlparse(self.rest_url)
        base_http = urlparse(self.http.base_url)

        if parsed.hostname in {"riot-backend-core", "backend-core"}:
            replacement = rest_parsed if rest_parsed.scheme and rest_parsed.netloc else base_http
            return urlunparse(
                (
                    replacement.scheme,
                    replacement.netloc,
                    parsed.path,
                    parsed.params,
                    parsed.query,
                    parsed.fragment,
                )
            )
        return absolute_url
