// watcher.go
package watcher

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

type WatcherAPI interface {
	Start(ctx context.Context) error
	AddWatch(filepath string) error
}

// ファイルの有無の確認
func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

// type Snapshotter interface { //エラーありかも
// 	Snapshot(trackID int)
// }

// Watcher が使うDB関数のインターフェース
// (index.Index がこのインターフェースを満たします)
// type Indexer interface {
// 	GetTrackID(filepath string) (int, error)
// 	// (AddWatch/RemoveWatchのために) GetAllTrackedFiles() (map[int]string, error)
// }

type Watcher struct {
	fsWatcher   *fsnotify.Watcher     // fsnotifyのWatcher
	snapshotter *snapshot.Snapshotter // スナップショットを作成するためのインターフェース
	index       index.Indexer         // trackIDを取得するためのインターフェース
	// デバウンス用のタイマーを管理するマップ
	// key: filepath, value: timer
	debounceTimers map[string]*time.Timer
	// マップを安全に操作するためのミューテックス
	mu sync.Mutex

	// 待機する時間
	debounceDuration time.Duration
}

func NewWatcher(snap *snapshot.Snapshotter, idx index.Indexer, debounceDuration time.Duration) (*Watcher, error) { //コンストラクタ
	watcher, err := fsnotify.NewWatcher() // fsnotifyのWatcherを作成
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsWatcher:        watcher,
		snapshotter:      snap,
		index:            idx,
		debounceTimers:   make(map[string]*time.Timer), // キーはファイルパス、値はタイマー
		debounceDuration: debounceDuration,
	}, nil
}

func (w *Watcher) AddWatch(filepath string) error { // ファイルを監視対象に追加 a
	return w.fsWatcher.Add(filepath)
}

// ファイルシステムの監視のループ
func (w *Watcher) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("watcherの停止")
			return ctx.Err()
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return nil
			}
			//bitmaskで書き込みを判定
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.triggerSnapshot(event.Name)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				// ファイルが削除された場合の処理

			}
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return nil
			}
			return err
		case msg := <-w.index.NotifyChan():
			fmt.Printf("Received notification: %s\n", msg)
			w.AddWatch(msg) // 監視対象に追加
		}

	}
}

// デバウンス処理をトリガーする関数

func (w *Watcher) triggerSnapshot(filepath string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if pretimer, ok := w.debounceTimers[filepath]; ok {
		pretimer.Stop() // 既存のタイマーを停止
	}
	//afterFuncで指定時間後に関数を実行

	timer := time.AfterFunc(w.debounceDuration, func() {
		w.finishDebounce(filepath)
	})

	// 新しいタイマーをマップに保存
	w.debounceTimers[filepath] = timer // filepathが未知の場合追加も行われる。

}

func (w *Watcher) finishDebounce(filepath string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.debounceTimers, filepath)

	w.executeSnapshot(filepath)
}

func (w *Watcher) executeSnapshot(filepath string) error {

	err := w.snapshotter.Snapshot(filepath)
	if err != nil {
		// エラーハンドリング
		return err
	}
	return nil
}
