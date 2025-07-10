package service

import (
	"context"
	"errors"

	"github.com/rainbowmga/timetravel/entity"
)

var (
	ErrRecordDoesNotExist    = errors.New("record with that id does not exist")
	ErrRecordIDInvalid       = errors.New("record id must >= 0")
	ErrRecordAlreadyExists   = errors.New("record already exists")
	ErrFailedToReadData      = errors.New("failed to read data")
	ErrFailedToOpenDatabase  = errors.New("failed to open SQLite database")
	ErrFailedToReadSchema    = errors.New("failed to read schema file")
	ErrFailedToExecuteSchema = errors.New("failed to execute schema SQL")
)

// Implements method to get, create, and update record data.
type RecordService interface {

	// GetRecord will retrieve an record.
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// CreateRecord will insert a new record.
	//
	// If it a record with that id already exists it will fail.
	CreateRecord(ctx context.Context, record entity.Record) error

	// UpdateRecord will update the `Map` values of the record if they exist.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, record entity.Record) error

	ListRecordVersions(ctx context.Context, id int) ([]entity.Record, error)

	GetRecordVersion(ctx context.Context, id int, version int) (entity.Record, error)
}
