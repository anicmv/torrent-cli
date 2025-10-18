package torrent

import (
	"crypto/sha1"
	"io"
	"os"
	"path/filepath"
	"time"
	"torrent-cli/pkg/bencode"
)

// CreateOptions 创建选项
type CreateOptions struct {
	InputPath      string
	OutputPath     string
	AnnounceURL    string
	Comment        string
	TorrentName    string
	PieceLength    int
	PrivateFlag    bool
	Source         string
	CreatedBy      string
	Publisher      string
	CreationDate   int64
	NoCreatedBy    bool
	NoCreationDate bool
	NoPublisher    bool
	NoSource       bool
}

// Creator 种子创建器
type Creator struct {
	options *CreateOptions
	encoder *bencode.BEncoder
}

// NewCreator 创建种子创建器
func NewCreator(options *CreateOptions) *Creator {
	return &Creator{
		options: options,
		encoder: bencode.NewBEncoder(),
	}
}

// Create 创建种子文件
func (c *Creator) Create() error {
	fileInfo, err := os.Stat(c.options.InputPath)
	if err != nil {
		return err
	}

	var info map[string]interface{}
	if fileInfo.IsDir() {
		info, err = c.createInfoForDirectory()
	} else {
		info, err = c.createInfoForFile()
	}

	if err != nil {
		return err
	}

	metaInfo := c.createMetaInfo(info)

	outputPath := c.getOutputPath()
	return c.encoder.Save(outputPath, metaInfo)
}

func (c *Creator) createInfoForFile() (map[string]interface{}, error) {
	file, err := os.Open(c.options.InputPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	fileInfo, _ := file.Stat()
	totalSize := fileInfo.Size()

	if c.options.PieceLength == 0 {
		c.options.PieceLength = c.calculateOptimalPieceLength(totalSize)
	}

	pieces, err := c.hashPiecesForFile(file, c.options.PieceLength)
	if err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"piece length": int64(c.options.PieceLength),
		"pieces":       pieces,
		"name":         c.getTorrentName(),
		"length":       totalSize,
	}

	if c.options.PrivateFlag {
		info["private"] = int64(1)
	}

	if !c.options.NoSource && c.options.Source != "" {
		info["source"] = c.options.Source
	}

	return info, nil
}

func (c *Creator) createInfoForDirectory() (map[string]interface{}, error) {
	var files []string
	var totalSize int64

	err := filepath.Walk(c.options.InputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if c.options.PieceLength == 0 {
		c.options.PieceLength = c.calculateOptimalPieceLength(totalSize)
	}

	pieces, err := c.hashPiecesForFiles(files, c.options.PieceLength)
	if err != nil {
		return nil, err
	}

	fileList := make([]interface{}, 0)
	for _, filePath := range files {
		relPath, _ := filepath.Rel(c.options.InputPath, filePath)
		fileInfo, _ := os.Stat(filePath)

		pathParts := filepath.SplitList(relPath)
		if len(pathParts) == 0 {
			pathParts = []string{filepath.Base(relPath)}
		}

		fileDict := map[string]interface{}{
			"length": fileInfo.Size(),
			"path":   c.convertToInterfaceSlice(pathParts),
		}
		fileList = append(fileList, fileDict)
	}

	info := map[string]interface{}{
		"piece length": int64(c.options.PieceLength),
		"pieces":       pieces,
		"name":         c.getTorrentName(),
		"files":        fileList,
	}

	if c.options.PrivateFlag {
		info["private"] = int64(1)
	}

	if !c.options.NoSource && c.options.Source != "" {
		info["source"] = c.options.Source
	}

	return info, nil
}

func (c *Creator) hashPiecesForFile(file *os.File, pieceLength int) ([]byte, error) {
	var pieces []byte
	buffer := make([]byte, pieceLength)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		hash := sha1.Sum(buffer[:n])
		pieces = append(pieces, hash[:]...)

		if err == io.EOF {
			break
		}
	}

	return pieces, nil
}

func (c *Creator) hashPiecesForFiles(files []string, pieceLength int) ([]byte, error) {
	var pieces []byte
	buffer := make([]byte, pieceLength)
	currentPos := 0

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}

		for {
			n, err := file.Read(buffer[currentPos:])
			if err != nil && err != io.EOF {
				err := file.Close()
				if err != nil {
					return nil, err
				}
				return nil, err
			}

			currentPos += n

			if currentPos == pieceLength {
				hash := sha1.Sum(buffer)
				pieces = append(pieces, hash[:]...)
				currentPos = 0
			}

			if err == io.EOF {
				break
			}
		}

		err = file.Close()
		if err != nil {
			return nil, err
		}
	}

	// 处理最后一个不完整的块
	if currentPos > 0 {
		hash := sha1.Sum(buffer[:currentPos])
		pieces = append(pieces, hash[:]...)
	}

	return pieces, nil
}

func (c *Creator) createMetaInfo(info map[string]interface{}) map[string]interface{} {
	metaInfo := make(map[string]interface{})

	metaInfo["announce"] = c.options.AnnounceURL
	metaInfo["info"] = info

	if !c.options.NoCreatedBy && c.options.CreatedBy != "" {
		metaInfo["created by"] = c.options.CreatedBy
	}

	if !c.options.NoCreationDate {
		if c.options.CreationDate == 0 {
			metaInfo["creation date"] = time.Now().Unix()
		} else {
			metaInfo["creation date"] = c.options.CreationDate
		}
	}

	if c.options.Comment != "" {
		metaInfo["comment"] = c.options.Comment
	}

	if !c.options.NoPublisher && c.options.Publisher != "" {
		metaInfo["publisher"] = c.options.Publisher
	}

	return metaInfo
}

func (c *Creator) calculateOptimalPieceLength(totalSize int64) int {
	const (
		minPieces      = 512
		maxPieces      = 1024
		minPieceLength = 32 * 1024        // 32 KB
		maxPieceLength = 16 * 1024 * 1024 // 16 MB
	)

	pieceLength := totalSize / minPieces
	pieceLength = c.nearestPowerOfTwo(pieceLength)

	if pieceLength < minPieceLength {
		pieceLength = minPieceLength
	}
	if pieceLength > maxPieceLength {
		pieceLength = maxPieceLength
	}

	return int(pieceLength)
}

func (c *Creator) nearestPowerOfTwo(n int64) int64 {
	if n < 1 {
		return 1
	}
	power := int64(1)
	for power < n {
		power *= 2
	}
	return power
}

func (c *Creator) getTorrentName() string {
	if c.options.TorrentName != "" {
		return c.options.TorrentName
	}
	return filepath.Base(c.options.InputPath)
}

func (c *Creator) getOutputPath() string {
	if c.options.OutputPath != "" {
		// 如果指定了输出路径
		if filepath.Ext(c.options.OutputPath) == ".torrent" {
			// 如果已经是 .torrent 文件，直接返回
			return c.options.OutputPath
		}
		// 如果是目录，添加文件名
		return filepath.Join(c.options.OutputPath, c.getTorrentName()+".torrent")
	}
	// 默认在当前目录生成
	return c.getTorrentName() + ".torrent"
}

func (c *Creator) convertToInterfaceSlice(strs []string) []interface{} {
	result := make([]interface{}, len(strs))
	for i, s := range strs {
		result[i] = s
	}
	return result
}
