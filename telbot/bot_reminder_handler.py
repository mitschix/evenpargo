import datetime
from typing import List, Tuple

import pytz
from telegram import InlineKeyboardMarkup, Update
from telegram.ext import (
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    ConversationHandler,
)

from bot_event_handler import EVENTS, format_events
from bot_keyboards import (
    keyboard_reminder_choice,
    keyboard_reminder_conf,
    keyboard_reminder_days,
)
from config import DB_NAME
from db_reminder import ReminderDB
from telegramtimepicker import TimePicker

timepicker = TimePicker()
reminder_db = ReminderDB(DB_NAME)

REMINDER_CONF, REMINDER_CHOICE, REMINDER_DAY, REMINDER_TIME = range(4)

DAY_MAPPING = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
]


async def reminder_events(context: ContextTypes.DEFAULT_TYPE) -> None:
    fri_events = format_events(EVENTS.get_events_per_day("Friday"))
    sat_events = format_events(EVENTS.get_events_per_day("Saturday"))
    sun_events = format_events(EVENTS.get_events_per_day("Sunday"))
    event_msg = f"""*=== Friday ===*

{fri_events}

*=== Saturday ===*

{sat_events}

*=== Sunday ===*

{sun_events}"""
    await context.bot.send_message(
        chat_id=context.job.chat_id,
        text=event_msg,
        parse_mode="Markdown",
        disable_web_page_preview=True,
    )


def check_if_exist(name: str, context: ContextTypes.DEFAULT_TYPE) -> Tuple[bool, List]:
    """Check if job with given name exists."""
    current_jobs = context.job_queue.get_jobs_by_name(name)
    return bool(current_jobs), current_jobs


def remove_job_if_exists(name: str, context: ContextTypes.DEFAULT_TYPE) -> bool:
    """Remove job with given name. Returns whether job was removed."""

    exists, cur_jobs = check_if_exist(name, context)

    if not exists:
        return False
    for job in cur_jobs:
        job.schedule_removal()
    return True


def set_reminder(
    chat_id: int,
    rem_day: int,
    rem_time: datetime.time,
    context: ContextTypes.DEFAULT_TYPE,
    set_new: bool = True,
) -> None:
    _ = remove_job_if_exists(str(chat_id), context)
    context.job_queue.run_daily(
        reminder_events,
        rem_time,
        days=(rem_day,),
        chat_id=chat_id,
        name=str(chat_id),
    )
    if set_new:
        reminder_db.set_reminder(chat_id, rem_day, str(rem_time))


