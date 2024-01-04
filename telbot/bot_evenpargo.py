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
from bot_msg import CLUB_MSG, HELP_MSG, WELCOME_MSG
from bot_report_conv_handler import conv_handler
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
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=HELP_MSG,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


async def get_club_list(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=CLUB_MSG,
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
    await context.bot.send_message(
        chat_id=MY_ID,
        text=f"New user {update.effective_user.first_name} {update.effective_user.last_name} - @{update.effective_user.username} . (:",
    )
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=WELCOME_MSG,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


async def echo(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=update.effective_chat.id, text="Sorry unrecogniced message :/"
    )


if __name__ == "__main__":
    application = ApplicationBuilder().token(TOKEN).build()

    start_handler = CommandHandler("start", start)
    update_h = CommandHandler("update", update_events, filters.User(MY_ID))
    help_h = CommandHandler("help", get_help_msg)
    list_h = CommandHandler("list", get_club_list)

    events_get_h = CommandHandler("events", get_events)
    event_show_h = CallbackQueryHandler(handle_events)
    echo_handler = MessageHandler(filters.TEXT, echo)

    # Add conversation handler with the states GENDER, PHOTO, LOCATION and BIO

    application.add_handler(start_handler)
    application.add_handler(update_h)
    application.add_handler(help_h)
    application.add_handler(list_h)
    application.add_handler(conv_handler)
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

    application.run_polling(allowed_updates=Update.ALL_TYPES)
