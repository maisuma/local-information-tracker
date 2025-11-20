package cli

import(
	//index.untrack(track_id)の使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
	"strings"
    "fmt"
    "os"
    "log"
)

//①track IDを使ってindex.goの関数を呼び出してcommit IDのリスト？を取得
//②全部表示

// Logは、指定されたファイルの変更履歴を表示します。
func Log(track_id int){ 
	log := 
	//listの中身を表示
	for _, entry := range log {
		fmt.Println(entry)
	}
	return
}