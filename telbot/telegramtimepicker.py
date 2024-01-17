# A library that allows to create an inline timepicker keyboard.
# mitschix (23629789+mitschix@users.noreply.github.com)
# https://github.com/mitschix
#

import datetime
from typing import Tuple

from telegram import InlineKeyboardButton, InlineKeyboardMarkup


class TimePicker:
    def __init__(
        self,
        start_hour: int = 18,
        start_minute: int = 0,
        max_hour: int = 23,
        min_hour: int = 0,
        steps_hour: int = 1,
        steps_minute: int = 15,
    ):
        self.start_hour = start_hour
        self.start_minute = start_minute
        self.max_hour = max_hour
        self.min_hour = min_hour
        self.min_hour = min_hour
        self.steps_hour = steps_hour
        self.steps_minute = steps_minute

    def _create_callback_data(self, action, hour, minute):
        return ";".join([action, str(hour), str(minute)])

    def create_timepicker(self, hour: int = None, minute: int = None):
        keyboard_time = []

        hour = hour if hour or isinstance(hour, int) else self.start_hour
        minute = minute if minute else self.start_minute

        if hour + self.steps_hour <= self.max_hour:
            key_h_up = InlineKeyboardButton(
                "↑", callback_data=self._create_callback_data("UP-Hour", hour, minute)
            )
        else:
            key_h_up = InlineKeyboardButton(
                "-", callback_data=self._create_callback_data("IGNORE", hour, minute)
            )

        if hour - self.steps_hour >= self.min_hour:
            key_h_down = InlineKeyboardButton(
                "↓", callback_data=self._create_callback_data("DOWN-Hour", hour, minute)
            )
        else:
            key_h_down = InlineKeyboardButton(
                "-", callback_data=self._create_callback_data("IGNORE", hour, minute)
            )

        key_m_up = InlineKeyboardButton(
            "↑", callback_data=self._create_callback_data("UP-Min", hour, minute)
        )
        key_m_down = InlineKeyboardButton(
            "↓", callback_data=self._create_callback_data("DOWN-Min", hour, minute)
        )

        key_h = InlineKeyboardButton(
            f"{hour:02d}",
            callback_data=self._create_callback_data("IGNORE", hour, minute),
        )
        key_m = InlineKeyboardButton(
            f"{minute:02d}",
            callback_data=self._create_callback_data("IGNORE", hour, minute),
        )

        key_accept = InlineKeyboardButton(
            "OK", callback_data=self._create_callback_data("ACCEPT", hour, minute)
        )

        keyboard_time = [
            [key_h_up, key_m_up],
            [key_h, key_m],
            [key_h_down, key_m_down],
            [key_accept],
        ]

        return InlineKeyboardMarkup(keyboard_time)

    async def process_time_selection(
        self, update, context
    ) -> Tuple[bool, datetime.time]:
        query = update.callback_query
        action, hour, minute = query.data.split(";")
        print(query.data)

        # fix since +/- does not work with datetime.time
        now = datetime.datetime.now()
        curr = now.replace(hour=int(hour), minute=int(minute))
        match action:
            case "UP-Hour":
                new_time = curr + datetime.timedelta(hours=self.steps_hour)
            case "UP-Min":
                new_time = curr + datetime.timedelta(minutes=self.steps_minute)
            case "DOWN-Hour":
                new_time = curr - datetime.timedelta(hours=self.steps_hour)
            case "DOWN-Min":
                new_time = curr - datetime.timedelta(minutes=self.steps_minute)
            case "ACCEPT":
                return True, curr.time()
            case _:
                return False, curr.time()
        await context.bot.edit_message_text(
            text=query.message.text,
            chat_id=query.message.chat_id,
            message_id=query.message.message_id,
            reply_markup=self.create_timepicker(
                int(new_time.hour), int(new_time.minute)
            ),
        )
        return False, curr.time()
