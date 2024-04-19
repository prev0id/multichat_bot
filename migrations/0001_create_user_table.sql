-- +goose Up
create table if not exists user_data
(
    key    text unique not null,
    twitch text unique not null,
    youtube text unique not null
);

-- +goose Down
drop table user;
