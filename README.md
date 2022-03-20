# onemembership
It is a telegram membership bot that is used for controlling telegram channels and groups.

User
- CRUD on profile { view, edit or disable/remove profile }
- CRUD on subscriptions { add, view, change or remove subscription }
- CRUD on subscribed to parties { add, view, edit or remove channel or group }
- User just pay using given payment methods
- Change language

Logging
- ServerLogFile contains log of [ Preference, Feedback, Language, Language Entries ]
- BotLogFile contains log of [ Temporary Service Provider, Temporary User ]

ATTENTION
- The tables languages and language_entries contain emoji's special characters therefore
    * It respective special character containing fields must be 'BLOB' type
    * The database must have character set of utf8mb4 ( ALTER DATABASE onemembership CHARACTER SET = utf8mb4    COLLATE = utf8mb4_unicode_ci; )
    * The table must have character set of utf8mb4 ( ALTER TABLE languages CONVERT TO CHARACTER SET utf8mb4        COLLATE utf8mb4_unicode_ci; ALTER TABLE languages MODIFY flag BLOB CHARSET utf8mb4; )