package cgo

//
//import (
//	"log"
//	"runtime"
//	"sirius/rocksdb4go/api"
//	"sync/atomic"
//	"time"
//)
//
//type RocksStorage struct {
//	Status    int32 // -1:closed, 0:opening, 1:opened
//	BaseDir   string
//	ReadOpts  *ReadOptions
//	WriteOpts *WriteOptions
//	StoreOpts *Options
//	Store     *DB
//	DefaultCF *ColumnFamilyHandle
//	DataCF    *ColumnFamilyHandle
//	IndexCFS  map[string]*ColumnFamilyHandle
//}
//
//func (storage *RocksStorage) GetStatus() int32 {
//	return storage.Status
//}
//
//func (storage *RocksStorage) IsOpen() bool {
//	return storage.Status == 1
//}
//
//func (storage *RocksStorage) Open(preOpen func(storage api.Storage) error) error {
//	// 防止重复初始化
//	if !atomic.CompareAndSwapInt32(&storage.Status, -1, 0) {
//		return nil
//	}
//
//	// default read opts
//	storage.ReadOpts = NewDefaultReadOptions()
//
//	// default write opts
//	storage.WriteOpts = NewDefaultWriteOptions()
//
//	// init storage opts
//	storage.StoreOpts = NewDefaultOptions()
//
//	// storage runtime env
//	env := NewDefaultEnv()
//	env.SetBackgroundThreads(runtime.NumCPU())
//	env.SetHighPriorityBackgroundThreads(runtime.NumCPU())
//	storage.StoreOpts.SetEnv(env)
//
//	storage.StoreOpts.IncreaseParallelism(runtime.NumCPU())
//	storage.StoreOpts.SetAllowConcurrentMemtableWrites(true)
//	storage.StoreOpts.SetAllowMmapReads(true)
//	storage.StoreOpts.SetAllowMmapWrites(true)
//	// opts.SetBlockBasedTableFactory(storage.BlockOpts())
//	storage.StoreOpts.SetCompactionStyle(LevelCompactionStyle)
//	storage.StoreOpts.SetCreateIfMissing(true)
//	storage.StoreOpts.SetCreateIfMissingColumnFamilies(true)
//	storage.StoreOpts.SetDbLogDir(storage.BaseDir + "/logs")
//	storage.StoreOpts.SetWalDir(storage.BaseDir + "/wal")
//
//	//opts.OptimizeForPointLookup(256)
//	storage.StoreOpts.SetMaxBackgroundCompactions(2)
//	storage.StoreOpts.SetMaxBackgroundFlushes(8)
//
//	// write buffer each_buffer=8MB buffer_number=32 total=1GB
//	storage.StoreOpts.SetMaxWriteBufferNumber(32)
//	storage.StoreOpts.SetWriteBufferSize(1024 * 1024 * 32)
//
//	// level file and size config
//	// max_level=5 max_level0_files=32MB each_level_files=16 each_storage_files=128
//	// file_size_l = target_file_size_base * (target_file_size_multiplier ^ (L-1))
//	// level_size_l = max_bytes_for_level_base * (max_bytes_for_level_multiplier ^ (L-1))
//	// level-1 : file=8MB    total=128MB
//	// level-2 : file=32MB   total=512MB
//	// level-3 : file=128MB  total=2GB
//	// level-4 : file=1GB    total=16GB
//	// level-5 : file=4GB    total=64GB
//
//	max_level := 5
//	file_base := uint64(1024 * 1024 * 8)
//	multiplier := 4
//	level_base := file_base * 16
//
//	storage.StoreOpts.SetNumLevels(max_level)
//	storage.StoreOpts.SetLevel0FileNumCompactionTrigger(4)
//	storage.StoreOpts.SetLevel0SlowdownWritesTrigger(16)
//	storage.StoreOpts.SetLevel0StopWritesTrigger(32)
//	storage.StoreOpts.SetTargetFileSizeBase(file_base)
//	storage.StoreOpts.SetTargetFileSizeMultiplier(multiplier)
//	storage.StoreOpts.SetMaxBytesForLevelBase(level_base)
//	storage.StoreOpts.SetMaxBytesForLevelMultiplier(float64(multiplier))
//
//	// log config
//	storage.StoreOpts.SetInfoLogLevel(DebugInfoLogLevel)
//	storage.StoreOpts.SetLogFileTimeToRoll(60 * 60)
//	storage.StoreOpts.SetKeepLogFileNum(24 * 30)
//
//	// trigger other config
//	if err := preOpen(storage); err != nil {
//		atomic.StoreInt32(&storage.Status, -1)
//		return err
//	} else {
//		// change open flag
//		atomic.StoreInt32(&storage.Status, 1)
//		return nil
//	}
//}
//
//func (storage *RocksStorage) Close(postClose func(api.Storage) error) {
//	// 防止重复关闭
//	if !atomic.CompareAndSwapInt32(&storage.Status, 1, -1) {
//		return
//	}
//
//	// trigger other close
//	if err := postClose(storage); err != nil {
//		log.Printf("post close storage due to error:%v", err)
//	}
//
//	// close rocks store
//	storage.ReadOpts.Destroy()
//	storage.WriteOpts.Destroy()
//
//	for i := 0; i < 3; i++ {
//		flushOpts := NewDefaultFlushOptions()
//		flushOpts.SetWait(true)
//
//		if err := storage.Store.Flush(flushOpts); err != nil {
//			log.Printf("flush storage due to error:%v", err)
//			time.Sleep(10 * time.Second)
//		} else {
//			break
//		}
//	}
//
//	storage.Store.Close()
//}
