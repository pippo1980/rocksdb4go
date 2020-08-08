package api

import (
	"bytes"
	"errors"
	"sirius/rocksdb/internel/cgo"
)

func (storage *CGOStorage) PrefixGet(key []byte, limit int) ([][]byte, [][]byte, error) {
	return storage.PrefixGet0(storage.ReadOpts, key, limit)
}

func (storage *CGOStorage) PrefixGet0(opts *cgo.ReadOptions, prefix []byte, limit int) ([][]byte, [][]byte, error) {
	if storage.Status != StorageStatusOpened.Code() {
		return nil, nil, errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.ReadOpts
	}

	iterator := storage.DB.NewIteratorCF(opts, storage.DataCF)
	defer iterator.Close()

	keys := make([][]byte, 0)
	values := make([][]byte, 0)

	iterator.Seek(prefix)
	for i := 0; i < limit && iterator.Valid(); i++ {
		key := cgo.RetrieveSliceData(iterator.Key())

		if !bytes.HasPrefix(key, prefix) {
			break
		}

		keys = append(keys, key)
		values = append(values, cgo.RetrieveSliceData(iterator.Value()))
		iterator.Next()
	}

	return keys, values, nil
}

func (storage *CGOStorage) PrefixDelete(key []byte, limit int) ([][]byte, error) {
	return storage.PrefixDelete0(storage.WriteOpts, key, limit)
}

func (storage *CGOStorage) PrefixDelete0(opts *cgo.WriteOptions, prefix []byte, limit int) ([][]byte, error) {
	if storage.Status != StorageStatusOpened.Code() {
		return nil, errors.New("storage not opened")
	}

	if opts == nil {
		opts = storage.WriteOpts
	}

	wb := cgo.NewWriteBatch()
	defer wb.Destroy()

	iterator := storage.DB.NewIteratorCF(storage.ReadOpts, storage.DataCF)
	defer iterator.Close()

	keys := make([][]byte, 0)
	iterator.Seek(prefix)

	for i := 0; i < limit && iterator.Valid(); i++ {
		key := cgo.RetrieveSliceData(iterator.Key())

		if !bytes.HasPrefix(key, prefix) {
			break
		}

		keys = append(keys, key)
		wb.DeleteCF(storage.DataCF, key)

		iterator.Next()
	}

	if err := storage.DB.Write(opts, wb); err != nil {
		return nil, err
	} else {
		return keys, nil
	}

}
