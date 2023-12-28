"""Utility class to handle host events"""

import json
import os
import sys
from pathlib import Path
from typing import Dict, List

T_EVENTS = List[Dict[str, str]]


class HostEventHandler(object):
    """docstring for ClassName."""

    def __init__(self, json_path: str):
        self.json_path = Path(json_path)
        self.events = []
        self._read_events()

    def _read_events(self) -> None:
        if not self.json_path.exists():
            print("file not found")
            sys.exit()

        with self.json_path.open(encoding="utf8") as j_f:
            content = json.load(j_f)
        self.events = content.get("host_events")

    def get_events_per_day(self, day: str) -> T_EVENTS:
        events = []
        for eve in self.events:
            if eve.get("day") == day:
                events.append(eve)
        return events

    def update(self) -> None:
        """scrapes all websites and rereads json"""
        os.system("./evenpargo")
        self._read_events()


def format_events(events: T_EVENTS) -> str:
    out = ""
    for club in events:
        out += f"{club.get('host')}"
        event_infos = club.get("events", [])
        for info in event_infos:
            val = list(info.values())
            print(val)
            out += f"\n- {val[1]}: {val[0]}\n({val[2]})\n"
        out += f"\n{10*'-----'}\n"
    return out
