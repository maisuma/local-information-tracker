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
}

func New(basePath string) (*Storage, error) {
	storagePath := filepath.Join(basePath, storageDirName)
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil,fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &Storage{basePath: storagePath}, nil
}


func (s *Storage) Write(data []byte) (int, int64, int64, error) {
	//複数のWrite()が同時に動作しないようにする
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(data) == 0 {
		return 0, 0, 0, fmt.Errorf("empty data")
	}

	dataSize := int64(len(data))
	//.packファイルを探す
	var packFilePath string
	var file *os.File
    var err error

    // .packファイルを探すまたは作成する
    for i := 1; ; i++ {
        packFilePath = filepath.Join(s.basePath, fmt.Sprintf("%d.pack", i))
        fileInfo, statErr := os.Stat(packFilePath)

        if os.IsNotExist(statErr) {
            // .packファイルが存在しない場合は新規作成
            file, err = os.OpenFile(packFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
            if err != nil {
                if os.IsExist(err) {
                    // 競合が発生したら再試行
                    continue
                }
                return 0, 0, 0, fmt.Errorf("failed to create packfile: %w", err)
            }
            defer file.Close()
            break
        }

        if statErr != nil {
            return 0, 0, 0, fmt.Errorf("failed to stat packfile: %w", statErr)
        }

        // ファイルサイズを確認し、データを追加できるかチェック
        if fileInfo.Size()+dataSize <= maxPackSize {
            file, err = os.OpenFile(packFilePath, os.O_WRONLY|os.O_APPEND, 0644)
            if err != nil {
                return 0, 0, 0, fmt.Errorf("failed to open packfile: %w", err)
            }
            defer file.Close()
            break
        }
    }
	//offsetを求める
	offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to seek to offset %d in packfile %s: %w", offset, packFilePath, err)
	}

	//データの書き込み
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to write data: %w", err)
	}
	
	//データが完全に書きこめたか確認
	if writtenBytes != len(data) {
		return 0, 0, 0, fmt.Errorf("incomplete write")
	}
	var packID int
	_, err = fmt.Sscanf(filepath.Base(packFilePath), "%d.pack", &packID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse packID: %w", err)
	}
	return packID, offset, int64(writtenBytes), nil
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
	readData, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d bytes from packfile %s at offset %d: %w", size, packFilePath, offset, err)
	}
	if int64(readData) != size {
		return nil, io.ErrUnexpectedEOF
	}
	return buffer[:readData], nil
}