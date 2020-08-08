package api

import (
	"log"
	"runtime"
	"sirius/rocksdb/internel/cgo"
	"sync/atomic"
	"time"
)

func NewCGOStorage(name, db_path string) (Storage, error) {
	storage := new(CGOStorage)
	storage.Name = name
	storage.Status = StorageStatusClosed.Code()
	storage.BaseDir = db_path

	preOpen := func(s Storage) error {
		storage.ConfigBlock(64*1024, 1000)
		storage.ConfigMemTable(true, 32, 8*1024*1024)
		storage.ConfigDBSize(5, 8*1024*1024, 4)
		// storage.ConfigLog(cgo.InfoInfoLogLevel, 24, 60*60)
		storage.ConfigLog(cgo.DebugInfoLogLevel, 24, 60*60)
		return nil
	}

	if err := storage.Open(preOpen); err != nil {
		return nil, err
	} else {
		return storage, nil
	}

}

type CGOStorage struct {
	Name      string
	Status    int32 // -1:closed, 0:opening, 1:opened
	BaseDir   string
	ReadOpts  *cgo.ReadOptions
	WriteOpts *cgo.WriteOptions
	StoreOpts *cgo.Options
	DB        *cgo.DB
	DefaultCF *cgo.ColumnFamilyHandle
	DataCF    *cgo.ColumnFamilyHandle
}

func (storage *CGOStorage) GetStatus() StorageStatus {
	return StorageStatus(storage.Status)
}

func (storage *CGOStorage) IsOpen() bool {
	return storage.Status == 1
}

func (storage *CGOStorage) ConfigBlock(block_size, cached_block_count int) {
	opts := cgo.NewDefaultBlockBasedTableOptions()
	opts.SetBlockSize(block_size)
	opts.SetBlockCache(cgo.NewLRUCache(block_size * cached_block_count))
	opts.SetIndexType(cgo.KBinarySearchIndexType)
	opts.SetFilterPolicy(cgo.NewBloomFilter(10))

	storage.StoreOpts.SetBlockBasedTableFactory(opts)
}

// write buffer: if each_table=8MB buffer_number=32 then total=1GB
func (storage *CGOStorage) ConfigMemTable(enable bool, table_count, table_bytes_size int) {
	storage.StoreOpts.SetAllowConcurrentMemtableWrites(enable)
	storage.StoreOpts.SetAllowMmapReads(enable)
	storage.StoreOpts.SetAllowMmapWrites(enable)

	if enable {
		storage.StoreOpts.SetMaxWriteBufferNumber(table_count)
		storage.StoreOpts.SetWriteBufferSize(table_bytes_size)
	}
}

// level file and size config
// max_level=5 max_level0_files=32MB each_level_files=16 each_storage_files=128
// file_size_l = target_file_size_base * (multiplier ^ (L-1))
// level_size_l = target_file_size_base * 16 * (multiplier ^ (L-1))
// level-1 : file=8MB    total=128MB
// level-2 : file=32MB   total=512MB
// level-3 : file=128MB  total=2GB
// level-4 : file=1GB    total=16GB
// level-5 : file=4GB    total=64GB
func (storage *CGOStorage) ConfigDBSize(max_level int, file_bytes_base uint64, multiplier int) {
	storage.StoreOpts.SetNumLevels(max_level)
	storage.StoreOpts.SetLevel0FileNumCompactionTrigger(4)
	storage.StoreOpts.SetLevel0SlowdownWritesTrigger(16)
	storage.StoreOpts.SetLevel0StopWritesTrigger(32)
	storage.StoreOpts.SetTargetFileSizeBase(file_bytes_base)
	storage.StoreOpts.SetTargetFileSizeMultiplier(multiplier)
	storage.StoreOpts.SetMaxBytesForLevelBase(file_bytes_base * 16)
	storage.StoreOpts.SetMaxBytesForLevelMultiplier(float64(multiplier))
	storage.StoreOpts.SetMaxOpenFiles(multiplier ^ (max_level + 1))
}

