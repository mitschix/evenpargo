from telegram import InlineKeyboardMarkup, Update
from telegram.ext import (
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    ConversationHandler,
    MessageHandler,
    filters,
)

from bot_keyboards import keyboard_report_types, keyboard_report_user
from config import SUPPORT_ID

REPORT_TYPE, REPORT_USER, REPORT_INFO = range(3)


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
    await context.bot.edit_message_text(
        chat_id=update.effective_chat.id,
        message_id=update.callback_query.message.message_id,
        text=f"Please tell me what's on your mind? ({rep_type}) 🤔",
    )

    return REPORT_USER


async def get_rep_user(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    rep_type = context.user_data.get("type")
    context.user_data["text"] = update.message.text
    reply_markup = InlineKeyboardMarkup(keyboard_report_user)
    await context.bot.send_message(
        text=f"Thank you for your message! I hope I can help you with this {rep_type}. (:\n\n"
        + f"If you don't mind, I'll send your username along with the {rep_type} - in case I have further questions - OK?",
        chat_id=update.effective_chat.id,
        reply_markup=reply_markup,
    )
    return REPORT_INFO


async def send_report(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    rep_type = context.user_data.get("type")
    rep_text = context.user_data.get("text")
    user_check = bool(int(update.callback_query.data))
    rep_msg = f"""New Report! (:

Topic: {rep_type}

Text:
{rep_text}

{f'User: @{update.effective_user.username}' if user_check else ''}
#{rep_type.lower()}"""

    await context.bot.send_message(
        chat_id=SUPPORT_ID,
        text=rep_msg,
    )

    await context.bot.edit_message_text(
        chat_id=update.effective_chat.id,
        message_id=update.callback_query.message.message_id,
        text=f"Thank you! (:\n\n*{rep_type}* sent with the following text:\n_{rep_text}_",
        parse_mode="Markdown",
    )

    return ConversationHandler.END


async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    await update.message.reply_text("Reporting cancelled.")
    return ConversationHandler.END


conv_handler = ConversationHandler(
    entry_points=[CommandHandler("report", submit_report)],
    states={
        REPORT_TYPE: [CallbackQueryHandler(get_rep_type)],
        REPORT_USER: [MessageHandler(filters.TEXT, get_rep_user)],
        REPORT_INFO: [CallbackQueryHandler(send_report)],
    },
    fallbacks=[CommandHandler("cancel", cancel)],
)
