from __future__ import annotations


class DailyProfile:
    def multiplier_for_minute_of_day(self, minute_of_day: int) -> float:
        raise NotImplementedError


class WazeDailyProfile(DailyProfile):
    def multiplier_for_minute_of_day(self, minute_of_day: int) -> float:
        hour = minute_of_day // 60
        if 0 <= hour < 5:
            return 0.35
        if 5 <= hour < 7:
            return 0.55
        if 7 <= hour < 10:
            return 1.25
        if 10 <= hour < 14:
            return 0.90
        if 14 <= hour < 18:
            return 1.35
        if 18 <= hour < 22:
            return 0.85
        return 0.50


class MHDDailyProfile(DailyProfile):
    def multiplier_for_minute_of_day(self, minute_of_day: int) -> float:
        hour = minute_of_day // 60
        minute = minute_of_day % 60
        if hour == 0 and minute >= 30:
            return 0.04
        if 1 <= hour < 4:
            return 0.02
        if hour == 4 and minute < 30:
            return 0.06
        if 5 <= hour < 7:
            return 0.65
        if 7 <= hour < 10:
            return 1.20
        if 10 <= hour < 14:
            return 0.95
        if 14 <= hour < 18:
            return 1.15
        if 18 <= hour < 22:
            return 0.90
        return 0.45


class FlatProfile(DailyProfile):
    def multiplier_for_minute_of_day(self, minute_of_day: int) -> float:
        _ = minute_of_day
        return 1.0
