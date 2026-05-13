"""
@file __init__.py
@brief Balíček experimentů ověřujících ingest, historii, reprocessing a ukládání dat.

@author Vojtěch Hubáček

@par Autorský podíl
- Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.

@ingroup riot_testing
"""

"""Experiment implementations."""

from testing.experiments.history_read_experiment import HistoryReadExperiment
from testing.experiments.ingest_experiment import IngestExperiment
from testing.experiments.reprocess_experiment import ReprocessExperiment
from testing.experiments.storage_experiment import StorageExperiment


def build_experiment(ctx, experiment_type: str):
    if experiment_type == "ingest":
        return IngestExperiment()
    if experiment_type == "reprocess":
        return ReprocessExperiment()
    if experiment_type == "history_read":
        return HistoryReadExperiment()
    if experiment_type == "storage":
        return StorageExperiment()
    raise ValueError(f"Unsupported experiment type: {experiment_type}")
