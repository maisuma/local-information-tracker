package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "help", "--help", "-h":
		printHelp()
	case "version", "--version", "-v":
		fmt.Printf("lit version %s\n", version)
	case "add":
		fmt.Println("add command: ファイルを追跡対象に追加します")
		fmt.Println("使用方法: lit add <filepath>")
	case "list":
		fmt.Println("list command: 追跡中のファイル一覧を表示します")
		fmt.Println("使用方法: lit list")
	case "log":
		fmt.Println("log command: ファイルの変更履歴を表示します")
		fmt.Println("使用方法: lit log <filepath>")
	case "remove":
		fmt.Println("remove command: ファイルを追跡対象から除外します")
		fmt.Println("使用方法: lit remove <filepath>")
	case "restore":
		fmt.Println("restore command: ファイルを過去の状態に復元します")
		fmt.Println("使用方法: lit restore <commitID>")
	case "gc":
		fmt.Println("gc command: 不要なデータをクリーンアップします")
		fmt.Println("使用方法: lit gc")
	default:
		fmt.Printf("エラー: 不明なコマンド '%s'\n", command)
		fmt.Println("使用可能なコマンドを確認するには 'lit help' を実行してください")
		os.Exit(1)
	}
}

func printHelp() {
	help := `lit - Local Information Tracker

使用方法:
  lit <command> [arguments]

利用可能なコマンド:
  add <filepath>     ファイルを追跡対象に追加します
  list               追跡中のファイル一覧を表示します
  log <filepath>     ファイルの変更履歴を表示します
  remove <filepath>  ファイルを追跡対象から除外します
  restore <commitID> ファイルを過去の状態に復元します
  gc                 不要なデータをクリーンアップします
  help               このヘルプメッセージを表示します
  version            バージョン情報を表示します

説明:
  litはローカルファイルの変更を追跡し、スナップショットを作成するツールです。
  ファイルが変更されるたびに自動的にバージョンを保存し、必要に応じて
  過去の状態に復元できます。

使用例:
  lit add myfile.txt        # myfile.txtを追跡開始
  lit list                  # 追跡中のファイルを表示
  lit log myfile.txt        # myfile.txtの変更履歴を表示
  lit restore 42            # commitID 42の状態に復元
  lit remove myfile.txt     # myfile.txtの追跡を停止
  lit gc                    # 不要なデータを削除

詳細情報:
  https://github.com/maisuma/local-information-tracker
`
	fmt.Print(help)
}
