package api

import (
	"errors"
	"sirius/rocksdb/internel/cgo"
)

func (storage *CGOStorage) Get(key []byte) ([]byte, error) {
	return storage.Get0(storage.ReadOpts, key)
}

func (storage *CGOStorage) Get0(opts *cgo.ReadOptions, key []byte) ([]byte, error) {
	if storage.Status != StorageStatusOpened.Code() {
		return nil, errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.ReadOpts
	}

	if result, err := storage.DB.GetCF(opts, storage.DataCF, key); err != nil {
		return nil, err
	} else {
		return cgo.RetrieveSliceData(result), nil
	}
}

func (storage *CGOStorage) GetBatch(keys ...[]byte) ([][]byte, error) {
	return storage.GetBatch0(storage.ReadOpts, keys...)
}

func (storage *CGOStorage) GetBatch0(opts *cgo.ReadOptions, keys ...[]byte) ([][]byte, error) {
	if storage.Status != StorageStatusOpened.Code() {
		return nil, errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.ReadOpts
	}

	if results, err := storage.DB.MultiGetCF(opts, storage.DataCF, keys...); err != nil {
		return nil, err
	} else {
		values := make([][]byte, len(keys))

		for index, result := range results {
			if result.Size() > 0 {

				values[index] = cgo.RetrieveSliceData(result)
			}
		}

		return values, nil
	}
}
