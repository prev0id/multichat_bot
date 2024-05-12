create table if not exists user
(
    id integer primary key not null
);

create table if not exists platform
(
    user_id        integer not null,

    name           text    not null,
    id             text    not null,
    channel        text    not null,

    access_token   text    not null,
    refresh_token  text    not null,
    expires_in     text    not null,

    disabled_users text    not null,
    banned_words   text    not null
);

