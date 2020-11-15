BEGIN;

CREATE TABLE IF NOT EXISTS "category"(
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL
);

INSERT INTO "category"(id, name) VALUES
('5bf5983c-2694-4467-b755-6c43aa870a34', 'Boucherie'),
('871f95e1-0719-4d45-a011-6c2331fa4fd3', 'Epicerie')
;

COMMIT;
