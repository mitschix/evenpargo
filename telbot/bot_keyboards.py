from telegram import InlineKeyboardButton

keyboard_days = [
    [
        InlineKeyboardButton("Friday", callback_data="check_friday"),
        InlineKeyboardButton("Saturday", callback_data="check_saturday"),
        InlineKeyboardButton("Sunday", callback_data="check_sunday"),
    ],
]
keyboard_report_types = [
    [
        InlineKeyboardButton("Feedback", callback_data="Feedback"),
        InlineKeyboardButton("Request", callback_data="Request"),
        InlineKeyboardButton("Issue", callback_data="Issue"),
    ],
]
keyboard_report_user = [
    [InlineKeyboardButton("✅ Sure/OK!", callback_data=1)],
    [InlineKeyboardButton("❌ I would prefer not to.", callback_data=0)],
]

keyboard_reminder_conf = [
    [
        InlineKeyboardButton("Toggle", callback_data="rem_toggle"),
        InlineKeyboardButton("Change", callback_data="rem_change"),
    ],
]
keyboard_reminder_choice = [
    [
        InlineKeyboardButton("Default", callback_data="rem_default"),
        InlineKeyboardButton("Custom", callback_data="rem_custom"),
    ],
]


keyboard_reminder_days = [
    [
        InlineKeyboardButton("Monday", callback_data=1),
        InlineKeyboardButton("Tuesday", callback_data=2),
    ],
    [
        InlineKeyboardButton("Wednesday", callback_data=3),
        InlineKeyboardButton("Thursday", callback_data=4),
    ],
    [
        InlineKeyboardButton("Friday", callback_data=5),
        InlineKeyboardButton("Saturday", callback_data=6),
    ],
    [
        InlineKeyboardButton("Sunday", callback_data=0),
        InlineKeyboardButton("Cancel", callback_data=-1),
    ],
]
