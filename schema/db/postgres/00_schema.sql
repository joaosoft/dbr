-- migrate up

CREATE SCHEMA IF NOT EXISTS dbr;

CREATE TABLE dbr.address (
  id_address    SERIAL PRIMARY KEY,
  street        TEXT,
  number        INTEGER,
  country       TEXT
);


CREATE TABLE dbr.person (
  id_person     SERIAL PRIMARY KEY,
  first_name    TEXT,
  last_name     TEXT,
  age           INTEGER,
  active        BOOLEAN,
  fk_address    INTEGER REFERENCES address (id_address)
);

-- migrate down