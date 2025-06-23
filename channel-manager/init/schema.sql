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