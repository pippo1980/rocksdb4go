package api

import (
	"errors"
	"sirius/rocksdb/internel/cgo"
)

func (storage *CGOStorage) Put(key, value []byte) error {
	return storage.Put0(storage.WriteOpts, key, value)
}

func (storage *CGOStorage) Put0(opts *cgo.WriteOptions, key, value []byte) error {
	if storage.Status != StorageStatusOpened.Code() {
		return errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.WriteOpts
	}

	if err := storage.DB.PutCF(opts, storage.DataCF, key, value); err != nil {
		return err
	} else {
		return nil
	}
}

func (storage *CGOStorage) PutBatch(keys, values [][]byte) error {
	return storage.PutBatch(keys, values)
}

func (storage *CGOStorage) PutBatch0(opts *cgo.WriteOptions, keys, values [][]byte) error {

	if storage.Status != StorageStatusOpened.Code() {
		return errors.New("storage not opened")
	}

	if len(keys) != len(values) {
		return errors.New("invalid key/val pairs")
	}

	if opts == nil {
		opts = storage.WriteOpts
	}

	// 更新数据
	wb := cgo.NewWriteBatch()
	defer wb.Destroy()

	for index, key := range keys {
		wb.PutCF(storage.DataCF, key, values[index])
	}

	if err := storage.DB.Write(opts, wb); err != nil {
		return err
	} else {
		return nil
	}

}
