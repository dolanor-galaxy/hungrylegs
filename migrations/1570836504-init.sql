-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS activity (
  uuid TEXT,
  suuid TEXT,
  sport TEXT,
  "time" TEXT,
  device TEXT,
  UNIQUE(uuid)
);

CREATE TABLE IF NOT EXISTS trackpoint (
  "time"  TEXT,    -- time is the pk for this table
  lat     FLOAT,
	long    FLOAT,
  alt     FLOAT,
  dist    FLOAT,
  hr      FLOAT,
  cad     FLOAT,
  speed   FLOAT,
  "power" FLOAT,
  activity_uuid TEXT,
  FOREIGN KEY(activity_uuid) REFERENCES activity(uuid)
);

CREATE TABLE IF NOT EXISTS lap (
  "time"      TEXT, -- time is the pk for this table
  start       TEXT,
  total_time  FLOAT,
  dist        FLOAT,
  calories    FLOAT,
  max_speed   FLOAT,
  avg_hr      FLOAT,
  max_hr      FLOAT,
  intensity   TEXT,
  trigger     TEXT,
  activity_uuid TEXT,
  FOREIGN KEY(activity_uuid) REFERENCES activity(uuid)
);

CREATE UNIQUE INDEX idx_act_uuid ON activity(uuid);
CREATE UNIQUE INDEX idx_act_suuid ON activity(suuid);
CREATE INDEX idx_act_sport ON activity(sport);
CREATE INDEX idx_act_time ON activity("time");

CREATE UNIQUE INDEX idx_tp_time ON trackpoint("time");
CREATE UNIQUE INDEX idx_lap_time ON lap("time");

-- Table to know if we've already imported something
CREATE TABLE IF NOT EXISTS fileimport (
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
