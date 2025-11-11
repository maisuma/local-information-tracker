package storage

type PackFiler interface {
	Write(path string, data []byte) (offset int64, size int64, err error)
	Read(path string, offset int64, size int64) ([]byte, error)
}