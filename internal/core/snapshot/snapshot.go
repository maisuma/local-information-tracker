// snapshot.go

package snapshot

import (
	"fmt"
	"os"

	"github.com/maisuma/local-information-tracker/internal/engine/chunker"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"github.com/maisuma/local-information-tracker/internal/engine/storage"
)

// 依存するチームAのコンポーネントを保持する
type Snapshotter struct {
	chunker *chunker.Chunker
	storage *storage.Storage
	index   index.Indexer
}

type SnapshotterAPI interface {
	Snapshot(filepath string) error
	Restore(trackID int, commitID int) error
}

// コンストラクタ (依存性の注入)
func NewSnapshotter(c *chunker.Chunker, s *storage.Storage, i index.Indexer) *Snapshotter {
	return &Snapshotter{chunker: c, storage: s, index: i}
}

// スナップショットを作成するメソッド
func (s *Snapshotter) Snapshot(filepath string) error {
	// チャンク化してハッシュを取得
	hashes, err := s.chunker.ChunkAndSave(filepath)
	if err != nil {
		return err
	}
	// trackIDを取得してコミットを追加
	trackID, err := s.index.GetTrackIDByFile(filepath)
	if err != nil {
		return err
	}
	commitID, err := s.index.AddCommit(trackID, hashes)
	if err != nil {
		return err
	}
	fmt.Printf("Snapshot created for file: %s\n", filepath)
	fmt.Printf("trackID: %d, chunk hashes: %v\n", trackID, hashes) //デバッグ用出力
	fmt.Printf("New commitID: %d\n", commitID)
	return nil
}

func (s *Snapshotter) Restore(commitID int) error { //復元用
	trackID, err := s.index.GetTrackIDByCommit(commitID)
	if err != nil {
		return err
	}
	filepath, err := s.index.GetFilepath(trackID)
	if err != nil {
		return err
	}
	Hashes, err := s.index.GetHashes(commitID)
	if err != nil {
		return err
	}

	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0644)
	//O_CREATE追加の可能性あり 追跡対象が消えたら新規作成するかも

	if err != nil {
		return err
	}
	defer out.Close()

	for _, hash := range Hashes {
		packID, offset, size, err := s.index.GetPack(hash) //GetPackがまだない
		if err != nil {
			return err
		}
		data, err := s.storage.Read(packID, offset, size)
		if err != nil {
			return err
		}
		if _, err := out.Write(data); err != nil {
			return err
		}
	}
	return nil
}
