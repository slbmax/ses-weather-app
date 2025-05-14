-- +migrate Up

-- simplifying to use one table for confirmed and unconfirmed subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL, -- RFC 5321 and RFC 5322
    city VARCHAR(100) NOT NULL,
    -- can be bad in some cases, but for now it is strictly defined
    frequency VARCHAR(6) NOT NULL CHECK (frequency IN ('daily', 'hourly')),

    confirmed BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_notified_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_email UNIQUE (email)
);


-- from the spec, 409 error when confirming a subscription is not defined, meaning the token will be deleted once used.
-- this means we need to have two tokens: for confirmation and unsubscription.
CREATE TABLE IF NOT EXISTS tokens (
    -- as there isn't any confirmation email resending defined by the spec,
    -- we can simplify the process by using a token without expiration.
    token VARCHAR(32) NOT NULL,
    subscription_id BIGINT NOT NULL,
    is_confirmation BOOLEAN NOT NULL,

    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    CONSTRAINT unique_token UNIQUE (token)
);

-- +migrate Down

DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS subscriptions;