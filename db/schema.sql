create table user_post (
    post_id BIGINT,
    time_posted BIGINT,
    content TEXT,
    primary key (post_id)
);
create sequence post_id_sequence start 1;