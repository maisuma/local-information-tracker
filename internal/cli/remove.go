package cli

import (
	"fmt"
	//RemoveTrack(track_id)の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

func Remove(track_id int) {
	err := new(index.DBIndexer).RemoveTrack(track_id) //追跡対象からの除外または停止
	if err != nil {
		fmt.Println("Error occured in removing track")
		return
	}
	fmt.Println("Removing track is complete")
	return
}
