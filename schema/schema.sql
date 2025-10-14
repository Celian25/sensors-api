CREATE TABLE IF NOT EXISTS devices (
                                       id TEXT PRIMARY KEY,
                                       name TEXT NOT NULL,
                                       location TEXT NOT NULL,
                                       token TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sensor_data (
                                           id TEXT PRIMARY KEY,
                                           sensor_id TEXT NOT NULL,
                                           temperature REAL NOT NULL,
                                           humidity REAL NOT NULL,
                                           timestamp TEXT NOT NULL,
                                           FOREIGN KEY (sensor_id) REFERENCES devices(id)
    ON UPDATE CASCADE ON DELETE CASCADE
);