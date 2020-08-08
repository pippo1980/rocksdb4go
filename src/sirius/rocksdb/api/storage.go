package api

type Storage interface {
	Open(preOpen func(Storage) error) error
	Close(postClose func(Storage) error)
	IsOpen() bool
	GetStatus() StorageStatus

	PrefixGet(key []byte, limit int) ( /*keys*/ [][]byte /*values*/, [][]byte, error)
	PrefixDelete(key []byte, limit int) ( /*keys*/ [][]byte, error)
	Get(key []byte) ( /*value*/ []byte, error)
	GetBatch(keys ...[]byte) ( /*values*/ [][]byte, error)
	Put(key, value []byte) error
	PutBatch(keys, values [][]byte) error
	Delete(keys ...[]byte) error
	Flush(block bool)
}
