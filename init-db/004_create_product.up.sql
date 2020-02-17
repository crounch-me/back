BEGIN;

CREATE TABLE IF NOT EXISTS "product"(
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  creation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id CHAR(36) NOT NULL,
  CONSTRAINT FK_product_user_id FOREIGN KEY (user_id) REFERENCES "user"(id)
);

COMMIT;
