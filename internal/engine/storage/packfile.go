package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type PackFiler interface {
	Write(data []byte) (int, int64, int64, error)
	Read(packID int, offset int64, size int64) ([]byte, error)
}
const (
    storageDirName = "storage"
    maxPackSize    = 8 * 1024 * 1024 * 1024 //閾値:8GB
)

type Storage struct {
	basePath string
	mutex sync.Mutex
	currentPackID int //packIDの探索時間を短縮
}

func New(basePath string) (*Storage, error) {
	storagePath := filepath.Join(basePath, storageDirName)
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil,fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &Storage{
		basePath: storagePath,
		currentPackID: 1,
	}, nil
}


func (s *Storage) Write(data []byte) (int, int64, int64, error) {
	//複数のWrite()が同時に動作しないようにする
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(data) == 0 {
		return 0, 0, 0, fmt.Errorf("empty data")
	}

	dataSize := int64(len(data))
	//最後に書き込んだpakcIDから探索
	for {
		packFilePath := filepath.Join(s.basePath, fmt.Sprintf("%d.pack", s.currentPackID))
		fileInfo, statErr := os.Stat(packFilePath)

		var file *os.File
		var err error
		var offset int64 = 0

		if os.IsNotExist(statErr) {
			file, err = os.OpenFile(packFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err != nil {
				if os.IsExist(err) {
					continue // 競合したら再試行
				}
				return 0, 0, 0, fmt.Errorf("failed to create packfile: %w", err)
			}
			//新規ファイルのoffsetは0
			offset = 0

		} else if statErr != nil {
			return 0, 0, 0, fmt.Errorf("failed to stat packfile: %w", statErr)

		} else {
			if fileInfo.Size()+dataSize > maxPackSize {
				//容量オーバーなら次のIDへ進んでループ継続
				s.currentPackID++
				continue
			}

			file, err = os.OpenFile(packFilePath, os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return 0, 0, 0, fmt.Errorf("failed to open packfile: %w", err)
			}

			offset, err = file.Seek(0, io.SeekEnd)
			if err != nil {
				file.Close()
				return 0, 0, 0, fmt.Errorf("failed to seek: %w", err)
			}
		}

		n, err := file.Write(data)
		closeErr := file.Close()

		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to write data: %w", err)
		}
		if closeErr != nil {
			return 0, 0, 0, fmt.Errorf("failed to close file: %w", closeErr)
		}
		if n != len(data) {
			return 0, 0, 0, fmt.Errorf("incomplete write")
		}

		// 成功したら現在のPackIDとオフセットを返す
		return s.currentPackID, offset, int64(n), nil
	}
}

func (s *Storage) Read(packID int, offset int64, size int64) ([]byte, error) {
	if offset < 0{
		return nil, fmt.Errorf("invalid size: %d", offset)
	}
	if size <= 0 {
		return nil, fmt.Errorf("invalid size: %d", size)
	}
	packFilePath := filepath.Join(s.basePath, fmt.Sprintf("%d.pack", packID))
	//ファイルを開く
	file, err := os.Open(packFilePath)

	if err != nil{
		return nil, fmt.Errorf("failed to open packfile %s: %w", packFilePath, err)
	}
	defer file.Close()

	//指定されたoffsetへ移動
	_, err = file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	//データの読み取り
	buffer := make([]byte, size)
	_, err = io.ReadFull(file, buffer)
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, fmt.Errorf("failed to read data: %w", err)
		}

		return buffer, nil
}