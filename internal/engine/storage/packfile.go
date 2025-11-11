package storage

type PackFiler interface {
	Write(data []byte) (int, int64, int64, error)
	Read(packID int, offset int64, size int64) ([]byte, error)
}