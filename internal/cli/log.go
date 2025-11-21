package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Log(targetFilepath string) {
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Println("Error resolving database path:%w", err)
		return
	}
	idx, err := index.NewDBIndexer(dbPath) //構造体を生成
	if err != nil {
		fmt.Println("Error occured in creating indexer:%w", err)
		return
	}
	track_id, err := idx.GetTrackIDByFile(targetFilepath) //filepathからtrackIDを取得
	if err != nil {
		fmt.Println("Error in getting track ID by filepath")
		return
	}
	commits, err := idx.GetCommitsList(track_id) //trackIDを使って変更履歴のリストを取得
	if err != nil {
		fmt.Println("Error occured in getting commits list")
		return
	} else {
		for _, commitID := range commits {
			fmt.Printf("Commit ID:%d\n", commitID)
		}
		return
	}
}
