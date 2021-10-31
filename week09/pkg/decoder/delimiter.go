package decoder

import (
	"errors"
	"io"
	"net"
)

type DelimiterBasedFrameDecoder struct {
	reader    io.Reader
	maxLength uint32
	delimiter []byte
	buf       []byte
}

func (decoder *DelimiterBasedFrameDecoder) Read() ([]byte, error) {
	length := int(decoder.maxLength)
	data := make([]byte, 0, length)
	data = append(decoder.buf)
	for i := 0; len(data) < length; {
		size, err := decoder.reader.Read(data[i:length])
		if err != nil {
			return nil, err
		}
		idx := decoder.delimiterIndex(data)
		if idx != -1 {
			decoder.buf = data[idx+len(decoder.buf):]
			return data[0:idx], nil
		}
		i += size
	}
	return nil, errors.New("frame length exceeds the maxLength")
}

func (decoder *DelimiterBasedFrameDecoder) delimiterIndex(data []byte) int {
	for i := 0; i < len(data); i++ {
		haystackIndex, needleIndex := i, 0
		for ; needleIndex < len(decoder.delimiter); i++ {
			if data[haystackIndex] != decoder.delimiter[needleIndex] {
				break
			} else {
				haystackIndex++
				if haystackIndex == len(data) && needleIndex != len(decoder.delimiter)-1 {
					return -1
				}
			}
		}
		if needleIndex == len(data) {
			return i
		}
	}
	return -1
}

func NewDelimiterBasedFrameDecoder(conn net.Conn, maxLength uint32, delimiter string) FrameDecoder {
	if delimiter == "" {
		panic("delimiter must be non-empty")
	}
	return &DelimiterBasedFrameDecoder{
		reader:    conn,
		delimiter: []byte(delimiter),
		maxLength: maxLength,
	}
}
