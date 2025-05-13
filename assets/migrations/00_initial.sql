-- +migrate Up

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    city VARCHAR(255) NOT NULL,
    -- can be bad in some cases, but for now it is strictly defined
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('daily', 'hourly')),
    confirmed BOOLEAN DEFAULT FALSE,
    unsubscribe_token VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_notified_at TIMESTAMP WITH TIME ZONE,
);

CREATE TABLE IF NOT EXISTS confirmation_tokens (
    token VARCHAR(64) PRIMARY KEY,
    subscription_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);

-- Here some useful triggers can be added. But, for this task, it will be better to manually
-- execute required actions (like token removal). It would be much easier to read and understand
-- the written code without keeping the triggers in mind (and thinking some logic is simply missed)

-- +migrate Down

DROP TABLE IF EXISTS confirmation_tokens;
DROP TABLE IF EXISTS subscriptions;