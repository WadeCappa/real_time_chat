create database chat_db;

\c chat_db

create table if not exists messages (
    user_id bigint,
    message_id bigint,
    channel_id bigint,
    time_posted timestamp,
    content text,

    primary key(channel_id, message_id)
);

create index idx_ordered_messages_by_channel on messages (channel_id, message_id desc);

create sequence message_id_sequence start with 101;