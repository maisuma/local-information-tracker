package strage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type PackLocation struct {
	packID int
	offset int64
	size   int64
}

type PackStorage struct {
	baseDir       string
	maxFileSize   int64
	mtx           sync.Mutex
	currentPackID int
	currentFile   *os.File
}

// PackIdから.packファイルを作成する
func (s *PackStorage) getPackPath(packID int) string {
	return filepath.Join(s.baseDir, fmt.Sprintf("%d.pack", packID))
}

func (s *PackStorage) openCurrentFileForWrite() error {
	path := s.getPackPath(s.currentPackID)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return err
	}
	s.currentFile = f
	return nil
}

func New(baseDir string, maxFileSize int64) (*PackStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	currentID := 0
	files, err := filepath.Glob(filepath.Join(baseDir, "*.pack"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		baseName := filepath.Base(file)
		idStr := strings.TrimSuffix(baseName, ".pack")
		id, err := strconv.Atoi(idStr)
		if err == nil {
			if id > currentID {
				currentID = id
			}
		}
	}
	if currentID == 0 {
		currentID = 1
	}
	storage := &PackStorage{
		baseDir:       baseDir,
		maxFileSize:   maxFileSize,
		currentPackID: currentID,
	}
	if err := storage.openCurrentFileForWrite(); err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *PackStorage) Save(data []byte) (*PackLocation, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	dataSize := int64(len(data))
	if dataSize == 0 {
		return nil, nil
	}
	offset, err := s.currentFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	if (offset+dataSize) > s.maxFileSize && offset != 0 {
		//ファイルを閉じる
		if err := s.currentFile.Close(); err != nil {
			return nil, err
		}
		//次のファイル
		s.currentPackID++

		//新しいファイルを開く
		if err := s.openCurrentFileForWrite(); err != nil {
			return nil, err
		}
		offset = 0
	}
	written, err := s.currentFile.Write(data)
	if int64(written) != dataSize {
		return nil, err
	}
	return &PackLocation{
		packID: s.currentPackID,
		offset: offset,
		size:   dataSize,
	}, nil
}

func (s *PackStorage) Load(packID int, offset int64, size int64) ([]byte, error) {
	if packID <= 0 {
		return nil, fmt.Errorf("invalid packID")
	}
	if offset < 0 {
		return nil, fmt.Errorf("invalid offset")
	}
	if size < 0 {
		return nil, fmt.Errorf("invalid size")
	}

	path := s.getPackPath(packID)

	f, err := os.Open(path)

	if err != nil {
		s.mtx.Lock()
		isCurrent := (packID == s.currentPackID)
		s.mtx.Unlock()

		if isCurrent {
			f, err = os.Open(path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	//読み込み終了
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if offset+size > fileInfo.Size() {
		return nil, fmt.Errorf("over file size")
	}

	buffer := make([]byte, size)
	n, err := f.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return nil, err
	}
	//不完全なファイルの検出
	if int64(n) != size {
		return nil, err
	}
	return buffer[:n], nil
}

// 現在開いている書き込み用のファイルを閉じる
func (s *PackStorage) Close() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.currentFile != nil {
		err := s.currentFile.Close()
		s.currentFile = nil
		if err != nil {
			return err
		}
	}
	return nil
}
