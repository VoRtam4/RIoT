from __future__ import annotations


class DatasetRegistry:
    def __init__(self, state_manager):
        self.state_manager = state_manager

    def exists(self, dataset_id: str) -> bool:
        return dataset_id in self.state_manager.load_dataset_registry()
