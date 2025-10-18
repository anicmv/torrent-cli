package bencode

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// BEncoder 编码解码器
type BEncoder struct {
	charset string
}

// NewBEncoder 创建新的编码器
func NewBEncoder() *BEncoder {
	return &BEncoder{
		charset: "UTF-8",
	}
}

// Decode 解码 bencode 数据
func (b *BEncoder) Decode(data []byte) (interface{}, error) {
	decoder := NewDecoder(data)
	return decoder.Decode()
}

// Encode 编码数据为 bencode 格式
func (b *BEncoder) Encode(data interface{}) ([]byte, error) {
	encoder := NewEncoder()
	return encoder.Encode(data)
}

// Save 保存种子文件
func (b *BEncoder) Save(filename string, dict map[string]interface{}) error {
	data, err := b.Encode(dict)
	if err != nil {
		return err
	}

	// 获取文件所在目录
	dir := filepath.Dir(filename)

	// 如果目录不是当前目录，则创建
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return os.WriteFile(filename, data, 0644)
}

// SHA1 计算 info 字典的 SHA1 哈希
func (b *BEncoder) SHA1(torrentData []byte) (string, error) {
	decoded, err := b.Decode(torrentData)
	if err != nil {
		return "", err
	}

	dict, ok := decoded.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid torrent format")
	}

	info, ok := dict["info"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("info not found")
	}

	infoEncoded, err := b.Encode(info)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum(infoEncoded)
	return hex.EncodeToString(hash[:]), nil
}
