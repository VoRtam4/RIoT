from __future__ import annotations

from abc import ABC, abstractmethod


class BaseExperiment(ABC):
    @abstractmethod
    def prepare(self, ctx, scenario) -> None:
        raise NotImplementedError

    @abstractmethod
    def run_once(self, ctx, scenario, repetition):
        raise NotImplementedError

    def run(self, ctx, scenario):
        self.prepare(ctx, scenario)
        results = []
        for repetition in range(1, scenario.repetitions + 1):
            results.append(self.run_once(ctx, scenario, repetition))
        return results
