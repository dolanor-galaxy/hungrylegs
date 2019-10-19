-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS athlete (
  "name" TEXT,
  age INTEGER,
  "weight" INTEGER,
  vo2_max INTEGER,
  rest_hr INTEGER,

  run_max_hr INTEGER,
  run_ft_pace FLOAT,      -- functional threshold pace
  run_zone5 INTEGER,  -- VO2 max
  run_zone4 INTEGER,  -- Threshold
  run_zone3 INTEGER,  -- Tempo
  run_zone2 INTEGER,  -- Endurance
  run_zone1 INTEGER,  -- Active Recovery

  bike_max_hr INTEGER,
  bike_ft_power FLOAT,     -- functional threshold power
  --                                ft_hr *  x  / 100 
  bike_zone5 INTEGER,  -- VO2 max          (121)
  bike_zone4 INTEGER,  -- Threshold        (105)
  bike_zone3 INTEGER,  -- Tempo            (94)
  bike_zone2 INTEGER,  -- Endurance        (83)
  bike_zone1 INTEGER,  -- Active Recovery  (68)

  swim_max_hr INTEGER,
  swim_ft_pace FLOAT,      -- functional threshold pace (meters/min)
  swim_zone5 INTEGER,
  swim_zone4 INTEGER,
  swim_zone3 INTEGER,
  swim_zone2 INTEGER,
  swim_zone1 INTEGER,

  ftp_watts FLOAT,  -- functional threshold watts
  --                                        ftp_watts *  x  / 100
  watts_zone7 INTEGER, -- Anaerobic capacity           (150)
  watts_zone6 INTEGER, -- VO2 max                      (120)
  watts_zone5 INTEGER, -- Lactate Threshold            (105)
  watts_zone4 INTEGER, -- Tempo                        (93)
  watts_zone3 INTEGER, -- Sweet Spot                   (90)
  watts_zone2 INTEGER, -- Endurance                    (76)
  watts_zone1 INTEGER  -- Active Recovery              (55)
  -- power_to_weight FLOAT,  -- ftp_watts / weight * 100 / 100
);


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
DROP TABLE athlete;

-- DROP INDEX salary_index;
