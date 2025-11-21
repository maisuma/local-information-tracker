package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func List() {
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
	list, err := idx.GetTracksList() //追跡対象になっているファイルのtrackID一覧を取得
	if err != nil {
		fmt.Println("Error occured in listing tracks")
		return
	} else {
		for _, track_id := range list {
			filepath, err := idx.GetFilepath(track_id) //trackIDからファイルパスを取得
			if err != nil {
				fmt.Println("Error in getting filepath by track ID")
				return
			}
			fmt.Printf("Filepath:%s\n", filepath)
		}
		return
	}
}
