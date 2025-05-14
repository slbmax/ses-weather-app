-- +migrate Up

-- simplifying to use only one table for all subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL, -- RFC 5321 and RFC 5322
    city VARCHAR(100) NOT NULL,
    -- can be bad in some cases, but for now it is strictly defined
    frequency VARCHAR(6) NOT NULL CHECK (frequency IN ('daily', 'hourly')),

    confirmed BOOLEAN DEFAULT FALSE,
    -- From the spec, 409 error when confirming a subscription is not defined, meaning the token will be changed once used.
    -- This means we need to have two tokens: for confirmation and unsubscription. But there is no need to issue both at the same time.
    -- This means we can use the same field for both tokens, and just update the token to unsubscription when the confirmation token is used.
    -- Also, as there isn't any confirmation email resending defined by the spec, we can simplify the process by using a token without expiration.
    token VARCHAR(32) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_notified_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT unique_email UNIQUE (email)
);


-- +migrate Down

DROP TABLE IF EXISTS subscriptions;