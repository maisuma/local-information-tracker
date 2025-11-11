package index

type Indexer interface {

    AddTrack(filepath string) (int, error)

    RemoveTrack(trackID int) error

	GetTrackID(filepath string) (int, error)

	SaveHash(hash []byte, packID int, offset int64, size int64) error
}