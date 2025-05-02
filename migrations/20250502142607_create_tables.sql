-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS jobs_raw (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    source_link VARCHAR(2048) NOT NULL,
    main_technology VARCHAR(255),
    date_posted TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    date_parsed TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_raw_main_technology ON jobs_raw(main_technology);

CREATE TABLE IF NOT EXISTS telegram_channels (
    id BIGSERIAL PRIMARY KEY,
    tag VARCHAR(255) NOT NULL,
    last_post_id BIGINT NOT NULL DEFAULT 0,
    date_channel_added TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    posts_parsed BIGINT NOT NULL DEFAULT 0,
    date_last_parsed TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_last_post FOREIGN KEY (last_post_id) REFERENCES jobs_raw(id)
);

CREATE INDEX idx_telegram_channels_tag ON telegram_channels(tag);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS telegram_channels;
DROP TABLE IF EXISTS jobs_raw;
-- +goose StatementEnd
