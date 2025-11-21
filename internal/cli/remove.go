package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Remove(targetFilepath string) {
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Printf("Error resolving database path: %v\n", err)
		return
	}
	idx, err := index.NewDBIndexer(dbPath) //構造体を生成
	if err != nil {
		fmt.Printf("Error occurred in creating indexer: %v\n", err)
		return
	}
	defer idx.Close()
	track_id, err := idx.GetTrackIDByFile(targetFilepath) //filepathからtrackIDを取得
	if err != nil {
		fmt.Printf("Error in getting track ID by filepath: %v\n", err)
		return
	}
	err = idx.RemoveTrack(track_id) //追跡対象からの除外または停止
	if err != nil {
		fmt.Printf("Error occurred in removing track: %v\n", err)
		return
	}
	fmt.Println("Removing track is complete")
	return
}
