package bencode

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

// Encoder bencode 编码器
type Encoder struct {
	buf *bytes.Buffer
}

// NewEncoder 创建编码器
func NewEncoder() *Encoder {
	return &Encoder{
		buf: new(bytes.Buffer),
	}
}

// Encode 编码数据
func (e *Encoder) Encode(data interface{}) ([]byte, error) {
	e.buf.Reset()
	if err := e.encode(data); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

func (e *Encoder) encode(data interface{}) error {
	switch v := data.(type) {
	case int:
		return e.encodeInt(int64(v))
	case int64:
		return e.encodeInt(v)
	case string:
		return e.encodeString(v)
	case []byte:
		return e.encodeBytes(v)
	case []interface{}:
		return e.encodeList(v)
	case map[string]interface{}:
		return e.encodeDict(v)
	default:
		return fmt.Errorf("unsupported type: %T", data)
	}
}

func (e *Encoder) encodeInt(n int64) error {
	e.buf.WriteByte('i')
	e.buf.WriteString(strconv.FormatInt(n, 10))
	e.buf.WriteByte('e')
	return nil
}

func (e *Encoder) encodeString(s string) error {
	return e.encodeBytes([]byte(s))
}

func (e *Encoder) encodeBytes(b []byte) error {
	e.buf.WriteString(strconv.Itoa(len(b)))
	e.buf.WriteByte(':')
	e.buf.Write(b)
	return nil
}

func (e *Encoder) encodeList(list []interface{}) error {
	e.buf.WriteByte('l')
	for _, item := range list {
		if err := e.encode(item); err != nil {
			return err
		}
	}
	e.buf.WriteByte('e')
	return nil
}

func (e *Encoder) encodeDict(dict map[string]interface{}) error {
	// 按键排序
	keys := make([]string, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	e.buf.WriteByte('d')
	for _, k := range keys {
		if err := e.encodeString(k); err != nil {
			return err
		}
		if err := e.encode(dict[k]); err != nil {
			return err
		}
	}
	e.buf.WriteByte('e')
	return nil
}
