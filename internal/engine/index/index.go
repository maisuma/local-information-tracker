package index

type Indexer interface {
    AddTrack(filepath string) (int, error)
    RemoveTrack(trackID int) error
	GetTrackIDByFile(filepath string) (int, error)
	GetTrackIDByCommit(commitID int) (int, error)
	GetFilepath(trackID int) (string, error)
	SaveHash(hash []byte, packID int, offset int64, size int64) error
	GetPack(hash []byte) (int, int64, int64, error)
	LookupHash(hash []byte) (bool, error)
	AddCommit(trackID int, hashes [][]byte) (int, error)
	GetHashes(commitID int) ([][]byte, error)
}