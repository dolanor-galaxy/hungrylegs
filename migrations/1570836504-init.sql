-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS activity (
  id INTEGER PRIMARY KEY,
  uuid TEXT,
  full_uuid TEXT,
  sport TEXT,
  "time" TEXT,
  device TEXT
);

CREATE TABLE IF NOT EXISTS trackpoint (
  id INTEGER PRIMARY KEY,
  "time"  TEXT,
  lat     FLOAT,
	long    FLOAT,
  alt     FLOAT,
  dist    FLOAT,
  hr      FLOAT,
  cad     FLOAT,
  speed   FLOAT,
  "power" FLOAT,
  activity_id INTEGER,
  FOREIGN KEY(activity_id) REFERENCES activity(id)
);

CREATE TABLE IF NOT EXISTS Lap (
  id INTEGER PRIMARY KEY,
  "time"      TEXT,
  start       TEXT,
  total_time  FLOAT,
  dist        FLOAT,
  calories    FLOAT,
  max_speed   FLOAT,
  avg_hr      FLOAT,
  max_hr      FLOAT,
  intensity   TEXT,
  trigger     TEXT,
  activity_id INTEGER,
  FOREIGN KEY(activity_id) REFERENCES activity(id)
);

CREATE INDEX idx_uuid ON activity(uuid);
CREATE UNIQUE INDEX idx_full_uuid ON activity(full_uuid);
CREATE INDEX idx_sport ON activity(sport);

-- Table to know if we've already imported something
CREATE TABLE IF NOT EXISTS fileimport (
  id INTEGER PRIMARY KEY,
  import_time  TEXT,
  "file_name"  TEXT
);
CREATE INDEX idx_file_name ON fileimport("file_name");

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE trackpoint;
DROP TABLE lap;
DROP TABLE activity;
DROP TABLE fileimport;

-- DROP INDEX salary_index;
