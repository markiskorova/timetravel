package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"github.com/rainbowmga/timetravel/entity"

	_ "modernc.org/sqlite"
)

// SqLiteRecordService is an SqLite implementation of RecordService.
type SqLiteRecordService struct {
	db *sql.DB
}

func NewSqLiteRecordService(dbPath string) (SqLiteRecordService, error) {
	db, err := InitSQLiteDB(dbPath, "migrations/schema.sql")
	if err != nil {
		return SqLiteRecordService{}, err
	}
	return SqLiteRecordService{db: db}, nil
}

func InitSQLiteDB(dbPath string, schemaPath string) (*sql.DB, error) {
	_, err := os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, ErrFailedToOpenDatabase
	}

	if !dbExists {
		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, ErrFailedToReadSchema
		}

		// Create and initailize database from schema
		_, err = db.Exec(string(schema))
		if err != nil {
			return nil, ErrFailedToExecuteSchema
		}
	}

	return db, err
}

func (s *SqLiteRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	row := s.db.QueryRow("SELECT version, created_at, data FROM record_latest WHERE record_id = ?", id)
	var record entity.Record
	var version int
	var dataJSON, createdAt string
	var data map[string]string

	err := row.Scan(&version, &createdAt, &dataJSON)
	if err != nil {
		return record, ErrRecordDoesNotExist
	}

	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		return record, ErrFailedToReadData
	}

	record.Version = version
	record.Timestamp = createdAt
	record.Data = data

	return record, nil
}

func (s *SqLiteRecordService) GetRecordVersion(ctx context.Context, id int, version int) (entity.Record, error) {
	row := s.db.QueryRow("SELECT created_at, data FROM record_history WHERE record_id = ? AND version = ?", id, version)
	var record entity.Record
	var dataJSON, createdAt string
	var data map[string]string

	err := row.Scan(&dataJSON, &createdAt)
	if err != nil {
		return record, ErrRecordDoesNotExist
	}

	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		return record, ErrFailedToReadData
	}

	record.ID = id
	record.Version = version
	record.Timestamp = createdAt
	record.Data = data

	return record, nil
}

func (s *SqLiteRecordService) ListRecordVersions(ctx context.Context, id int) ([]entity.Record, error) {
	rows, err := s.db.Query("SELECT version, created_at, data FROM record_history WHERE record_id = ? ORDER BY version DESC ", id)
	if err != nil {
		return nil, err
	}

	var records []entity.Record
	for rows.Next() {
		var record entity.Record
		var version int
		var dataJSON, createdAt string
		var data map[string]string

		err := rows.Scan(&version, &createdAt, &dataJSON)
		if err != nil {
			return records, ErrRecordDoesNotExist
		}

		err = json.Unmarshal([]byte(dataJSON), &data)
		if err != nil {
			return records, ErrFailedToReadData
		}

		record.ID = id
		record.Version = version
		record.Timestamp = createdAt
		record.Data = data

		records = append(records, record)
	}

	return records, nil
}

func (s *SqLiteRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}
	dataJSON, _ := json.Marshal(record.Data)

	_, err := s.db.Exec("INSERT OR IGNORE INTO records (id) VALUES (?)", id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO record_history (record_id, version, data, is_latest) VALUES (?, ?, ?, TRUE)", id, 1, dataJSON)

	return err
}

func (s *SqLiteRecordService) UpdateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}
	dataJSON, _ := json.Marshal(record.Data)
	var version int

	s.db.QueryRow("SELECT IFNULL(MAX(version), 0) FROM record_history WHERE record_id = ?", id).Scan(&version)

	_, err := s.db.Exec("INSERT INTO record_history (record_id, version, data, is_latest) VALUES (?, ?, ?, TRUE)", id, version+1, dataJSON)

	return err
}
