from telegram import InlineKeyboardButton

keyboard_days = [
    [InlineKeyboardButton("Friday", callback_data="check_friday")],
    [InlineKeyboardButton("Saturday", callback_data="check_saturday")],
    [InlineKeyboardButton("Sunday", callback_data="check_sunday")],
]
