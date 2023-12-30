from telegram import InlineKeyboardMarkup, Update
from telegram.ext import (
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    ConversationHandler,
    MessageHandler,
    filters,
)

from bot_keyboards import keyboard_report_types
from config import MY_ID

REPORT_TYPE, REPORT_INFO = range(2)


async def submit_report(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    reply_markup = InlineKeyboardMarkup(keyboard_report_types)
    await update.message.reply_text(
        text="Please choose the type of report?",
        reply_markup=reply_markup,
    )

    return REPORT_TYPE


async def get_rep_type(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    rep_type = update.callback_query.data
    context.user_data["type"] = rep_type
    await context.bot.send_message(
        chat_id=update.effective_chat.id,
        text="Please tell me what's up?",
    )

    return REPORT_INFO


async def send_report(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    await update.message.reply_text("Thank you! I hope I can help you. (:")

    await context.bot.send_message(
        chat_id=MY_ID,
        text=f"New Issue reported from: @{update.effective_user.username} . (:\n\nTopic: {context.user_data.get('type')}\n\nText:\n{update.message.text}",
    )

    return ConversationHandler.END


async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    await update.message.reply_text("Reporting cancelled.")
    return ConversationHandler.END


conv_handler = ConversationHandler(
    entry_points=[CommandHandler("report", submit_report)],
    states={
        REPORT_TYPE: [CallbackQueryHandler(get_rep_type)],
        REPORT_INFO: [MessageHandler(filters.TEXT, send_report)],
    },
    fallbacks=[CommandHandler("cancel", cancel)],
)
