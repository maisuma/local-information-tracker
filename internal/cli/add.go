package cli

import (
	//watcher.FileExist(filepath)の使用 watcherパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	//index.AddTrack(filepath)の使用 indexパッケージ
	"fmt"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Add(filepath string) {
	if !watcher.FileExist(filepath) {
		fmt.Println("File not exist")
	} else {
		track_id, err := index.AddTrack(filepath) //トラックIDの発行と取得
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
