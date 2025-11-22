package cli

import (
	"fmt"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func List() {
	dbPath, err := filepath.Abs("./lit.db") //データベースファイルへの絶対パスの取得
	if err != nil {
		fmt.Printf("Error resolving database path: %v\n", err)
		return
	}
	idx, err := index.NewDBIndexer(dbPath) //構造体を生成
	if err != nil {
		fmt.Printf("Error occurred in creating indexer")
		return
	}
	defer idx.Close()
	list, err := idx.GetTracksList() //追跡対象になっているファイルのtrackID一覧を取得
	if err != nil {
		fmt.Printf("Error occurred in listing tracks: %v\n", err)
		return
	}
	fmt.Println("Tracked Files:")
    fmt.Println("==============================")
	for i, track_id := range list {
		filepath, err := idx.GetFilepath(track_id) //trackIDからファイルパスを取得
		if err != nil {
			fmt.Printf("Error in getting filepath by track ID: %v\n", err)
			return
		}
        fmt.Printf("%d. %s\n", i+1, filepath)
	}
    fmt.Println("==============================")
	return
}
