BEGIN;

CREATE TABLE IF NOT EXISTS "user"(
  id UUID PRIMARY KEY,
  password CHAR(60) NOT NULL,
  email VARCHAR(300) UNIQUE NOT NULL
);

INSERT INTO "user"(id, password, email) VALUES('00000000-0000-0000-0000-000000000000', '$2a$10$Rg.NL4wI72A14m3mApK7t.eqnX57c9a2vnva/.iVzpIkffmHQpG.a', 'admin@crounch.me');

COMMIT;