def set_default(chat_id: int, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Set default reminder (Thursday, 1800)"""

    default_rem_day = 4
    default_rem_time = datetime.time(hour=18, tzinfo=pytz.timezone("Europe/Berlin"))
    set_reminder(chat_id, default_rem_day, default_rem_time, context)


async def handle_time_picker(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    selected, new_time = await timepicker.process_time_selection(update, context)
    if selected:
        if new_time:
            chat_id = update.effective_chat.id
            rem_day = context.user_data.get("day")
            hour = new_time.hour
            minute = new_time.minute
            rem_time = datetime.time(
                hour=hour, minute=minute, tzinfo=pytz.timezone("Europe/Berlin")
            )
            set_reminder(chat_id, rem_day, rem_time, context)

            await context.bot.edit_message_text(
                chat_id=update.effective_chat.id,
                message_id=update.callback_query.message.message_id,
                text=f"🟢 New reminder set to: _{DAY_MAPPING[rem_day]}, {hour:02d}:{minute:02d}_!",
                parse_mode="Markdown",
            )
            return ConversationHandler.END
    return REMINDER_TIME


async def cancel(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    await update.message.reply_text("Reminder settings cancelled.")
    return ConversationHandler.END


async def handle_day(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    chat_id = update.effective_chat.id
    q_data = int(update.callback_query.data)
    if q_data < 0:
        await context.bot.edit_message_text(
            chat_id=chat_id,
            message_id=update.callback_query.message.message_id,
            text="Customization cancelled.",
        )
        return ConversationHandler.END

    context.user_data["day"] = q_data
    await context.bot.edit_message_text(
        chat_id=chat_id,
        message_id=update.callback_query.message.message_id,
        text="🕔 Please select the time you wish to be reminded at:",
        reply_markup=timepicker.create_timepicker(),
    )
    return REMINDER_TIME


async def handle_rem_choice(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    chat_id = update.effective_chat.id
    q_data = str(update.callback_query.data)
    option = q_data.split("_")[1]
    if option == "default":
        set_default(chat_id, context)
        await context.bot.edit_message_text(
            chat_id=chat_id,
            message_id=update.callback_query.message.message_id,
            text="🟢 New reminder set to default: _Thursday, 18:00_!",
            parse_mode="Markdown",
        )
        return ConversationHandler.END
    elif option == "cancel":
        await context.bot.edit_message_text(
            chat_id=chat_id,
            message_id=update.callback_query.message.message_id,
            text="Reminder settings cancelled.",
        )
        return ConversationHandler.END

    await context.bot.edit_message_text(
        chat_id=chat_id,
        message_id=update.callback_query.message.message_id,
        text="📅 Please choose the day you wish to be reminded:",
        reply_markup=InlineKeyboardMarkup(keyboard_reminder_days),
    )
    return REMINDER_DAY


async def handle_reminder(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    chat_id = update.effective_chat.id

    q_data = str(update.callback_query.data)
    option = q_data.split("_")[1]
    text_msg = ""
    reply_markup = None
    if option == "toggle":
        job_removed = remove_job_if_exists(str(chat_id), context)
        if job_removed:
            reminder_db.toggle_state(chat_id)
            text_msg = "❌ Reminder deactivated."
            return_code = ConversationHandler.END
        else:
            current_rem = reminder_db.get_reminder(chat_id)
            if current_rem:
                rem_day = current_rem.get("day", 0)
                rem_time = current_rem.get("time", datetime.time(hour=18))
                set_reminder(
                    current_rem.get("userid", 0),
                    rem_day,
                    rem_time,
                    context,
                    set_new=False,
                )
                exists, _ = check_if_exist(str(chat_id), context)
                if exists:
                    reminder_db.toggle_state(chat_id)
                    text_msg = f"🟢 Reminder activated: _{DAY_MAPPING[rem_day]}, {str(rem_time)}_ !"
                    return_code = ConversationHandler.END
                else:
                    text_msg = "🙈 Something went wrong. Please set a new reminder!\n\nDo you wish to active the *default* reminder (_Thursday, 18:00_) or set a *custom* time?"
                    reply_markup = InlineKeyboardMarkup(keyboard_reminder_choice)
                    return_code = REMINDER_CHOICE
            else:
                text_msg = "🤷 No reminder found.\n\nDo you wish to active the *default* reminder (_Thursday, 18:00_) or set a *custom* time?"
                reply_markup = InlineKeyboardMarkup(keyboard_reminder_choice)
                return_code = REMINDER_CHOICE
    elif option == "change":
        text_msg = "📅 Please choose the day you wish to be reminded:"
        reply_markup = InlineKeyboardMarkup(keyboard_reminder_days)
        return_code = REMINDER_DAY
    elif option == "cancel":
        text_msg = "Reminder settings cancelled."
        return_code = ConversationHandler.END
    elif option == "delete":
        reminder_db.delete_reminder_by_user_id(chat_id)
        text_msg = "🔕 Reminder deleted."
        return_code = ConversationHandler.END
    else:
        text_msg = "Something went wrong. :/"
        return_code = ConversationHandler.END

    await context.bot.edit_message_text(
        chat_id=chat_id,
        message_id=update.callback_query.message.message_id,
        text=text_msg,
        reply_markup=reply_markup,
        parse_mode="Markdown",
    )

    return return_code


async def start_rem_conf(update: Update, context: ContextTypes.DEFAULT_TYPE) -> int:
    chat_id = update.message.chat_id

    curr_rem_info = reminder_db.get_reminder(chat_id)
    if curr_rem_info:
        curr_rem_text = f"\n🗒️ Current reminder: \
_{DAY_MAPPING[curr_rem_info.get('day',0)]}, \
{curr_rem_info.get('time')}_ !"

        exists, _ = check_if_exist(str(chat_id), context)
        state_text, state_icon = ("*ON*", "🟢") if exists else ("*OFF*", "❌")
        state_info = f"{state_icon} Your reminder is currently set to {state_text}.\n"
    else:
        curr_rem_text = "\n🤐 No reminder configured!"
        state_info = ""

    reply_markup = InlineKeyboardMarkup(keyboard_reminder_conf)
    await update.message.reply_text(
        text=f"What do you want to do?\n{state_info}{curr_rem_text}",
        parse_mode="Markdown",
        reply_markup=reply_markup,
    )

    return REMINDER_CONF


rem_handler = ConversationHandler(
    entry_points=[CommandHandler("reminder", start_rem_conf)],
    states={
        REMINDER_CONF: [CallbackQueryHandler(handle_reminder)],
        REMINDER_CHOICE: [CallbackQueryHandler(handle_rem_choice)],
        REMINDER_DAY: [CallbackQueryHandler(handle_day)],
        REMINDER_TIME: [CallbackQueryHandler(handle_time_picker)],
    },
    fallbacks=[CommandHandler("cancel", cancel)],
)
