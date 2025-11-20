package cli

import (
	//GetCommitsList(trackID)の使用 indexパッケージ
	"fmt"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

//①trackIDを使ってindex.goの関数を呼び出してcommit IDのリスト？を取得
//②全部表示

// Logは、指定されたファイルの変更履歴を表示します。
func Log(track_id int) {
	commit_list, err := new(index.DBIndexer).GetCommitsList(track_id) //trackIDからコミットID一覧を取得
	if err != nil {
		fmt.Println("Error occured in getting commit list")
		return
	} else {
		fmt.Printf("Change history for Track ID:%d\n", track_id)
		for _, commit_id := range commit_list {
			fmt.Printf("Commit ID:%d\n", commit_id)
		}
		return
	}
}
