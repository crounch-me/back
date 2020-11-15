BEGIN;

CREATE TABLE IF NOT EXISTS "contributor"(
  list_id UUID,
  user_id UUID,
  CONSTRAINT PK_contributor PRIMARY KEY (list_id, user_id),
  CONSTRAINT FK_contributor_list_id FOREIGN KEY (list_id) REFERENCES "list"(id) ON DELETE CASCADE,
  CONSTRAINT FK_contributor_user_id FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);

COMMIT;
