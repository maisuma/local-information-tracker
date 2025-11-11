// watcher.go
package watcher

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ファイルの有無の確認
func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)

	return err == nil
}

type Snapshotter interface { //エラーありかも
	Snapshot(trackID int)
}

// Watcher が使うDB関数のインターフェース
// (index.Index がこのインターフェースを満たします)
type Index interface {
	GetTrackID(filepath string) (int, error)
	// (AddWatch/RemoveWatchのために) GetAllTrackedFiles() (map[int]string, error)
}

type Watcher struct {
	fsWatcher   *fsnotify.Watcher // fsnotifyのWatcher
	snapshotter Snapshotter       // スナップショットを作成するためのインターフェース
	index       Index             // trackIDを取得するためのインターフェース

	// デバウンス用のタイマーを管理するマップ
	// key: filepath, value: timer
	debounceTimers map[string]*time.Timer
	// マップを安全に操作するためのミューテックス
	mu sync.Mutex

	// 待機する時間
	debounceDuration time.Duration
}

func NewWatcher(snap Snapshotter, idx Index, debounceDuration time.Duration) (*Watcher, error) { //コンストラクタ
	watcher, err := fsnotify.NewWatcher() // fsnotifyのWatcherを作成
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsWatcher:        watcher,
		snapshotter:      snap,
		index:            idx,
		debounceTimers:   make(map[string]*time.Timer),
		debounceDuration: debounceDuration,
	}, nil
}

// ファイルシステムの監視
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
			//bitmaskで書き込みと削除を判定
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

		}

	}
}

// デバウンス処理をトリガーする関数

func (w *Watcher) triggerSnapshot(filepath string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if timer, ok := w.debounceTimers[filepath]; ok {
		timer.Stop()
	}
	w.debounceTimers[filepath] = time.AfterFunc(w.debounceDuration, func() {
		w.cleanupTimer(filepath)
		w.executeSnapshot(filepath)
	})
}

func (w *Watcher) cleanupTimer(filepath string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.debounceTimers, filepath)
}
func (w *Watcher) executeSnapshot(filepath string) {
	trackID, err := w.index.GetTrackID(filepath)
	if err != nil {
		// エラーハンドリング
		return
	}
	w.snapshotter.Snapshot(trackID)
	if err != nil {
		// エラーハンドリング
		return
	}
}
