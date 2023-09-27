"""Utility class to handle host events"""

import sys
import json
import os
from pathlib import Path

from typing import List, Dict

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

        with self.json_path.open(encoding='utf8') as j_f:
            content = json.load(j_f)
        self.events = content.get('host_events')

    def get_events_per_day(self, day: str) -> T_EVENTS:
        events = []
        for eve in self.events:
            if eve.get('day') == day:
                events.append(eve)
        return events

    def update(self) -> None:
        """scrapes all websites and rereads json"""
        os.system("./evenpargo")
        self._read_events()


def format_events(events: T_EVENTS) -> str:
    out = ""
    for eve in events:
        event_infos = '\n- '.join(eve.get('events'))
        out += f"{eve.get('host')}\n- {event_infos}"
        out += f"\n{10*'-----'}\n"
    return out
