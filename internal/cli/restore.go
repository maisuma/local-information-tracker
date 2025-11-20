package cli

import (
	"fmt"
	//snapshot.Restore(track_id, commitID)の使用 snapshotパッケージ
	"github.com/maisuma/local-information-tracker/internal/core/snapshot"
)

func Restore(commitID int) {
	err := new(snapshot.Snapshotter).Restore(commitID)
	if err != nil {
		fmt.Println("Error occurred in restoring")
		return
	}
	fmt.Println("Restoring is complete")
	return
}
