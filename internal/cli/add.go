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
	_, err = idx.AddTrack(targetFilepath) //トラックIDの発行と取得
	if err != nil {
		fmt.Printf("Error occurred in adding track: %v\n", err)
		return
	}
	fmt.Println("File successfully added to tracking!")
	fmt.Println("==============================")
	fmt.Println("The file at the following filepath has been added to tracking")
	fmt.Printf("Filepath: %s\n", targetFilepath)
	fmt.Println("==============================")
	return
}
