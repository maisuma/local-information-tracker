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
	} else {
		dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
		if err != nil {
			fmt.Println("Error resolving database path:", err)
			return
		}
		idx, err := index.NewDBIndexer(dbPath) //構造体を生成
		if err != nil {
			fmt.Println("Error occured in creating indexer:%w", err)
			return
		}
		track_id, err := idx.AddTrack(targetFilepath) //トラックIDの発行と取得
		if err != nil {
			fmt.Println("Error occured in adding track")
			return
		}
		fmt.Println("Adding file is complete")
		fmt.Printf("Track ID:%d\n", track_id)
		return
	}
}
