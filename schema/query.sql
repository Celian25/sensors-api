-- name: GetDevices :many
SELECT * FROM devices
ORDER BY name;

-- name: InsertDevice :exec
INSERT INTO devices (id, name, location, token)
VALUES (?, ?, ?, ?);

-- name: InsertSensorData :one
INSERT INTO sensor_data (id, sensor_id, temperature, humidity, timestamp)
VALUES (?, ?, ?, ?, ?)
RETURNING *;


-- name: GetSensorDataByTime :many
SELECT * FROM sensor_data s
WHERE s.sensor_id = ?
and s.timestamp >= ?
and s.timestamp < ?;