func (storage *CGOStorage) ConfigLog(level cgo.InfoLogLevel, file_count, file_roll_seconds int) {
	storage.StoreOpts.SetInfoLogLevel(level)
	storage.StoreOpts.SetKeepLogFileNum(file_count)
	storage.StoreOpts.SetLogFileTimeToRoll(file_roll_seconds)
}

func (storage *CGOStorage) Flush(block bool) {
	opts := cgo.NewDefaultFlushOptions()
	opts.SetWait(block)

	defer opts.Destroy()

	// for i := 0; i < 3; i++ {
	if err := storage.DB.Flush(opts); err != nil {
		log.Printf("flush storage due to error:%v", err)
		time.Sleep(10 * time.Second)
	} else {
		// break
	}
	// }
}

func (storage *CGOStorage) Open(preOpen func(storage Storage) error) error {
	// 防止重复初始化
	if !atomic.CompareAndSwapInt32(&storage.Status, StorageStatusClosed.Code(), StorageStatusOpening.Code()) {
		return nil
	}

	// default read opts
	storage.ReadOpts = cgo.NewDefaultReadOptions()

	// default write opts
	storage.WriteOpts = cgo.NewDefaultWriteOptions()
	storage.WriteOpts.DisableWAL(false)
	storage.WriteOpts.SetSync(false)

	// init storage opts
	storage.StoreOpts = cgo.NewDefaultOptions()

	// storage runtime env
	env := cgo.NewDefaultEnv()
	env.SetBackgroundThreads(runtime.NumCPU())
	env.SetHighPriorityBackgroundThreads(runtime.NumCPU())
	storage.StoreOpts.SetEnv(env)

	storage.StoreOpts.IncreaseParallelism(runtime.NumCPU())
	storage.StoreOpts.SetCompactionStyle(cgo.LevelCompactionStyle)
	storage.StoreOpts.SetCreateIfMissing(true)
	storage.StoreOpts.SetCreateIfMissingColumnFamilies(true)
	storage.StoreOpts.SetDbLogDir(storage.BaseDir + "/logs")
	storage.StoreOpts.SetWalDir(storage.BaseDir + "/wal")

	//opts.OptimizeForPointLookup(256)
	storage.StoreOpts.SetMaxBackgroundCompactions(2)
	storage.StoreOpts.SetMaxBackgroundFlushes(8)

	column_families := []string{"default", storage.Name}

	column_families_option := make([]*cgo.Options, len(column_families))
	for index := range column_families_option {
		column_families_option[index] = storage.StoreOpts
	}

	// trigger other config
	if preOpen != nil {
		if err := preOpen(storage); err != nil {
			atomic.StoreInt32(&storage.Status, StorageStatusClosed.Code())
			return err
		}
	}

	// open store
	if db, columnFamilyHandlers, err := cgo.OpenDbColumnFamilies(storage.StoreOpts,
		storage.BaseDir,
		column_families,
		column_families_option); err != nil {

		atomic.StoreInt32(&storage.Status, StorageStatusClosed.Code())
		return err
	} else {
		storage.DB = db
		storage.DefaultCF = columnFamilyHandlers[0]
		storage.DataCF = columnFamilyHandlers[1]
	}

	log.Printf("open storage:[%v] with store dir:[%v]", storage.Name, storage.BaseDir)

	// change open flag
	atomic.StoreInt32(&storage.Status, StorageStatusOpened.Code())
	return nil

}

func (storage *CGOStorage) Close(postClose func(Storage) error) {
	// 防止重复关闭
	if !atomic.CompareAndSwapInt32(&storage.Status, StorageStatusOpened.Code(), StorageStatusClosed.Code()) {
		return
	}

	// trigger other close
	if postClose != nil {
		if err := postClose(storage); err != nil {
			log.Printf("post close storage due to error:%v", err)
		}
	}

	storage.Flush(true)

	// close rocks store
	storage.ReadOpts.Destroy()
	storage.WriteOpts.Destroy()
	storage.DataCF.Destroy()
	storage.DefaultCF.Destroy()

	storage.DB.Close()

	log.Printf("close storage:[%v] with store dir:[%v]", storage.Name, storage.BaseDir)
}
