package cli

import (
	"fmt"
	//GetTracksList()の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

// index.何かの使用 indexパッケージ
// GetFilepath(track_id)で一意なファイルパスが取得できる。
// ①index.goから追跡対象の ファイルのtrackIDをまとめて入手
// ②一つずつtrack IDからファイルパスに変換して標準出力
func List() {
	list, err := new(index.DBIndexer).GetTracksList() //追跡対象のファイルのtrackID一覧を取得
	if err != nil {
		fmt.Println("Error occured in listing tracks")
		return
	} else {
		for _, track_id := range list {
			filepath, err := new(index.DBIndexer).GetFilepath(track_id) //trackIDからファイルパスを取得
			if err != nil {
				fmt.Println("Error occured in getting filepath")
				return
			} else {
				fmt.Printf("Track ID:%d, Filepath:%s\n", track_id, filepath)
			}
		}
		return
	}
}
