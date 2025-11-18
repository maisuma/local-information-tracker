package cli

import(
	//watcher.FileExist(filepath)の使用 watcherパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/watcher"
	//index.AddTrack(filepath)の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"strings"
    "fmt"
    "os"
    "log"
)

func Add(filepath String) {

	if !watcher.FileExist(filepath){
		fmt.Println("File not exist")
	}else {
		var track_id int
		track_id, err = index.AddTrack(filepath)//トラックIDの発行と取得
		if err != nil {
			fmt.Println("Error occured in adding track")
			return
		}else{
		fmt.Println("Adding file is complete")
		fmt.Printf("Track ID:%d\n", track_id)
		return
		}
	}
}