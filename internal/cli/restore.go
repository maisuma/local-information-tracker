package cli

import(
	//snapshot.restore(track_id, commitID)の使用 snapshotパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
	"strings"
    "fmt"
    "os"
    "log"
)

func Restore(track_id int, commitID int) {
	snapshot.Restore(track_id, commitID)
	fmt.Println("Restoring is complete")
	return
}