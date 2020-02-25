BEGIN;

CREATE TABLE IF NOT EXISTS "product_in_list"(
  product_id CHAR(36) NOT NULL,
  list_id CHAR(36) NOT NULL,
  CONSTRAINT PK_product_in_list PRIMARY KEY (product_id,list_id),
  CONSTRAINT FK_product_in_list_product_id FOREIGN KEY (product_id) REFERENCES "product"(id),
  CONSTRAINT FK_product_in_list_list_id FOREIGN KEY (list_id) REFERENCES "list"(id)
);

COMMIT;
