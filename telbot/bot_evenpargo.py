#!/usr/bin/env python
"""Telegram bot to read events.json and returns events"""
import logging

from telegram import Update
from telegram.ext import filters, ApplicationBuilder, ContextTypes, CommandHandler, MessageHandler

from config import TOKEN
from parse_eve import HostEventHandler, format_events

EVENTS = HostEventHandler("../events.json")

logging.basicConfig(
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    level=logging.INFO
)


async def get_fri(update: Update, context: ContextTypes.DEFAULT_TYPE):
    event_msg = format_events(EVENTS.get_events_per_day("Friday"))
    await context.bot.send_message(chat_id=update.effective_chat.id, text=event_msg)


async def get_sat(update: Update, context: ContextTypes.DEFAULT_TYPE):
    event_msg = format_events(EVENTS.get_events_per_day("Saturday"))
    await context.bot.send_message(chat_id=update.effective_chat.id, text=event_msg)


async def get_sun(update: Update, context: ContextTypes.DEFAULT_TYPE):
    event_msg = format_events(EVENTS.get_events_per_day("Sunday"))
    await context.bot.send_message(chat_id=update.effective_chat.id, text=event_msg)


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(chat_id=update.effective_chat.id,
                                   text="I'm a bot, please talk to me!")


async def echo(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(chat_id=update.effective_chat.id,
                                   text="Sorry unrecogniced message :/")


if __name__ == '__main__':
    application = ApplicationBuilder().token(TOKEN).build()

    start_handler = CommandHandler('start', start)
    fri_handler = CommandHandler('event_fri', get_fri)
    sat_handler = CommandHandler('event_sat', get_sat)
    sun_handler = CommandHandler('event_sun', get_sun)
    echo_handler = MessageHandler(filters.TEXT, echo)

    application.add_handler(start_handler)
    application.add_handler(fri_handler)
    application.add_handler(sat_handler)
    application.add_handler(sun_handler)
    application.add_handler(echo_handler)

    application.run_polling()
