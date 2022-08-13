drop table if exists posts;
drop table if exists threads;
drop table if exists sessions;
drop table if exists users;
drop table if exists likedPosts;
drop table if exists likedThreads;

create table users (
  id          serial primary key,
  uuid        varchar(64) not null unique,
  name        varchar(255),
  email       varchar(255) not null unique,
  password    varchar(255) not null,
  created_at  timestamp not null
);

create table sessions (
  id          serial primary key,
  uuid        varchar(64) not null unique,
  email       varchar(255) not null unique,
  user_id     integer references users(id),
  created_at  timestamp not null
);

create table threads (
  id          serial primary key,
  uuid        varchar(64) not null unique,
  topic       text,
  user_id     integer references users(id),
  author      varchar(255),
  created_at  timestamp not null,
  likes       integer
);

create table posts (
  id          serial primary key,
  uuid        varchar(64) not null unique,
  body        text,
  user_id     integer references users(id),
  thread_id   integer references threads(id),
  created_at  timestamp not null,
  likes       integer
);

create table likedThreads (
  user_uuid     varchar(64),
  thread_uuid   varchar(64),
  primary key   (user_uuid, thread_uuid)
);

create table likedPosts (
  user_uuid       varchar(64),
  post_uuid       varchar(64),
  primary key     (user_uuid, post_uuid)
);