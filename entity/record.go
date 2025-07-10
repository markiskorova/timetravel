package entity

type Record struct {
	ID        int               `json:"id"`
	Version   int               `json:"version"`
	Timestamp string            `json:"timestamp"` // or time.Time
	Data      map[string]string `json:"data"`
}

func (r *Record) Copy() Record {
	newMap := make(map[string]string, len(r.Data))
	for key, value := range r.Data {
		newMap[key] = value
	}

	return Record{
		ID:        r.ID,
		Version:   r.Version,
		Timestamp: r.Timestamp,
		Data:      newMap,
	}
}
