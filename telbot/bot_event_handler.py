import datetime

import pytz
from telegram import InlineKeyboardMarkup, Update
from telegram.ext import (
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    JobQueue,
    filters,
)

from bot_keyboards import keyboard_days
from config import SUPPORT_ID
from parse_eve import HostEventHandler, format_events

EVENTS = HostEventHandler("./events.json")


async def update_events_job(context: ContextTypes.DEFAULT_TYPE):
    EVENTS.update()
    await context.bot.send_message(
        chat_id=SUPPORT_ID, text="(Job) Events should be up-to-date. (: "
    )


async def update_events(update: Update, context: ContextTypes.DEFAULT_TYPE):
    EVENTS.update()
    await context.bot.send_message(
        chat_id=update.effective_chat.id, text="Events should be up-to-date. (: "
    )


async def get_events(update: Update, context: ContextTypes.DEFAULT_TYPE):
    reply_markup = InlineKeyboardMarkup(keyboard_days)
    await context.bot.send_message(
        chat_id=update.message.chat_id,
        text="Which day you want to choose?",
        reply_markup=reply_markup,
    )


async def handle_events(update: Update, context: ContextTypes.DEFAULT_TYPE):
    q_data = update.callback_query.data
    q_day = q_data.split("_")[1].title()
    event_msg = format_events(EVENTS.get_events_per_day(q_day))
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=event_msg,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


def set_update_jobs(jobq: JobQueue) -> None:
    daily_upadte_time = datetime.time(hour=22, tzinfo=pytz.timezone("Europe/Berlin"))
    _ = jobq.run_daily(update_events_job, daily_upadte_time)
    daily_upadte_time = datetime.time(hour=10, tzinfo=pytz.timezone("Europe/Berlin"))
    _ = jobq.run_daily(update_events_job, daily_upadte_time)


event_update_h = CommandHandler("update", update_events, filters.User(SUPPORT_ID))
event_get_h = CommandHandler("events", get_events)
event_show_h = CallbackQueryHandler(handle_events)
