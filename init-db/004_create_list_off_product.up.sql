BEGIN;

CREATE TABLE IF NOT EXISTS "list_off_product"(
  list_id CHAR(36) NOT NULL,
  code CHAR(13) NOT NULL,
  CONSTRAINT PK_list_off_product PRIMARY KEY (list_id,code),
  CONSTRAINT FK_list_off_product_list_id FOREIGN KEY (list_id) REFERENCES "list"(id)
);

COMMIT;
