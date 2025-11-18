package cli

import(
	//index.untrack(track_id)の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"strings"
    "fmt"
    "os"
    "log"
)

func Log(track_id int) {
	log := index.GetList(track_id)
	//listの中身を表示
	for _, entry := range log {
		fmt.Println(entry)
	}
	return
}