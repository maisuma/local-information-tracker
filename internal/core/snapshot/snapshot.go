// snapshot.go

package snapshot

import (
	"internal/core/chunker"
	"internal/core/index"
	"internal/core/storage"
)

// 依存するチームAのコンポーネントを保持する
type Snapshotter struct {
	chunker *chunker.Chunker
	storage *storage.Storage
	index   *index.Index
}

// コンストラクタ (依存性の注入)
func NewSnapshotter(c *chunker.Chunker, s *storage.Storage, i *index.Index) *Snapshotter {
	return &Snapshotter{chunker: c, storage: s, index: i}
}

/*
 * スナップショットを作成する
 * watcherから呼ばれるメイン関数。
 * 内部で差分比較や削除ハンドリングを行う。
 */
func (s *Snapshotter) Create(trackID int, filepath string) (commitID int64, err error) {
	// 1. (追加) ファイルが存在するかチェック
	// 2. (追加) 直前のコミットのハッシュリストをDB (index) から取得
	// 3. チームA (chunker) を呼び出し、現在のファイルのハッシュリストを取得
	// 4. (追加) 直前のハッシュリストと比較
	// 5. 差分があれば、新しいチャンクをチームA (storage) に保存依頼
	// 6. チームA (index) にコミット情報を書き込み依頼
}

/*
 * ファイルを指定したコミットIDの状態に復元する
 * lit restore から呼ばれる。
 */
func (s *Snapshotter) Restore(trackID int, commitID int64) error {
	// 1. チームA (index) から、commitIDが持つハッシュのリストと順序を取得
	// 2. 各ハッシュについて、チームA (index) から保存場所 (packID, offset, length) を取得
	// 3. チームA (storage) から実データを順に読み出す
	// 4. ファイルに書き込んで復元する
}
