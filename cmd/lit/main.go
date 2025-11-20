// package lit 本番用
package main //デバッグ用

import (
	//add等の関数の使用
	"strconv"

	"github.com/maisuma/local-information-tracker/internal/cli"
	//GetTrackIDByFileの使用 indexパッケージ
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

// FilenameToAbsFilepathは、ファイル名を絶対ファイルパスに変換します。
//
//	カレントディレクトリとファイル名から絶対パスを生成します。
func FilenameToAbsFilepath(filename string) (filepath string, err error) {
	// filepath.Abs は、カレントディレクトリを基準に絶対パスを返します。
	// これにより、ユーザーが "file.txt" や "./file.txt" のように入力しても、
	// DBには一貫した絶対パスが保存されます。
	//nilはnull
	filepath, err := filepath.Abs(filename)
	if err != nil {
		log.Printf("エラー: ファイルパス '%s' の絶対パスを取得できませんでした: %v", filename, err)
		return "", err
	}
	return filepath, nil
}

//trackIDはファイルごとに一意に発行されるID
//commitIDは変更ごとに一意に発行されるID

// コマンド一覧
//lit add <filename>
//lit remove <filename>
//lit gc
//lit log <filename>
//lit restore <commitID>
//lit list

//呼び出し一覧
//Add(filepath string)
//Remove(track_id int)
//Gc()
//Log(track_id int)
//Restore(commitID int)
//List()

func main() {

	for true { //デバッグ用無限ループ

		scanner := bufio.NewScanner(os.Stdin)

		// Fieldsは、スペース、タブ、改行で区切られたトークンを返す
		words := strings.Fields(scanner)

		//コマンドライン引数を取得
		// os.Args は、コマンドライン引数全体を要素ごとに分割して保持しています。
		//args := os.Args　本番用

		//正しいコマンド入力かの確認
		//
		if args[0] != "lit" {
			//エラー文
			fmt.Println("Command not exist")
			fmt.Println("Please inclode 'lit' ")
			return
		} else {
			var filename string
			var filepath string
			var track_id int
			var commitID int

			switch args[1] {
			case "add":
				if !args[3] { //コマンドライン引数3つのみ
					filename = args[2]
					//filepathに変換
					filepath, err = FilenameToAbsFilepath(filename)
					cli.Add(filepath)
				} else {
					fmt.Println("many Command to use add")
					return
				}
			case "remove":
				if !args[3] { //コマンドライン引数3つのみ
					filename = args[2]
					//filepathに変換
					filepath, err = FilenameToAbsFilepath(filename)
					//track_idに変換
					track_id = index.GetTrackIDByFile(filepath)
					cli.Remove(track_id)
				} else {
					fmt.Println("many Command to use remove")
					return
				}
			case "gc": //詳細不明
			case "log":
				if !args[3] { //コマンドライン引数3つのみ
					filename = args[2]
					//filepathに変換
					filepath, err = FilenameToAbsFilepath(filename)
					//track_idに変換
					track_id = index.GetTrackIDByFile(filepath)
					cli.Log(track_id)
				} else {
					fmt.Println("many Command to use log")
					return
				}
			case "restore":
				if !args[4] { //コマンドライン引数4つのみ
					filename = args[2]
					//filepathに変換
					filepath, err = FilenameToAbsFilepath(filename)
					//track_idに変換
					track_id = index.GetTrackIDByFile(filepath)
					commitID, err = strconv.Atoi(args[3])
					if err != nil {
						fmt.Println("Invalid commit ID")
						return
					}
					cli.Restore(track_id, commitID)
				} else {
					fmt.Println("many Command to use restore")
					return
				}
			case "list": //-a でallいらなくない？
				if !args[2] { //コマンドライン引数3つのみ
					cli.List()
				} else {
					fmt.Println("many Command to use list")
					return
				}
			default:
				fmt.Println("Please include Command")
			}
		}

	} //デバッグ用無限ループ終わり
}
