package api

type StorageStatus int32

const (
	StorageStatusClosed  StorageStatus = -1
	StorageStatusOpening StorageStatus = 0
	StorageStatusOpened  StorageStatus = 1
)

func (status StorageStatus) Code() int32 {
	switch status {
	case StorageStatusClosed:
		return -1
	case StorageStatusOpening:
		return 0
	case StorageStatusOpened:
		return 1
	}

	return -2
}
