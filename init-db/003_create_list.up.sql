BEGIN;

CREATE TABLE IF NOT EXISTS "list"(
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE NOT NULL,
  user_id CHAR(36) NOT NULL,
  CONSTRAINT FK_list_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

COMMIT;
