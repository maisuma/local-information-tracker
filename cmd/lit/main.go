package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/maisuma/local-information-tracker/internal/cli"
)

// filepath.Abs は、カレントディレクトリを基準に絶対パスを返します
func FilenameToAbsFilepath(filename string) (abspath string, err error) {
	abspath, err = filepath.Abs(filename)
	if err != nil {
		log.Printf("エラー: ファイルパス '%s' の絶対パスを取得できませんでした: %v", filename, err)
		return "", err
	}
	return abspath, nil
}

func main() {

	// コマンドライン引数を取得
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Please provide a command.")
		return
	}

	parts := args[1:]

	switch parts[0] {
	case "add":
		if len(parts) != 2 {
			fmt.Println("Usage:add <filepath>")
			return
		}
		abspath, err := FilenameToAbsFilepath(parts[1])
		if err != nil {
			fmt.Println("Error in getting absolute filepath")
			return
		}
		cli.Add(abspath)

	case "remove":
		if len(parts) != 2 {
			fmt.Println("Usage:remove <filepath>")
			return
		}
		abspath, err := FilenameToAbsFilepath(parts[1])
		if err != nil {
			fmt.Println("Error in getting absolute filepath")
			return
		}
		cli.Remove(abspath)

	case "gc":
		if len(parts) != 1 {
			fmt.Println("Usage:gc")
			return
		}
		cli.Gc()

	case "log":
		if len(parts) != 2 {
			fmt.Println("Usage:log <filepath>")
			return
		}
		abspath, err := FilenameToAbsFilepath(parts[1])
		if err != nil {
			fmt.Println("Error in getting absolute filepath")
			return
		}
		cli.Log(abspath)

	case "restore":
		if len(parts) != 2 {
			fmt.Println("Usage:restore <commitID>")
			return
		}
		var commitID int
		_, err := fmt.Sscanf(parts[1], "%d", &commitID)
		if err != nil {
			fmt.Println("Error in parsing commit ID")
			return
		}
		cli.Restore(commitID)

	case "list":
		if len(parts) != 1 {
			fmt.Println("Usage:list")
			return
		}
		cli.List()

	default:
		fmt.Println("Unknown command:")
		fmt.Println("Please use one of the following commands")
		fmt.Println("add <filepath>")
		fmt.Println("remove <filepath>")
		fmt.Println("gc")
		fmt.Println("log <filepath>")
		fmt.Println("restore <commitID>")
		fmt.Println("list")
		return
	}
}
