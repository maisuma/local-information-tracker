// package lit 本番用
package main //デバッグ用

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	//add等の関数の使用
	"github.com/maisuma/local-information-tracker/internal/cli"
	//GetTrackIDByFileの使用 indexパッケージ
	"github.com/maisuma/local-information-tracker/internal/engine/index"
)

// FilenameToAbsFilepathは、ファイル名を絶対ファイルパスに変換します。
//
//	カレントディレクトリとファイル名から絶対パスを生成します。
func FilenameToAbsFilepath(filename string) (abspath string, err error) {
	fmt.Println("Filename関数まできた") //デバッグ用
	// filepath.Abs は、カレントディレクトリを基準に絶対パスを返します。
	// これにより、ユーザーが "file.txt" や "./file.txt" のように入力しても、
	// DBには一貫した絶対パスが保存されます
	//nilはnull
	abspath, err = filepath.Abs(filename)
	if err != nil {
		log.Printf("エラー: ファイルパス '%s' の絶対パスを取得できませんでした: %v", filename, err)
		return "", err
	}
	fmt.Println("Filename関数を抜ける") //デバッグ用
	return abspath, nil
}

//trackIDはファイルごとに一意に発行されるID
//commitIDは変更ごとに一意に発行されるID

// コマンド一覧
//lit add <filename> 3
//lit remove <filename> 3
//lit gc 2
//lit log <filename> 3
//lit restore <commitID> 3
//lit list 2

//呼び出し一覧
//Add(filepath string)
//Remove(track_id int)
//Gc()
//Log(track_id int)
//Restore(commitID int)
//List()

func main() {

	for true { //デバッグ用無限ループ

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("コマンドを入力してください: ")

		// 標準入力から一行を読み込む
		inputLine, err := reader.ReadString('\n')
		if err != nil && err != os.ErrClosed { // EOF（Ctrl+D/Z）以外のエラーを処理
			fmt.Println("入力読み込みエラー:", err)
			return
		}

		// 読み込んだ文字列から改行や前後の空白を除去
		line := strings.TrimSpace(inputLine)
		if line == "" {
			fmt.Println("入力がありませんでした。")
			return
		}

		// strings.Fieldsを使って、すべての単語をスライスに分割する
		// これにより、単語の数に関係なくすべてが取得される
		parts := strings.Fields(line)

		//ここからは本実装
		//デーモンでスペースごとに区切られている想定
		if parts[0] != "lit" {
			fmt.Println("コマンドは 'lit' から始めてください。")
			return
		} else {
			switch parts[1] {
			case "add":
				if len(parts) != 3 {
					fmt.Println("Usage: lit add <filename>")
					return
				} else {
					abspath, err := FilenameToAbsFilepath(parts[2])
					if err != nil {
						fmt.Println("Error in getting absolute filepath")
						return
					}
					fmt.Println("絶対パスは:", abspath) //デバッグ用
					fmt.Println("今からAdd()を呼びだす")   //デバッグ用
					cli.Add(abspath)
				}
			case "remove":
				if len(parts) != 3 {
					fmt.Println("Usage: lit remove <filename>")
					return
				} else {
					abspath, err := FilenameToAbsFilepath(parts[2])
					if err != nil {
						fmt.Println("Error in getting absolute filepath")
						return
					}
					track_id, err := new(index.DBIndexer).GetTrackIDByFile(abspath)
					if err != nil {
						fmt.Println("Error in getting track ID by filepath")
						return
					}
					cli.Remove(track_id)
				}
			case "gc": //実装不明
				fmt.Println("here is gc") //デバッグ用
				if len(parts) != 2 {
					fmt.Println("Usage: lit gc")
					return
				} else {
					cli.Gc()
				}
			case "log":
				if len(parts) != 3 {
					fmt.Println("Usage: lit log <filename>")
					return
				} else {
					abspath, err := FilenameToAbsFilepath(parts[2])
					if err != nil {
						fmt.Println("Error in getting absolute filepath")
						return
					}
					track_id, err := new(index.DBIndexer).GetTrackIDByFile(abspath)
					if err != nil {
						fmt.Println("Error in getting track ID by filepath")
						return
					}
					cli.Log(track_id)
				}
			case "restore":
				if len(parts) != 3 {
					fmt.Println("Usage: lit restore <commitID>")
					return
				} else {
					var commitID int
					_, err := fmt.Sscanf(parts[2], "%d", &commitID)
					if err != nil {
						fmt.Println("Error in parsing commit ID")
						return
					}
					cli.Restore(commitID)
				}
			case "list":
				if len(parts) != 2 {
					fmt.Println("Usage: lit list")
					return
				} else {
					cli.List()
				}
			default:
				fmt.Println("Unknown command:")
				fmt.Println("Pleadease use one of the following commands")
				fmt.Println("lit add <filename>")
				fmt.Println("lit remove <filename>")
				fmt.Println("lit gc")
				fmt.Println("lit log <filename>")
				fmt.Println("lit restore <commitID>")
				fmt.Println("lit list")
				return
			}
		}
	} //デバッグ用無限ループ終わり
}
