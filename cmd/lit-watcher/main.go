package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	"github.com/maisuma/local-information-tracker/internal/engine/chunker"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"github.com/maisuma/local-information-tracker/internal/engine/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // ゴールーチンの
	defer cancel()

	// 停止シグナルの処理
	sigCh := make(chan os.Signal, 1)
	//killによる強制終了やCtrl+Cによる割り込みをキャッチ
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)
	go func() {
		<-sigCh
		cancel()
	}()

	// 初期化処理
	basePath := "./testdata"
	dbPath := "./testdata/index.db"

	// 必要なディレクトリを作成
	if err := os.MkdirAll(basePath, 0755); err != nil {
		log.Fatalf("Failed to create base path: %v", err)
	}

	// Indexerを初期化
	idx, err := index.NewDBIndexer(dbPath)
	fmt.Println("1")
	if err != nil {
		log.Fatalf("Failed to initialize indexer: %v", err)
	}
	defer idx.Close()

	// Storageを初期化
	stor, err := storage.New(basePath)
	fmt.Println("2")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Chunkerを初期化
	ch := chunker.NewChunker(idx, stor, 8192, 4096, 16384)
	fmt.Println("3")

	// Snapshotterを初期化
	snap := snapshot.NewSnapshotter(ch, stor, idx)
	fmt.Println("4")

	defaultDuration := 2 * time.Second

	watcherImpl, err := watcher.NewWatcher(snap, idx, defaultDuration)
	if err != nil {
		log.Fatalf("watcher init failed: %v", err)
	}
	var watch watcher.WatcherAPI = watcherImpl

	go func() {
		if err := watch.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("watcher stopped: %v", err)
			cancel()
		}
	}()

}
