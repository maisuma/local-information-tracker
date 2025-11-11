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

//.packを保存するファイルを管理する構造体
var (
	globalMutex sync.Mutex
	basePath string
)

func New(basePath string) error {
	storagePath := filepath.Join(basePath, storageDirName)
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return fmt.Errorf("faild to create storage directory: %w",err)
	}
	return nil
}

func findPackFile(dataSize int64) (string, error) {
    // 新しい .pack ファイルを作成または既存ファイルを選択
    var packFilePath string
    for i := 1; ; i++ {
        packFilePath = filepath.Join(basePath, fmt.Sprintf("%d.pack", i))
        fileInfo, err := os.Stat(packFilePath)

        if os.IsNotExist(err) {
            // .pack ファイルが存在しない場合は新規作成
            file, err := os.OpenFile(packFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
            if err != nil {
                if os.IsExist(err) {
					//競合が発生したら、再度ファイルを確認する
					continue
				}
				return "", fmt.Errorf("faild to create packfile: %w",err)
            }
            file.Close()
            return packFilePath, nil
        }
		if err != nil {
			return "", fmt.Errorf("faild to stat packfile: %w",err)
		}
        // ファイルサイズを確認し、データを追加できるかチェック
        if fileInfo.Size()+dataSize <= maxPackSize {
            return packFilePath, nil
        }
    }
}

func Write(data []byte) (int, int64, int64, error) {
	//複数のWrite()が同時に動作しないようにする
	globalMutex.Lock()
	defer globalMutex.Unlock()
	if len(data) == 0 {
		return 0, 0, 0, fmt.Errorf("emptydata")
	}

	dataSize := int64(len(data))
	//.packファイルを探す
	packFilePath, err := findPackFile(dataSize)
	if err != nil{
		return  0, 0, 0, err
	}

	file, err := os.OpenFile(packFilePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return 0, 0, 0,err
	}
	
	defer file.Close()

	//offsetを求める
	offset , err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("faild to seek: %w",err)
	}

	//データの書き込み
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("faild to write data: %w",err)
	}
	
	//データが完全に書きこめたか確認
	if writtenBytes != len(data) {
		return 0, 0, 0, fmt.Errorf("incompe write")
	}
	var packID int
	_, err = fmt.Sscanf(filepath.Base(packFilePath), "%d.pack", &packID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("faild to parse packID:%w",err)
	}
	return packID, offset, int64(writtenBytes), nil
}

func Read(packID int, offset int64, size int64) ([]byte, error) {
	if size <= 0 {
		return nil, fmt.Errorf("invalid size:%d",size)
	}
	packFilePath := filepath.Join(basePath, fmt.Sprintf("%d.pack", packID))
	//ファイルを開く
	file, err := os.Open(packFilePath)

	if err != nil{
		return nil, err
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
		return nil, err
	}
	return buffer[:readData], nil
}