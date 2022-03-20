CREATE TABLE subscription_histories (
    id INTEGER PRIMARY KEY UNIQUE AUTO INCREMENT NOT NULL,
    subscriber_id VARCHAR(255) NOT NULL,
    subscriber_first_name VARCHAR(255) NOT NULL,
    subscriber_last_name VARCHAR(255),
    subscriber_phone_number VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    provider_phone_number VARCHAR(255) NOT NULL,
    provider_email VARCHAR(255),
    subscription_plan_id VARCHAR(255) NOT NULL,
    subscription_plan_name VARCHAR(255) NOT NULL,
    subscription_plan_duration INTEGER NOT NULL,
    subscription_plan_price DOUBLE PRECISION(11, 2) NOT NULL,
    subscription_plan_currency VARCHAR(255) NOT NULL,
    CreatedAt time.Time
);