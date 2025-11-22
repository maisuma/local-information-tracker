package index

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Indexer interface {
	AddTrack(filepath string) (int, error)
	RemoveTrack(trackID int) error
	GetTrackIDByFile(filepath string) (int, error)
	//GetTrackIDByCommit(commitID int) (int, error)
	GetFilepath(trackID int) (string, error)
	SaveHash(hash []byte, packID int, offset int64, size int64) error
	GetPack(hash []byte) (int, int64, int64, error)
	LookupHash(hash []byte) (bool, error)
	AddCommit(trackID int, hashes [][]byte) (int, error)
	GetHashes(commitID int) ([][]byte, error)
	GetTracksList() ([]int, error)
	GetCommitsList(trackID int) ([]int, error)
	GetCommit(commitID int) (Commit, error)
}
//indexerインターフェースの実装構造体
type DBIndexer struct {
	db *sql.DB
}

type Commit struct {
	TrackID int
	CommitID int
	Created_at time.Time
}

//データベース接続を初期化
func NewDBIndexer(dbPath string) (*DBIndexer, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	//外部キーの制約を有効化
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	indexer := &DBIndexer{db: db}
	if err := indexer.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}
	return indexer, nil
}

//データベースの接続を閉じる
func (i *DBIndexer) Close() error {
	return i.db.Close()
}

