#!/usr/bin/env python
"""Telegram bot to read events.json and returns events"""
import datetime
import logging

from telegram import InlineKeyboardMarkup, Update
from telegram.ext import (
    ApplicationBuilder,
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    MessageHandler,
    filters,
)

from bot_keyboards import keyboard_days
from config import MY_ID, TOKEN
from parse_eve import HostEventHandler, format_events

EVENTS = HostEventHandler("./events.json")

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)


async def update_events_job(context: ContextTypes.DEFAULT_TYPE):
    EVENTS.update()
    await context.bot.send_message(
        chat_id=MY_ID, text="(Job) Events should be up-to-date. (: "
    )


async def get_help_msg(update: Update, context: ContextTypes.DEFAULT_TYPE):
    help_msg = """*How to use the bot?*
The following list shows the currently available commands with a short description.

- To get the events list run /events and choose the day you want to check.
- To display this help message again, you can run /help.

*What's more to come?*
Here is a list of features I have in mind that will be implemented sooner or later.

- Add a custom event to distribute it to others.
- Prefilter the clubs you wish to get updates about.
- Get updates about the next week.
- Open Issues/Request new clubs/Give Feedback via the bot.

*Feedback/Issues/Requests?*
If you got any of the above you can use one of the following methods:

- Use the bot builtin functionality. (TBD)
- Open up an Issue on [GitHub](https://github.com/mitschix/evenpargo/issues). (if you know what that is and how to use it :D)
- Contact me directly and tell me what's on your mind - @mitschix (:
"""
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=help_msg,
        parse_mode="Markdown",
        disable_web_page_preview=True,
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
    print(q_data)
    q_day = q_data.split("_")[1].title()
    print(q_day)
    event_msg = format_events(EVENTS.get_events_per_day(q_day))
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=event_msg,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    welcome_msg = """Welcome to *EvenParVIE*! (:

This is a tiny little bot that tries to visit the website of a bunch of clubs in vienna to get the latest events of the current weekend. This can be useful to get a brief overview and see where you want to go out.

To get the events - run /events and choose the day you wish to get information about.
To get more information or if you need any help you can run /help.

Feedback is very much appreciated. (:

Have a nice day/night and KEEP RAVING. 😁
- @mitschix"""
    await context.bot.send_message(chat_id=update.effective_chat.id, text=welcome_msg)


async def echo(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=update.effective_chat.id, text="Sorry unrecogniced message :/"
    )


if __name__ == "__main__":
    application = ApplicationBuilder().token(TOKEN).build()

    start_handler = CommandHandler("start", start)
    update_h = CommandHandler("update", update_events)
    help_h = CommandHandler("help", get_help_msg)

    events_get_h = CommandHandler("events", get_events)
    event_show_h = CallbackQueryHandler(handle_events)
    echo_handler = MessageHandler(filters.TEXT, echo)

    application.add_handler(start_handler)
    application.add_handler(update_h)
    application.add_handler(help_h)
    application.add_handler(events_get_h)
    application.add_handler(event_show_h)
    application.add_handler(echo_handler)

    job = application.job_queue  # pip install "python-telegram-bot[job-queue]
    daily_update_eve = job.run_daily(
        update_events_job, datetime.time.fromisoformat("22:00:00+02:00")
    )
    daily_update_mor = job.run_daily(
        update_events_job, datetime.time.fromisoformat("10:00:00+02:00")
    )

    application.run_polling()
