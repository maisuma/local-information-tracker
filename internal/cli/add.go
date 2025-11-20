package cli

import (
	//FileExist(filepath)の使用 watcherパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	//AddTrack(filepath)の使用 indexパッケージ
	"fmt"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Add(filepath string) {
	if !watcher.FileExist(filepath) {
		fmt.Println("File not exist")
	} else {

		idx, err := index.NewDBIndexer(filepath)             //構造体を生成
		track_id, err := idx.AddTrack(filepath)              //トラックIDの発行と取得
		fmt.Println("AddTrack returned track_id:", track_id) //デバッグ用
		if err != nil {
			fmt.Println("Error occured in adding track")
			return
		} else {
			err = new(watcher.Watcher).AddWatch(filepath)
			if err != nil {
				fmt.Println("Error occured in adding watch")
				return
			}
		}
		fmt.Println("Adding file is complete")
		fmt.Printf("Track ID:%d\n", track_id)
		return
	}
}
