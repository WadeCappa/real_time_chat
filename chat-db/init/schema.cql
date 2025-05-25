
create keyspace if not exists posts_db with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };

create table if not exists posts_db.messages (
    userId bigint,
    offset bigint,
    channelId bigint, 
    time_posted timestamp,
    content text,

    primary key(channelId, offset)
)  with clustering order by (offset desc);