//テーブルを作成
func (i *DBIndexer) initSchema() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS tracked_files (
			track_id INTEGER PRIMARY KEY AUTOINCREMENT,
			filepath TEXT NOT NULL UNIQUE
		);`,
		`CREATE TABLE IF NOT EXISTS commits (
			commit_id INTEGER PRIMARY KEY AUTOINCREMENT,
			track_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (track_id) REFERENCES tracked_files(track_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS hash_storage (
			hash TEXT PRIMARY KEY,
			pack_id INTEGER,
			offset INTEGER,
			size INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS commit_contents (
			commit_id INTEGER,
			hash_order INTEGER,
			hash TEXT,
			PRIMARY KEY (commit_id, hash_order),
			FOREIGN KEY (commit_id) REFERENCES commits(commit_id) ON DELETE CASCADE,
			FOREIGN KEY (hash) REFERENCES hash_storage(hash)
		);`,
	}

	for _, query := range queries {
		if _, err := i.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
//新しいファイルをトラック対象に追加する
func (i *DBIndexer) AddTrack(filepath string) (int, error) {
	query := `INSERT INTO tracked_files (filepath) VALUES (?)`
	res, err := i.db.Exec(query, filepath)
	if err != nil {
		id, getErr := i.GetTrackIDByFile(filepath)
		if getErr == nil {
			return id, nil
		}
		return 0, fmt.Errorf("failed to add track: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return int(id), nil
}

//トラック対象から削除する
func (i *DBIndexer) RemoveTrack(trackID int) error {
	res, err := i.db.Exec(`DELETE FROM tracked_files WHERE track_id = ?`, trackID)
	if err != nil {
		return fmt.Errorf("failed to remove track: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("track_id %d not found", trackID)
	}
	return nil
}

//ファイルパスからtrackIDを取得する。
func (i *DBIndexer) GetTrackIDByFile(filepath string) (int, error) {
	var id int
	err := i.db.QueryRow(`SELECT track_id FROM tracked_files WHERE filepath = ?`, filepath).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("filepath '%s' not found", filepath)
	}
	return id, err
}

//commitIDからtrackIDを取得する。
func (i *DBIndexer) GetTrackIDByCommit(commitID int) (int, error) {
	var trackID int
	err := i.db.QueryRow(`SELECT track_id FROM commits WHERE commit_id = ?`, commitID).Scan(&trackID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("commit_id %d not found", commitID)
	}
	return trackID, err
}

//trackIDからファイルパスを取得する。
func (i *DBIndexer) GetFilepath(trackID int) (string, error) {
	var path string
	err := i.db.QueryRow(`SELECT filepath FROM tracked_files WHERE track_id = ?`, trackID).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("track_id %d not found", trackID)
	}
	return path, err
}

//ハッシュとpackfileの情報の対応関係を保存する。
func (i *DBIndexer) SaveHash(hash []byte, packID int, offset int64, size int64) error {
	hashStr := hex.EncodeToString(hash)
	//既に存在する場合は情報を更新する
	query := `INSERT OR REPLACE INTO hash_storage (hash, pack_id, offset, size) VALUES (?, ?, ?, ?)`
	_, err := i.db.Exec(query, hashStr, packID, offset, size)
	if err != nil {
		return fmt.Errorf("failed to save hash: %w", err)
	}
	return nil
}

//ハッシュからpackfileの情報を取得する。
func (i *DBIndexer) GetPack(hash []byte) (int, int64, int64, error) {
	hashStr := hex.EncodeToString(hash)
	var packID int
	var offset, size int64
	
	query := `SELECT pack_id, offset, size FROM hash_storage WHERE hash = ?`
	err := i.db.QueryRow(query, hashStr).Scan(&packID, &offset, &size)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, 0, 0, fmt.Errorf("hash not found")
	}
	if err != nil {
		return 0, 0, 0, err
	}
	return packID, offset, size, nil
}

//ハッシュが追加済みか確認する。
func (i *DBIndexer) LookupHash(hash []byte) (bool, error) {
	hashStr := hex.EncodeToString(hash)
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM hash_storage WHERE hash = ?)`
	err := i.db.QueryRow(query, hashStr).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

//コミット（スナップショット）を行う。そのコミットのハッシュをDBに保存する。
func (i *DBIndexer) AddCommit(trackID int, hashes [][]byte) (int, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return 0, err
	}
	//関数終了時にCommitされていなければロールバック
	defer tx.Rollback()
	
	query := `INSERT INTO commits (track_id, created_at) VALUES (?, ?)`

	res, err := tx.Exec(query, trackID, time.Now()) 
	if err != nil {
		return 0, fmt.Errorf("failed to insert commit record: %w", err)
	}
	commitID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	//commit_contentsテーブルへハッシュリストを順序付きで挿入
	stmt, err := tx.Prepare(`INSERT INTO commit_contents (commit_id, hash_order, hash) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for order, h := range hashes {
		hashStr := hex.EncodeToString(h)
		
		//外部キー制約により、hash_storageにこのハッシュが存在しないとエラーになる。
		_, err = stmt.Exec(commitID, order, hashStr)
		if err != nil {
			return 0, fmt.Errorf("failed to link hash %s to commit: %w", hashStr, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(commitID), nil
}

//コミットのハッシュを取得する。
func (i *DBIndexer) GetHashes(commitID int) ([][]byte, error) {
	query := `SELECT hash FROM commit_contents WHERE commit_id = ? ORDER BY hash_order ASC`
	rows, err := i.db.Query(query, commitID)
	if err != nil {
		return nil, fmt.Errorf("failed to query commit hashes: %w", err)
	}
	defer rows.Close()

	var hashes [][]byte
	for rows.Next() {
		var hashStr string
		if err := rows.Scan(&hashStr); err != nil {
			return nil, err
		}
		
		decoded, err := hex.DecodeString(hashStr)
		if err != nil {
			return nil, fmt.Errorf("database contains invalid hex hash: %w", err)
		}
		hashes = append(hashes, decoded)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hashes, nil
}

//全てのtrackIDのリストを返す
func (i *DBIndexer) GetTracksList() ([]int, error) {
	query := `SELECT track_id FROM tracked_files`
	rows, err := i.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tracked files: %w", err)
	}
	defer rows.Close()

	var trackIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan track ID: %w", err)
		}
		trackIDs = append(trackIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracked files: %w", err)
	}
	return trackIDs, nil
}

//trackIDのコミット履歴を返す
func (i *DBIndexer) GetCommitsList(trackID int) ([]int, error) {
	query := `SELECT commit_id FROM "commits" WHERE track_id = ? ORDER BY commit_id DESC`
	rows, err := i.db.Query(query, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to query commits: %w", err)
	}
	defer rows.Close()

	var commitIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan commit ID: %w", err)
		}
		commitIDs = append(commitIDs, id)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating commits: %w", err)
	}
	return commitIDs, nil
}

func (i *DBIndexer) GetCommit(commitID int) (Commit, error) {
	var c Commit
	var created_atStr string
	query := `SELECT commit_id, track_id, created_at FROM "commits" WHERE commit_id = ?`
	err := i.db.QueryRow(query, commitID).Scan(&c.CommitID, &c.TrackID, &created_atStr)
	if errors.Is(err, sql.ErrNoRows) {
		return Commit{}, fmt.Errorf("commitID %d not found: %w", commitID, err)
	}
	if err != nil {
		return Commit{}, fmt.Errorf("failed to query commit: %w", err)
	}
	if t, perr := time.Parse(time.RFC3339, created_atStr); perr == nil {
		c.Created_at =t
	}else if t, perr := time.Parse("2006-01-02_15:04:05", created_atStr); perr == nil {
		c.Created_at = t
	}else {
		c.Created_at = time.Time{}
	}
	return c, nil
}