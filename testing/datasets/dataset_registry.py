"""
@file dataset_registry.py
@brief Registr dostupných datových sad a jejich parametrů pro experimenty.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

from __future__ import annotations


class DatasetRegistry:
    def __init__(self, state_manager):
        self.state_manager = state_manager

    def exists(self, dataset_id: str) -> bool:
        return dataset_id in self.state_manager.load_dataset_registry()
