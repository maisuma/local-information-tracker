package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Remove(targetFilepath string) {
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Println("Error resolving database path:", err)
		return
	}
	idx, err := index.NewDBIndexer(dbPath) //構造体を生成
	if err != nil {
		fmt.Println("Error occured in creating indexer")
		return
	}
	track_id, err := idx.GetTrackIDByFile(targetFilepath) //filepathからtrackIDを取得
	if err != nil {
		fmt.Println("Error in getting track ID by filepath")
		return
	}
	err = idx.RemoveTrack(track_id) //追跡対象からの除外または停止
	if err != nil {
		fmt.Println("Error occured in removing track")
		return
	}
	fmt.Println("Removing track is complete")
	return
}
