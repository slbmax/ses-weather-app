-- +migrate Up

-- simplifying to use one table for confirmed and unconfirmed subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL, -- RFC 5321 and RFC 5322
    city VARCHAR(100) NOT NULL,
    -- can be bad in some cases, but for now it is strictly defined
    frequency VARCHAR(6) NOT NULL CHECK (frequency IN ('daily', 'hourly')),

    -- as there isn't any confirmation email resending defined by the spec,
    -- we can simplify the process by using a token without expiration
    -- and reuse it for the unsubscribing process. Also, using
    -- uuid for the token to avoid collisions
    token VARCHAR(36) NOT NULL,
    confirmed BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_notified_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_email UNIQUE (email)
);

-- +migrate Down

DROP TABLE IF EXISTS subscriptions;