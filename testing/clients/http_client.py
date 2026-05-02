from __future__ import annotations

from dataclasses import dataclass

import requests


@dataclass
class HTTPResponse:
    status_code: int
    body: object
    headers: dict


class HTTPClient:
    def __init__(self, base_url: str, api_key: str, timeout_seconds: int = 30):
        self.base_url = base_url.rstrip("/")
        self.timeout_seconds = timeout_seconds
        self.session = requests.Session()
        self.session.headers.update(
            {
                "X-API-Key": api_key,
                "Content-Type": "application/json",
            }
        )

    def get(self, path: str, **kwargs) -> HTTPResponse:
        response = self.session.get(self.base_url + path, timeout=self.timeout_seconds, **kwargs)
        return self._wrap(response)

    def post(self, path: str, json: dict | None = None, **kwargs) -> HTTPResponse:
        response = self.session.post(self.base_url + path, json=json, timeout=self.timeout_seconds, **kwargs)
        return self._wrap(response)

    def _wrap(self, response: requests.Response) -> HTTPResponse:
        try:
            body = response.json()
        except Exception:
            body = response.text
        return HTTPResponse(
            status_code=response.status_code,
            body=body,
            headers=dict(response.headers),
        )
