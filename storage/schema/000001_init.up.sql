CREATE TABLE chats (
                       id serial not null unique,
                       name TEXT NOT NULL
);

CREATE TABLE jobs (
                      id serial not null unique,
                      task TEXT NOT NULL,
                      reminder_time TIMESTAMP,
                      done boolean not null default false
);

CREATE TABLE members (
                         id serial not null unique,
                         name TEXT NOT NULL
);

CREATE TABLE jobs_members (
                              id serial not null unique,
                              job_id int references jobs (id) on delete cascade not null,
                              member_id int references members (id) on delete cascade not null
);