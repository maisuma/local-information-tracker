# lit - Local Information Tracker

## 概要

`lit`（Local Information Tracker）は、ローカルファイルの変更を追跡し、スナップショットを作成するツールです。ファイルが変更されるたびに自動的にバージョンを保存し、必要に応じて過去の状態に復元できます。

## 主な機能

- **ファイル追跡**: 指定したファイルの変更を自動的に監視
- **スナップショット作成**: ファイルが変更されると自動的にバージョンを保存
- **履歴管理**: ファイルの変更履歴を確認
- **復元機能**: 過去の任意のバージョンにファイルを復元
- **効率的なストレージ**: チャンク化により重複データを削減

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/maisuma/local-information-tracker
cd local-information-tracker

# ビルド
go build -o lit ./cmd/lit

# （オプション）パスに追加
sudo mv lit /usr/local/bin/
```

## 使用方法

### 基本コマンド

```bash
# ヘルプを表示
lit help

# バージョン情報を表示
lit version
```

### ファイルの追跡

```bash
# ファイルを追跡対象に追加
lit add <filepath>

# 例：
lit add myfile.txt
```

### 追跡中のファイルを確認

```bash
# 追跡中のファイル一覧を表示
lit list
```

### 変更履歴の確認

```bash
# ファイルの変更履歴を表示
lit log <filepath>

# 例：
lit log myfile.txt
```

### ファイルの復元

```bash
# ファイルを過去の状態に復元
lit restore <commitID>

# 例：
lit restore 42
```

### ファイルの追跡解除

```bash
# ファイルを追跡対象から除外
lit remove <filepath>

# 例：
lit remove myfile.txt
```

### ガベージコレクション

```bash
# 不要なデータをクリーンアップ
lit gc
```

## Watcher（自動監視）

`lit-watcher`を使用すると、追跡中のファイルを自動的に監視し、変更があった場合に自動的にスナップショットを作成します。

```bash
# Watcherをビルド
go build -o lit-watcher ./cmd/lit-watcher

# Watcherを起動
./lit-watcher
```

## 仕組み

1. **チャンク化**: ファイルを小さなチャンクに分割し、各チャンクのハッシュを計算
2. **重複排除**: 既存のチャンクと同じハッシュがあれば再利用
3. **スナップショット**: チャンクのハッシュリストをコミットとして保存
4. **復元**: 保存されたハッシュリストからチャンクを読み出してファイルを再構築

## 技術スタック

- **言語**: Go 1.25.4
- **ファイル監視**: fsnotify
- **ストレージ**: カスタムパックファイルフォーマット
- **インデックス**: SQLiteベースのインデックス（予定）

## 開発状況

このプロジェクトは現在開発中です。以下の機能が実装済みまたは実装予定です：

- [x] ファイル監視機能（Watcher）
- [x] チャンク化とハッシュ計算
- [x] スナップショット作成
- [x] 復元機能
- [ ] CLI完全実装
- [ ] インデックスの永続化
- [ ] ガベージコレクション
- [ ] テストスイート

## ライセンス

このプロジェクトはオープンソースです。

## 貢献

プルリクエストやイシューの報告を歓迎します！

## リンク

- リポジトリ: https://github.com/maisuma/local-information-tracker
