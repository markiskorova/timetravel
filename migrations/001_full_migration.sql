-- Rename original table
ALTER TABLE records RENAME TO old_records;

-- Recreate full schema
CREATE TABLE IF NOT EXISTS records (
    id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS record_history (
    record_history_id INTEGER PRIMARY KEY AUTOINCREMENT,
    record_id INTEGER NOT NULL,
    version INTEGER NOT NULL,
    data TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_latest BOOLEAN DEFAULT TRUE,

    FOREIGN KEY (record_id) REFERENCES records(id)
);

CREATE INDEX IF NOT EXISTS idx_record_versions_record_id_version
ON record_history (record_id, version);

CREATE INDEX IF NOT EXISTS idx_record_versions_latest
ON record_history (record_id, is_latest);

CREATE VIEW IF NOT EXISTS record_latest AS
SELECT *
FROM record_history
WHERE is_latest = TRUE;

-- Migrate data
INSERT INTO records (id)
SELECT id FROM old_records;

INSERT INTO record_history (record_id, version, data, is_latest)
SELECT id, 1, data, TRUE FROM old_records;

-- Drop the old table
DROP TABLE old_records;