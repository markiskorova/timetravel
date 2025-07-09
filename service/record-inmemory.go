package service

import (
	"context"

	"github.com/rainbowmga/timetravel/entity"
)

// InMemoryRecordService is an in-memory implementation of RecordService.
type InMemoryRecordService struct {
	data map[int]entity.Record
}

func NewInMemoryRecordService() InMemoryRecordService {
	return InMemoryRecordService{
		data: map[int]entity.Record{},
	}
}

func (s *InMemoryRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	record := s.data[id]
	if record.ID == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	record = record.Copy() // copy is necessary so modifations to the record don't change the stored record
	return record, nil
}

func (s *InMemoryRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	existingRecord := s.data[id]
	if existingRecord.ID != 0 {
		return ErrRecordAlreadyExists
	}

	s.data[id] = record
	return nil
}

func (s *InMemoryRecordService) UpdateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	s.data[id] = record

	return nil
}
