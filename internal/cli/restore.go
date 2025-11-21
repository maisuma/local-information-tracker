package cli

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
	"github.com/maisuma/local-information-tracker/internal/engine/chunker"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"github.com/maisuma/local-information-tracker/internal/engine/storage"
)

func Restore(commitID int) {
	// データベースファイルへの絶対パスの取得
	dbPath, err := filepath.Abs("./lit.db")
	if err != nil {
		log.Fatalf("Failed to resolve database path: %v", err)
	}

	// ベースパスを取得（ストレージ用）
	basePath, err := filepath.Abs(".")
	if err != nil {
		log.Fatalf("Failed to resolve base path: %v", err)
	}

	// Indexerを初期化
	idx, err := index.NewDBIndexer(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize indexer: %v", err)
	}
	defer idx.Close()

	// Storageを初期化
	stor, err := storage.New(basePath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Chunkerを初期化
	ch := chunker.NewChunker(idx, stor, 8192, 4096, 16384)

	// Snapshotterを初期化
	snap := snapshot.NewSnapshotter(ch, stor, idx)

	// 指定された commitID を使って復元
	err = snap.Restore(commitID)
	if err != nil {
		log.Fatalf("Failed to restore commit: %v", err)
	}

	fmt.Println("Restore completed successfully")
}
