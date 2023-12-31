from telegram import InlineKeyboardButton

keyboard_days = [
    [InlineKeyboardButton("Friday", callback_data="check_friday")],
    [InlineKeyboardButton("Saturday", callback_data="check_saturday")],
    [InlineKeyboardButton("Sunday", callback_data="check_sunday")],
]
keyboard_report_types = [
    [InlineKeyboardButton("Feedback", callback_data="Feedback")],
    [InlineKeyboardButton("Request", callback_data="Request")],
    [InlineKeyboardButton("Issue", callback_data="Issue")],
]
keyboard_report_user = [
    [InlineKeyboardButton("Sure/OK! ✅", callback_data=1)],
    [InlineKeyboardButton("I would prefer not to. ❌", callback_data=0)],
]
