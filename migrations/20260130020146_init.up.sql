CREATE TABLE node(
    id TEXT PRIMARY KEY,
    address TEXT,
    load_percent INTEGER
);

CREATE TABLE "user"(
    id UUID PRIMARY KEY,
    node_id TEXT REFERENCES node,
    tariff_id INTEGER,
    traffic_used BIGINT,
    subscription_expires_at TIMESTAMPTZ,
    subscription_is_active BOOLEAN,
    created_at TIMESTAMPTZ
);

CREATE TABLE identity(
    provider TEXT,
    provider_id TEXT,
    user_id UUID REFERENCES "user",
    PRIMARY KEY (provider, provider_id)
);