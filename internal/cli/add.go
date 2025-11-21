package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Add(targetFilepath string) {
	if !watcher.FileExist(targetFilepath) {
		fmt.Println("File not exist")
		return
	}
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Printf("Error resolving database path:", err)
		return
	}
	idx, err := index.NewDBIndexer(dbPath) //構造体を生成
	if err != nil {
		fmt.Printf("Error occurred in creating indexer: %v\n", err)
		return
	}
	defer idx.Close()
	track_id, err := idx.AddTrack(targetFilepath) //トラックIDの発行と取得
	if err != nil {
		fmt.Printf("Error occurred in adding track: %v\n", err)
		return
	}
	fmt.Println("Adding file is complete")
	fmt.Printf("Track ID:%d\n", track_id)
	return
}
