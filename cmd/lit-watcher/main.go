package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
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

	//これら以外の初期化コードも必要
	var (
		idx  index.Indexer //indexの初期化未定
		snap *snapshot.Snapshotter
	)

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
	//追跡対象のファイルを追加するときのコードも必要

}
