BEGIN;

ALTER TABLE "product_in_list" ADD COLUMN bought BOOLEAN DEFAULT FALSE;

COMMIT;
