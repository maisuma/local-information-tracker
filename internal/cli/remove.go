package cli

import(
	//watcher.FileExist(filepath)の使用 watcherパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	//index.untrack(track_id)の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"strings"
    "fmt"
    "os"
    "log"
)

func Remove(track_id int) {
	index.Untrack(track_id)//追跡対象からの除外または停止
	fmt.Println("Removing track is complete")
	return
}