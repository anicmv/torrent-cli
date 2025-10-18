package torrent

// Torrent 种子文件结构
type Torrent struct {
	Path            string
	Name            string
	ProtocolVersion string
	InfoHash        string
	MagnetURI       string
	MetaInfo        map[string]interface{}
}

// NewTorrent 创建新的种子对象
func NewTorrent(path, name string) *Torrent {
	return &Torrent{
		Path:            path,
		Name:            name,
		ProtocolVersion: "v1",
		MetaInfo:        make(map[string]interface{}),
	}
}

// MetaInfo 种子元信息
type MetaInfo struct {
	Announce     string
	AnnounceList []string
	CreationDate int64
	Comment      string
	CreatedBy    string
	Info         *Info
}

// Info info 字典信息
type Info struct {
	PieceLength int64
	Pieces      []byte
	PrivateFlag int64
	Source      string
	Name        string
	Length      int64
	Files       []*FileInfo
}

// FileInfo 文件信息
type FileInfo struct {
	Length int64
	Path   []string
}
