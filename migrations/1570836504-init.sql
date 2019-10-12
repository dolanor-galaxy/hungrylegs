-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS TrackPoint (
  id INTEGER PRIMARY KEY,
  `time`  TEXT,
  lat     FLOAT,
	long    FLOAT,
  alt     FLOAT,
  dist    FLOAT,
  hr      FLOAT,
  cad     FLOAT,
  speed   FLOAT,
  `power` FLOAT,
  lap_id INTEGER,
  FOREIGN KEY(lap_id) REFERENCES Lap(id)
);

CREATE TABLE IF NOT EXISTS Lap (
  id INTEGER PRIMARY KEY,
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
  FOREIGN KEY(activity_id) REFERENCES Activity(id)
);

CREATE TABLE IF NOT EXISTS Activity (
  id INTEGER PRIMARY KEY,
  sport TEXT,
  `time` TEXT,
  device TEXT
);

CREATE INDEX idx_tp_lap ON TrackPoint(lap_id);
CREATE INDEX idx_lap_act ON Lap(activity_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE Trackpoint;
DROP TABLE Lap;
DROP TABLE Activity;

-- DROP INDEX salary_index;
