package chunker

import (
	"crypto/sha256"
	"io"
	"math/big"
	"os"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"github.com/maisuma/local-information-tracker/internal/engine/storage"
)

type ChunkerAPI interface {
	ChunkAndSave(filepath string) ([][]byte, error)
}

const (
	k = 48
	b = 257 // 基数
	m = (1 << 61) - 1 // 法
	mask = (1 << 13) - 1 // 8KB
)

// b^(k-1) % m (区間の一番古いやつ)
var p uint64

func init() {
	p = new(big.Int).Exp(big.NewInt(b), big.NewInt(k-1), big.NewInt(m)).Uint64()
}

type Chunker struct {
	idx index.Indexer 
	pak storage.PackFiler 

	avgChunkSize int
	minChunkSize int
	maxChunkSize int
}


func NewChunker(
	idx index.Indexer,
	pak storage.PackFiler,
	avgSize, minSize, maxSize int,
) *Chunker {
	return &Chunker{idx, pak, avgSize, minSize, maxSize}
}

func (c *Chunker) ChunkAndSave(filepath string) ([][]byte, error) {
	file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
	defer file.Close()

    hashes, err := c.findCutPoints(file)
    if err != nil {
        return nil, err
    }
	
	return hashes, nil
}

//chunkDataを保存する
func (c *Chunker) saveChunk(chunkData []byte) (hash []byte, err error) {
	h := sha256.Sum256(chunkData)
	hash = h[:]

	exists, err := c.idx.LookupHash(hash)
	if err != nil {
		return nil, err
	}
	if exists {
		return hash, nil
	}

	packID, offset, size, err := c.pak.Write(chunkData)
	if err != nil {
		return nil, err
	}

	err = c.idx.SaveHash(hash, packID, offset, size)
	if err != nil {
		return nil, err
	}

	return hash, nil
}


func (c *Chunker) findCutPoints(reader io.Reader) ([][]byte, error) {
	var res [][]byte

	//チャンクのバッファ
	chunkBuf := make([]byte, 0, c.maxChunkSize)

	//次に追い出す奴
	var window [k]byte
	pos := 0

	//ローリングハッシュ
	var hash uint64 = 0

	//読み込みバッファ
	readBuf := make([]byte, 4096*8)

	for {
		n, err := reader.Read(readBuf)
		if n > 0 {
			for i := 0; i < n; i++ {
				new := readBuf[i]

				chunkBuf = append(chunkBuf, new)

				old := window[pos]

				window[pos] = new

				pos = (pos + 1) % k

				hash = (hash - (uint64(old) * p) % m + m) % m
				hash = (hash * b) % m
				hash = (hash + uint64(new)) % m

				//今のチャンク幅
				size := len(chunkBuf)
				
				isHit := (hash & mask) == 1
				isMinOK := size >= c.minChunkSize
				isMaxReached := size >= c.maxChunkSize

				if (isHit && isMinOK) || isMaxReached {
					sha256Hash, err := c.saveChunk(chunkBuf)
					if err != nil {
						return nil, err
					}
					res = append(res, sha256Hash)

					chunkBuf = chunkBuf[:0]
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	if len(chunkBuf) > 0 {
		sha256Hash, err := c.saveChunk(chunkBuf)
		if err != nil {
			return nil, err
		}
		res = append(res, sha256Hash)
	}

	return res, nil
}