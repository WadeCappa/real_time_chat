
-- for the channel manager
create database channel_manager_db;

\c channel_manager_db

create table channels (
    id bigint,
    name text,
    public boolean,
    primary key (id)
);

create sequence channel_ids start with 101;

create index idx_public_channels on channels (id) where public=True;

create table channel_members (
    channel_id bigint,
    user_id bigint,
    foreign key (channel_id) references channels(id),
    primary key (channel_id, user_id)
);

create index idx_channel_members_by_channel_id on channel_members (channel_id);

-- for chat-db (the chat writer)
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