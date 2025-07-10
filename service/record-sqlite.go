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
	row := s.db.QueryRow("SELECT id, data FROM records WHERE id = ?", id)
	var record entity.Record
	var dataJSON string
	err := row.Scan(&record.ID, &dataJSON)
	if err != nil {
		return record, ErrRecordDoesNotExist
	}
	err = json.Unmarshal([]byte(dataJSON), &record.Data)
	if err != nil {
		return record, ErrFailedToReadData
	}
	return record, nil
}

func (s *SqLiteRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	dataJSON, _ := json.Marshal(record.Data)
	_, err := s.db.Exec("INSERT INTO records (id, data) VALUES (?, ?)", id, dataJSON)

	return err
}

func (s *SqLiteRecordService) UpdateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	dataJSON, _ := json.Marshal(record.Data)
	_, err := s.db.Exec("UPDATE records SET data=? WHERE id=", dataJSON, id)

	return err
}
