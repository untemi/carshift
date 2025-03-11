CREATE TABLE users (
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  username  TEXT    NOT NULL,

  firstname TEXT    NOT NULL,
  lastname  TEXT    NOT NULL,
  passhash  TEXT    NOT NULL,
  
  phone     TEXT    NOT NULL,
  email     TEXT    NOT NULL
);

CREATE TABLE cars (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,

  name  TEXT    NOT NULL,
  price REAL    NOT NULL,

  start_at  DATETIME,
  end_at    DATETIME,

  owner_id    INTEGER NOT NULL,
  district_id INTEGER NOT NULL,

  FOREIGN KEY(owner_id)    REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(district_id) REFERENCES districts(id) ON DELETE CASCADE
);

CREATE TABLE districts (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT    NOT NULL
);
