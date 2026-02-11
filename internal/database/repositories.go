package database

type EntryRecordRepository interface {
	Save(entry *EntryRecord) error
	Delete(key string) error
	Get(key string) (*EntryRecord, error)
	GetAll() ([]EntryRecord, error)
}
