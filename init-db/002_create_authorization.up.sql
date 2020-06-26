BEGIN;

CREATE TABLE IF NOT EXISTS "authorization"(
  user_id UUID NOT NULL,
  token CHAR(36) NOT NULL,
  CONSTRAINT PK_authorization PRIMARY KEY (user_id,token),
  CONSTRAINT FK_authorization_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

COMMIT;
