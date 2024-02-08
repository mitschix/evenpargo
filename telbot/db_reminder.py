import datetime
import sqlite3
from typing import Dict, List, Optional, Union

from config import DB_NAME


class ReminderDB:
    def __init__(self, db_name):
        self._db_name = db_name
        self._create_tables()

    def _create_tables(self) -> None:
        with self._get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute(
                """CREATE TABLE IF NOT EXISTS reminders \
                (userid INTEGER PRIMARY KEY, day INTEGER, time TEXT, \
                state INTEGER)"""
            )

    def _get_connection(self) -> sqlite3.Connection:
        conn = sqlite3.connect(self._db_name)
        return conn

    def _execute_query(self, query, params=None) -> None:
        with self._get_connection() as conn:
            cursor = conn.cursor()
            _ = cursor.execute(query, params) if params else cursor.execute(query)
            conn.commit()

    def _read_data(self, query, params=None) -> List:
        with self._get_connection() as conn:
            cursor = conn.cursor()
            result = (
                cursor.execute(query, params).fetchall()
                if params
                else cursor.execute(query).fetchall()
            )
            return result

    def get_reminder(self, userid) -> Optional[Dict[str, Union[int, datetime.time]]]:
        result = self._read_data("SELECT * FROM reminders WHERE userid=?", (userid,))
        if result:
            result = result[0]
            reminder = {
                "userid": result[0],
                "day": result[1],
                "time": datetime.datetime.strptime(result[2], "%H:%M:%S").time(),
                "state": result[3],
            }
            return reminder
        return None

    def set_reminder(self, userid: int, day: int, time: str) -> None:
        self._execute_query(
            "REPLACE INTO reminders (userid, day, time, state) VALUES (?, ?, ?, ?)",
            (userid, day, time, 1),
        )

    def toggle_state(self, userid: int) -> None:
        state = self._read_data(
            "SELECT state FROM reminders WHERE userid=?", (userid,)
        )[0][0]
        new_state = 1 if state == 0 else 0
        self._execute_query(
            "UPDATE reminders SET state=? WHERE userid=?",
            (
                new_state,
                userid,
            ),
        )

    def delete_reminder_by_user_id(self, userid: int) -> None:
        """
        Deletes an entry from the database table given its userid.
        """
        self._execute_query("DELETE FROM reminders WHERE userid = ?", (userid,))

    def list_reminders(self) -> None:
        print(self._read_data("SELECT * FROM reminders"))


if __name__ == "__main__":
    db = ReminderDB(DB_NAME)
    db.list_reminders()
