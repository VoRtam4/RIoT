"""
@file rest_client.py
@brief REST klient pro volání exportů, historie a dalších REST operací Backend Core.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations

import time
from pathlib import Path
from urllib.parse import urlparse, urlunparse


class RESTClient:
    def __init__(self, http_client, rest_url: str):
        self.http = http_client
        self.rest_url = rest_url.rstrip("/")
        self.export_timeout_seconds = max(self.http.timeout_seconds * 20, 1800)
        self.export_poll_interval_seconds = 0.05

    def start_time_series_export(self, payload: dict) -> dict:
        response = self.http.post(self._to_path("/time-series/export"), json=payload)
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict) or "id" not in response.body or "status" not in response.body:
            raise RuntimeError("Unexpected export response shape")
        return response.body

    def start_aggregate_kpi_export(self, payload: dict) -> dict:
        response = self.http.post(self._to_path("/time-series/export/aggregate-kpi"), json=payload)
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict) or "id" not in response.body or "status" not in response.body:
            raise RuntimeError("Unexpected aggregate export response shape")
        return response.body

    def get_time_series_export(self, export_id: int) -> dict:
        response = self.http.get(self._to_path(f"/time-series/export/{export_id}/status"))
        if response.status_code >= 400:
            raise RuntimeError(response.body)
        if not isinstance(response.body, dict) or "status" not in response.body:
            raise RuntimeError("Unexpected export status response shape")
        return response.body

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
            export_job = self.start_aggregate_kpi_export(payload)
        else:
            export_job = self.start_time_series_export(payload)
        export_id = int(export_job["id"])
        download_url = self._wait_for_export_download_url(export_id)
        export_ready = time.perf_counter()
        self.download_file(download_url, target_path)
        finished = time.perf_counter()
        return {
            "download_url": download_url,
            "export_ready_time_s": export_ready - started,
            "download_time_s": finished - export_ready,
            "total_time_s": finished - started,
        }

    def _wait_for_export_download_url(self, export_id: int) -> str:
        deadline = time.perf_counter() + self.export_timeout_seconds
        last_job = None
        while True:
            export_job = self.get_time_series_export(export_id)
            last_job = export_job
            status = str(export_job.get("status", "")).lower()
            download_url = export_job.get("downloadUrl")
            if status == "done":
                if not download_url:
                    raise RuntimeError(f"Export {export_id} finished without downloadUrl")
                return str(download_url)
            if status in {"failed", "cancelled", "expired"}:
                error = export_job.get("error")
                raise RuntimeError(f"Export {export_id} ended with status={status}: {error}")
            if time.perf_counter() > deadline:
                raise TimeoutError(
                    f"Export {export_id} did not finish within {self.export_timeout_seconds} seconds. "
                    f"Last job={last_job!r}"
                )
            time.sleep(self.export_poll_interval_seconds)

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
