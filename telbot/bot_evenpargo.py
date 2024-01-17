#!/usr/bin/env python
"""Telegram bot to read events.json and returns events"""
import logging

from telegram import Update
from telegram.ext import (
    ApplicationBuilder,
    CommandHandler,
    ContextTypes,
    MessageHandler,
    filters,
)

from bot_event_handler import event_get_h, event_show_h, event_update_h, set_update_jobs
from bot_msg import CLUB_MSG, HELP_MSG, WELCOME_MSG
from bot_reminder_handler import rem_handler, set_default
from bot_report_conv_handler import conv_handler
from config import SUPPORT_ID, TOKEN

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)


async def get_club_list(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=CLUB_MSG,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=SUPPORT_ID,
        text=f"New user {update.effective_user.first_name} {update.effective_user.last_name} - @{update.effective_user.username} . (:",
    )
    set_default(update.effective_chat.id, context)
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=WELCOME_MSG,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


async def get_help_msg(update: Update, context: ContextTypes.DEFAULT_TYPE):
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text=HELP_MSG,
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
    help_h = CommandHandler("help", get_help_msg)
    list_h = CommandHandler("list", get_club_list)

    echo_handler = MessageHandler(filters.TEXT, echo)

    application.add_handler(start_handler)
    application.add_handler(help_h)
    application.add_handler(list_h)

    application.add_handler(conv_handler)

    application.add_handler(rem_handler)

    application.add_handler(event_update_h)
    application.add_handler(event_get_h)
    application.add_handler(event_show_h)

    application.add_handler(echo_handler)

    jobq = application.job_queue  # pip install "python-telegram-bot[job-queue]
    set_update_jobs(jobq)

    application.run_polling(allowed_updates=Update.ALL_TYPES)
