package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Log(targetFilepath string) {
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Printf("Error resolving database path:%v\n", err)
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
	commits, err := idx.GetCommitsList(track_id) //trackIDを使って変更履歴のリストを取得
	if err != nil {
		fmt.Printf("Error occurred in getting commits list: %v\n", err)
		return
	}
 	fmt.Println("Commit History:")
    fmt.Println("==============================")
	for i, commitID := range commits {
		fmt.Printf("%d. Commit ID: %d\n", i+1, commitID)
	}
    fmt.Println("==============================")
	return
}
