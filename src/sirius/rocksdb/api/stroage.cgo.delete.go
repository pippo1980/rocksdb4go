package api

import (
	"errors"
	"sirius/rocksdb/internel/cgo"
)

func (storage *CGOStorage) Delete(keys ...[]byte) error {
	return storage.Delete0(storage.WriteOpts, keys...)
}

func (storage *CGOStorage) Delete0(opts *cgo.WriteOptions, keys ...[]byte) error {
	if storage.Status != StorageStatusOpened.Code() {
		return errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.WriteOpts
	}

	// 更新数据
	wb := cgo.NewWriteBatch()
	defer wb.Destroy()

	for _, key := range keys {
		wb.DeleteCF(storage.DataCF, key)
	}

	if err := storage.DB.Write(opts, wb); err != nil {
		return err
	} else {
		return nil
	}
}
