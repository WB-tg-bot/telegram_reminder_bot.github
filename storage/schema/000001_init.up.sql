CREATE TABLE tasks(
                      id SERIAL PRIMARY KEY,
                      content TEXT NOT NULL,
                      reminder_time TIMESTAMP NOT NULL,
                      chat_id BIGINT NOT NULL,
                      username TEXT NOT NULL
);
