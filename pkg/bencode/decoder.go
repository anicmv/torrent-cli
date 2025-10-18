package bencode

import (
	"fmt"
	"strconv"
)

// Decoder bencode 解码器
type Decoder struct {
	data  []byte
	index int
}

// NewDecoder 创建解码器
func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		data:  data,
		index: 0,
	}
}

// Decode 解码数据
func (d *Decoder) Decode() (interface{}, error) {
	if d.index >= len(d.data) {
		return nil, fmt.Errorf("unexpected end of data")
	}

	switch d.data[d.index] {
	case 'i':
		return d.decodeInt()
	case 'l':
		return d.decodeList()
	case 'd':
		return d.decodeDict()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return d.decodeString()
	default:
		return nil, fmt.Errorf("invalid bencode at position %d", d.index)
	}
}

// decodeInt 解码整数
func (d *Decoder) decodeInt() (int64, error) {
	d.index++ // 跳过 'i'
	start := d.index

	for d.index < len(d.data) && d.data[d.index] != 'e' {
		d.index++
	}

	if d.index >= len(d.data) {
		return 0, fmt.Errorf("unexpected end while decoding integer")
	}

	numStr := string(d.data[start:d.index])
	d.index++ // 跳过 'e'

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

// decodeString 解码字符串
func (d *Decoder) decodeString() (string, error) {
	start := d.index

	for d.index < len(d.data) && d.data[d.index] != ':' {
		d.index++
	}

	if d.index >= len(d.data) {
		return "", fmt.Errorf("unexpected end while decoding string length")
	}

	lengthStr := string(d.data[start:d.index])
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	d.index++ // 跳过 ':'
	if d.index+length > len(d.data) {
		return "", fmt.Errorf("string length exceeds data")
	}

	str := string(d.data[d.index : d.index+length])
	d.index += length

	return str, nil
}

// decodeList 解码列表
func (d *Decoder) decodeList() ([]interface{}, error) {
	d.index++ // 跳过 'l'
	list := make([]interface{}, 0)

	for d.index < len(d.data) && d.data[d.index] != 'e' {
		item, err := d.Decode()
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	if d.index >= len(d.data) {
		return nil, fmt.Errorf("unexpected end while decoding list")
	}

	d.index++ // 跳过 'e'
	return list, nil
}

// decodeDict 解码字典
func (d *Decoder) decodeDict() (map[string]interface{}, error) {
	d.index++ // 跳过 'd'
	dict := make(map[string]interface{})

	for d.index < len(d.data) && d.data[d.index] != 'e' {
		key, err := d.decodeString()
		if err != nil {
			return nil, err
		}

		value, err := d.Decode()
		if err != nil {
			return nil, err
		}

		dict[key] = value
	}

	if d.index >= len(d.data) {
		return nil, fmt.Errorf("unexpected end while decoding dictionary")
	}

	d.index++ // 跳过 'e'
	return dict, nil
}
